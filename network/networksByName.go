package network


import  (
	"time"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type NetworksByNameNode struct {
	utils.DynamicDir
	dockerClient  *client.Client
}


func NewNetworksByNameNode( pDockerClient *client.Client) (*NetworksByNameNode, error) {
	vRis:=&NetworksByNameNode {dockerClient: pDockerClient}
        vRis.DynamicDir.CacheInterval= time.Second*5
        vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *NetworksByNameNode)  populateChildren() (map[string] utils.DirentTyped,error) {
	vNetworks,vError:=vSelf.dockerClient.NetworkList(context.Background(),types.NetworkListOptions{})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list networks",vError)
	}
	vRis := make(map[string] utils.DirentTyped,0)
	for _,vCurNetwork := range vNetworks {
		vCurSymLink,_:=utils.NewSymLinkNode("../byId/"+vCurNetwork.ID)
		utils.AddDirentTo(vCurNetwork.Name, vCurSymLink, vRis)
 
	}
	return vRis,nil
}

