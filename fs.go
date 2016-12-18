package dockfs


import  (
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/container"
	"github.com/mysinmyc/dockfs/image"
)


type DockFs struct {
	dockerClient  *client.Client
}


func NewDockFs(pDockerClient *client.Client) (*DockFs, error) {
	return &DockFs {dockerClient: pDockerClient}, nil
}

func (vSelf *DockFs) Root() (fs.Node, error) {
	vRoot:= &fs.Tree{}

	vContainersById, vContainersByIdError:=container.NewContainersByIdNode(vSelf.dockerClient)
	if vContainersByIdError != nil {
		return nil,diagnostic.NewError("Failed to create ContainersById node",vContainersByIdError)
	}
	vRoot.Add("containers/byId",vContainersById)

	vContainersByState, vContainersByStateError:=container.NewContainersByStateNode(vSelf.dockerClient)
	if vContainersByStateError != nil {
		return nil,diagnostic.NewError("Failed to create ContainersByState node",vContainersByStateError)
	}
	vRoot.Add("containers/byState",vContainersByState)

	vImagesById, vImagesByIdError:=image.NewImagesByIdNode(vSelf.dockerClient)
	if vImagesByIdError != nil {
		return nil,diagnostic.NewError("Failed to create ImagesById node",vContainersByIdError)
	}
	vRoot.Add("images/byId",vImagesById)
	return vRoot, nil
}
