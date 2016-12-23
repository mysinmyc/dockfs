package container


import  (
	"golang.org/x/net/context"
	"time"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type NetworksOfContainerNode struct {
	utils.DynamicDir	
	dockerClient  		*client.Client
	containerId			string
}


func NewNetworksOfContainerNode(pContainerId string, pDockerClient *client.Client) (*NetworksOfContainerNode, error) {
	vRis:=&NetworksOfContainerNode{containerId:pContainerId, dockerClient: pDockerClient}
	vRis.DynamicDir.CacheInterval= time.Second*5
	vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *NetworksOfContainerNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vContainer,vError:=vSelf.dockerClient.ContainerInspect(context.Background(),vSelf.containerId)
	if vError!=nil {
		return nil, diagnostic.NewError("failed to inspect container %s",vError,vSelf.containerId)
	}

	vRis:= make(map[string]utils.DirentTyped,0)
	for vCurNetworkName,vCurNetwork := range vContainer.NetworkSettings.Networks {
		vSymLink,_:=utils.NewSymLinkNode("../../../../networks/byId/"+vCurNetwork.NetworkID)
		vRis[vCurNetworkName] = vSymLink
	}
	return vRis,nil
}

