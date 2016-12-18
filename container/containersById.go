package container


import  (
	"golang.org/x/net/context"
	"os"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type ContainersByIdNode struct {
	dockerClient  *client.Client
	containersNodes map[string]*ContainerNode
}


func NewContainersByIdNode( pDockerClient *client.Client) (*ContainersByIdNode, error) {
	vRis:=&ContainersByIdNode {dockerClient: pDockerClient}

	vInitError:=vRis.init(context.Background())

	if vInitError != nil {
		return nil,diagnostic.NewError("initialization failed",vInitError)
	}		
	return vRis,nil
}

func (vSelf *ContainersByIdNode)  Attr(pContext context.Context, pAttr *fuse.Attr) error {
	pAttr.Mode = os.ModeDir | 0555
	return nil
}

func (vSelf *ContainersByIdNode) Lookup(pContext context.Context, pName string) (fs.Node, error) {
	vRis:=vSelf.containersNodes[pName]

	if vRis != nil {
		return vRis,nil	
	}
	return nil, fuse.ENOENT
}

func (vSelf *ContainersByIdNode) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
	vRis:= make([]fuse.Dirent,0)
	for vCurContainerId, _ := range vSelf.containersNodes {
		vRis=append(vRis, fuse.Dirent{ Type: fuse.DT_Dir, Name:vCurContainerId})     	
	}
	return vRis,nil
}

func (vSelf *ContainersByIdNode) init(pContext context.Context) (error) {

	vContainers,vError:=vSelf.dockerClient.ContainerList(pContext,types.ContainerListOptions{All:true, Limit:-1})
	if vError!=nil {
		return diagnostic.NewError("failed to list containers",vError)
	}

	vContainersNodes:= make(map[string]*ContainerNode,len(vContainers))
	for _,vCurContainer := range vContainers {
		vCurContainerNode,vCurContainerNodeError:=NewContainerNode(vSelf.dockerClient,vCurContainer)
		if vCurContainerNodeError!= nil {
			return diagnostic.NewError("Failed to create container Node for container %#v", vCurContainerNodeError, vCurContainer)
		}
		
		vContainersNodes[vCurContainer.ID] = vCurContainerNode
	}
	vSelf.containersNodes = vContainersNodes
	return nil
}

