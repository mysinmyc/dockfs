package container


import  (
	"golang.org/x/net/context"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type ContainersByStateNode struct {
	fs.Tree
	dockerClient  *client.Client
}


func NewContainersByStateNode( pDockerClient *client.Client) (*ContainersByStateNode, error) {
	vRis:=&ContainersByStateNode {dockerClient: pDockerClient}

	vInitError:=vRis.init(context.Background())

	if vInitError != nil {
		return nil,diagnostic.NewError("initialization failed",vInitError)
	}		
	return vRis,nil
}


func (vSelf *ContainersByStateNode) init(pContext context.Context) (error) {

	vContainers,vError:=vSelf.dockerClient.ContainerList(pContext,types.ContainerListOptions{All:true, Limit:-1})
	if vError!=nil {
		return diagnostic.NewError("failed to list containers",vError)
	}

	for _,vCurContainer := range vContainers {
		vCurSymLink,_:=utils.NewSymLinkNode("../../byId/"+vCurContainer.ID)
		vSelf.Add(vCurContainer.State+"/"+vCurContainer.ID, vCurSymLink)
 
	}
	return nil
}

