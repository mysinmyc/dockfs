package container


import  (
	"time"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type ContainersByStateNode struct {
	utils.DynamicDir
	dockerClient  *client.Client
}


func NewContainersByStateNode( pDockerClient *client.Client) (*ContainersByStateNode, error) {
	vRis:=&ContainersByStateNode {dockerClient: pDockerClient}
        vRis.DynamicDir.CacheInterval= time.Second*5
        vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ContainersByStateNode)  populateChildren() (map[string] utils.DirentTyped,error) {
	vContainers,vError:=vSelf.dockerClient.ContainerList(context.Background(),types.ContainerListOptions{All:true, Limit:-1})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list containers",vError)
	}
	vRis := make(map[string] utils.DirentTyped,0)
	for _,vCurContainer := range vContainers {
		vCurSymLink,_:=utils.NewSymLinkNode("../../byId/"+vCurContainer.ID)
		utils.AddDirentTo(vCurContainer.State+"/"+vCurContainer.ID, vCurSymLink, vRis)
 
	}
	return vRis,nil
}

