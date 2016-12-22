package network


import  (
	"golang.org/x/net/context"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type NetworksByIdNode struct {
	utils.DynamicDir	
	dockerClient  *client.Client
}


func NewNetworksByIdNode( pDockerClient *client.Client) (*NetworksByIdNode, error) {
	vRis:=&NetworksByIdNode{dockerClient: pDockerClient}
	vRis.DynamicDir.CacheInterval= time.Second*5
	vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *NetworksByIdNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vNetworks,vError:=vSelf.dockerClient.NetworkList(context.Background(),types.NetworkListOptions{})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list networks",vError)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vNetworks))
	for _,vCurNetwork := range vNetworks {
		vCurNetworkNode,vCurNetworkNodeError:=NewNetworkNode(vSelf.dockerClient,vCurNetwork)
		if vCurNetworkNodeError!= nil {
			return nil,diagnostic.NewError("Failed to create network Node for network %#v", vCurNetworkNodeError, vCurNetwork)
		}
		
		vRis[vCurNetwork.ID] = vCurNetworkNode
	}
	return vRis,nil
}

