package dockfs


import  (
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/container"
	"github.com/mysinmyc/dockfs/image"
	"github.com/mysinmyc/dockfs/network"
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

	vContainersByName, vContainersByNameError:=container.NewContainersByNameNode(vSelf.dockerClient)
	if vContainersByNameError != nil {
		return nil,diagnostic.NewError("Failed to create ContainersByName node",vContainersByNameError)
	}
	vRoot.Add("containers/byName",vContainersByName)

	vContainersByState, vContainersByStateError:=container.NewContainersByStateNode(vSelf.dockerClient)
	if vContainersByStateError != nil {
		return nil,diagnostic.NewError("Failed to create ContainersByState node",vContainersByStateError)
	}
	vRoot.Add("containers/byState",vContainersByState)

	vImagesById, vImagesByIdError:=image.NewImagesByIdNode(vSelf.dockerClient)
	if vImagesByIdError != nil {
		return nil,diagnostic.NewError("Failed to create ImagesById node",vImagesByIdError)
	}
	vRoot.Add("images/byId",vImagesById)

	vImagesByTag, vImagesByTagError:=image.NewImagesByTagNode(vSelf.dockerClient)
	if vImagesByTagError != nil {
		return nil,diagnostic.NewError("Failed to create ImagesByTag node",vImagesByTagError)
	}
	vRoot.Add("images/byTag",vImagesByTag)
	
	vNetworksById, vNetworksByIdError:=network.NewNetworksByIdNode(vSelf.dockerClient)
	if vNetworksByIdError != nil {
		return nil,diagnostic.NewError("Failed to create NetworksById node",vNetworksByIdError)
	}
	vRoot.Add("networks/byId",vNetworksById)

	vNetworksByName, vNetworksByNameError:=network.NewNetworksByNameNode(vSelf.dockerClient)
	if vNetworksByNameError != nil {
		return nil,diagnostic.NewError("Failed to create NetworksByName node",vNetworksByNameError)
	}
	vRoot.Add("networks/byName",vNetworksByName)
	
	return vRoot, nil
}
