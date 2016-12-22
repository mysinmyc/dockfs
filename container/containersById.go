package container


import  (
	"golang.org/x/net/context"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type ContainersByIdNode struct {
	utils.DynamicDir	
	dockerClient  *client.Client
}


func NewContainersByIdNode( pDockerClient *client.Client) (*ContainersByIdNode, error) {
	vRis:=&ContainersByIdNode{dockerClient: pDockerClient}
	vRis.DynamicDir.CacheInterval= time.Second*5
	vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ContainersByIdNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vContainers,vError:=vSelf.dockerClient.ContainerList(context.Background(),types.ContainerListOptions{All:true, Limit:-1})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list containers",vError)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vContainers))
	for _,vCurContainer := range vContainers {
		vCurContainerNode,vCurContainerNodeError:=NewContainerNode(vSelf.dockerClient,vCurContainer)
		if vCurContainerNodeError!= nil {
			return nil,diagnostic.NewError("Failed to create container Node for container %#v", vCurContainerNodeError, vCurContainer)
		}
		
		vRis[vCurContainer.ID] = vCurContainerNode
	}
	return vRis,nil
}

