package container


import  (
	"golang.org/x/net/context"
	"time"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type ContainersOfNetworkNode struct {
	utils.DynamicDir	
	dockerClient  		*client.Client
	networkId			string
}


func NewContainersOfNetworkNode(pNetworkId string, pDockerClient *client.Client) (*ContainersOfNetworkNode, error) {
	vRis:=&ContainersOfNetworkNode{networkId:pNetworkId, dockerClient: pDockerClient}
	vRis.DynamicDir.CacheInterval= time.Second*5
	vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ContainersOfNetworkNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vNetwork,vError:=vSelf.dockerClient.NetworkInspect(context.Background(),vSelf.networkId)
	if vError!=nil {
		return nil, diagnostic.NewError("failed to inspect network %s",vError,vSelf.networkId)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vNetwork.Containers))
	for vCurContainerId,_ := range vNetwork.Containers {
		vSymLink,_:=utils.NewSymLinkNode("../../../../containers/byId/"+vCurContainerId)
		vRis[vCurContainerId] = vSymLink
	}
	return vRis,nil
}

