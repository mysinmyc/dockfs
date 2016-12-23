package main

import(
	"flag"
	//"golang.org/x/net/context"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs"
)


func buildClient() (*client.Client,error) {
	diagnostic.LogInfo("buildClient", "Creating docker client...")
	return client.NewEnvClient() 
}

func negotiateApiVersion(pClient *client.Client) error {

		//vPing,vPingError := pClient.Ping(context.Background())
		//diagnostic.LogFatalIfError(vPingError,"negotiateApiVersion", "An error occurred while asking server version")
		//pClient.UpdateClientVersion(vPing.APIVersion)
		pClient.UpdateClientVersion("") //Blank means no api check
		return nil
}

func mount(pMountPoint string, pClient *client.Client, pAllowOther bool) {

	 
	diagnostic.LogInfo("mount", "Creating filesystem...")
	vFs,vFsError:= dockfs.NewDockFs(pClient)	
	diagnostic.LogFatalIfError(vFsError,"mount","Failed to create filesystem")
		
	diagnostic.LogInfo("mount", "Mounting filesystem on %s...", pMountPoint)

	vMountOptions := make([]fuse.MountOption,0)
	vMountOptions=append(vMountOptions,fuse.FSName("dockfs"),
                fuse.ReadOnly(),
                fuse.LocalVolume())
	if (pAllowOther) {
		vMountOptions=append(vMountOptions,fuse.AllowOther())
	}
        vFileSystemConnection, vMountError := fuse.Mount(
                pMountPoint,
                vMountOptions...)
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

	vAllowOtherParameter := flag.Bool("allowOther", false, "Allow Other")
	vAutoApiVersionParameter := flag.Bool("autoApiVersion", true, "Automatic api version configuration")
	vLogLevelParameter := flag.Int("logLevel", diagnostic.LogLevel_Info, "Log level")
	vMountPointParameter := flag.String("mountPoint", "", "Target mountpoint")
	vActionParameter := flag.String("action", "mount", "action [mount|umount]")
	flag.Parse()

	diagnostic.SetLogLevel(diagnostic.LogLevel(*vLogLevelParameter))

	if *vMountPointParameter == "" {
		diagnostic.LogFatal("main","missing mountPoint parameter",nil)
	}



	switch *vActionParameter {
		case "mount":
			vClient,vClientError:= buildClient()	
			
			diagnostic.LogFatalIfError(vClientError,"main", "An error occurred while creating docker client")

			if *vAutoApiVersionParameter {
				vNegotiateApiError:= negotiateApiVersion(vClient)
				diagnostic.LogFatalIfError(vNegotiateApiError, "main", "Failed to negotiate the api version")
			}

			mount(*vMountPointParameter,vClient,*vAllowOtherParameter)
			break
		case "umount":
			umount(*vMountPointParameter)
			break
		default:
			diagnostic.LogFatal("main", "invalid action %s",nil, *vActionParameter)
	}

	
}
