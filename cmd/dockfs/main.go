package main

import(
	"flag"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs"
)


func mount(pMountPoint string) {

	diagnostic.LogInfo("mount", "Creating docker client...")
	vClient, vClientError := client.NewEnvClient() 
	diagnostic.LogFatalIfError(vClientError,"mount","Failed to create client")
	 
	diagnostic.LogInfo("mount", "Creating filesystem...")
	vFs,vFsError:= dockfs.NewDockFs(vClient)	
	diagnostic.LogFatalIfError(vFsError,"mount","Failed to create filesystem")
		
	diagnostic.LogInfo("mount", "Mounting filesystem on %s...", pMountPoint)
        vFileSystemConnection, vMountError := fuse.Mount(
                pMountPoint,
                fuse.FSName("dockfs"),
                fuse.ReadOnly(),
                fuse.LocalVolume())
	diagnostic.LogFatalIfError(vMountError,"mount","An error occurred while mounting dockfs filesystem")
       
        defer vFileSystemConnection.Close()

	diagnostic.LogInfo("mount", "Filesystem mounted")
        vServeError := fs.Serve(vFileSystemConnection, vFs )
	diagnostic.LogFatalIfError(vServeError,"mount", "An error occurred while serving filesystem")

        <-vFileSystemConnection.Ready
        if vFileSystemConnection.MountError != nil {
                diagnostic.LogFatal("mount", "Mount failed", vFileSystemConnection.MountError)

        }

}

func umount(pMountPoint string) {
	diagnostic.LogInfo("umount", "Unmounting %s",pMountPoint)
	vUmountError := fuse.Unmount(pMountPoint)
	diagnostic.LogFatalIfError(vUmountError,"umount","An error occurred while umounting dockfs filesystem at %s",pMountPoint)
}

func main() {

	vMountPointParameter := flag.String("mountPoint", "", "Target mountpoint")
	vActionParameter := flag.String("action", "mount", "action [mount|umount]")
	flag.Parse()

	if *vMountPointParameter == "" {
		diagnostic.LogFatal("main","missing mountPoint parameter",nil)
	}


	switch *vActionParameter {
		case "mount":
			mount(*vMountPointParameter)
			break
		case "umount":
			umount(*vMountPointParameter)
			break
		default:
			diagnostic.LogFatal("main", "invalid action %s",nil, *vActionParameter)
	}

	
}
