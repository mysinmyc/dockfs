package main

import(
	"flag"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs"
)

func main() {

	vMountPointParameter := flag.String("mountPoint", "", "Target mountpoint")
	flag.Parse()

	if *vMountPointParameter == "" {
		diagnostic.LogFatal("main","missing mountPoint parameter",nil)
	}

	diagnostic.LogInfo("main", "Creating docker client...")
	vClient, vClientError := client.NewEnvClient() 
	diagnostic.LogFatalIfError(vClientError,"main","Failed to create client")
	 
	diagnostic.LogInfo("main", "Creating filesystem...")
	vFs,vFsError:= dockfs.NewDockFs(vClient)	
	diagnostic.LogFatalIfError(vFsError,"main","Failed to create filesystem")
		
	diagnostic.LogInfo("main", "Mounting filesystem on %s...", *vMountPointParameter)
        vFileSystemConnection, vMountError := fuse.Mount(
                *vMountPointParameter,
                fuse.FSName("dockfs"),
                fuse.ReadOnly(),
                fuse.LocalVolume())
	diagnostic.LogFatalIfError(vMountError,"main","An error occurred while mounting onedrive filesystem")
       
        defer vFileSystemConnection.Close()

	diagnostic.LogInfo("main", "Filesystem mounted")
        vServeError := fs.Serve(vFileSystemConnection, vFs )
	diagnostic.LogFatalIfError(vServeError,"main", "An error occurred while serving filesystem")

        <-vFileSystemConnection.Ready
        if vFileSystemConnection.MountError != nil {
                diagnostic.LogFatal("main", "Mount failed", vFileSystemConnection.MountError)

        }

	
}
