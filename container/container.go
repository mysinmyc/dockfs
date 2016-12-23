package container


import  (
	"encoding/json" 
	"os"
	"io/ioutil"
	"strings"
	"golang.org/x/net/context"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/mysinmyc/dockfs/utils"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type ContainerNode struct {
	dockerClient  *client.Client
	containerId   string
	container     *types.Container 
	children      map[string] utils.DirentTyped
}


func NewContainerNode(pDockerClient *client.Client, pContainerId string) (*ContainerNode, error) {
	vRis:= &ContainerNode {dockerClient: pDockerClient, containerId: pContainerId}
        vInitError:=vRis.init()

        if vInitError != nil {
                return nil,diagnostic.NewError("initialization failed",vInitError)
        }
	
	return vRis,nil
}

func (vSelf *ContainerNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {
	/*vContainer,vContainerError:=vSelf.getContainer(true)
	if vContainerError != nil {
		return diagnostic.NewError("failed to get container",vContainerError)
	}
	pAttr.Mtime = time.Unix(vContainer.Created,0)
	*/

	pAttr.Mode = os.ModeDir | 0555
	return nil
}

func (vSelf *ContainerNode) Lookup(pContext context.Context, pName string) (fs.Node, error) {
        vRis:=vSelf.children[pName]

        if vRis != nil {
                return vRis.(fs.Node),nil
        }
        return nil, fuse.ENOENT
}

func (vSelf *ContainerNode) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
        vRis:= make([]fuse.Dirent,0)
        for vCurChildName, vCurChild := range vSelf.children {
                vRis=append(vRis, fuse.Dirent{ Type: vCurChild.GetDirentType(), Name:vCurChildName})   
        }
        return vRis,nil
}

func (vSelf *ContainerNode) init() (error) {

        vChildren:= make(map[string]utils.DirentTyped)
	//vChildren["."] = vSelf
	vChildren["name"],_ = utils.NewDynamicFileNode(vSelf.nameFunc)
	vChildren["json"],_ = utils.NewDynamicFileNode(vSelf.jsonFunc)
	vChildren["state"],_ = utils.NewDynamicFileNode(vSelf.stateFunc)
	vChildren["stdout"],_ = utils.NewDynamicFileNode(vSelf.stdOutFunc)
	vChildren["stderr"],_ = utils.NewDynamicFileNode(vSelf.stdErrFunc)
	vChildren["command"],_ = utils.NewDynamicFileNode(vSelf.commandFunc)

	vContainer,vContainerError:=vSelf.getContainer(true)
	if vContainerError != nil {
		return diagnostic.NewError("failed to get container", vContainerError)
	}

	vChildren["image"],_ = utils.NewSymLinkNode("../../../images/byId/"+vContainer.ImageID)
        vSelf.children = vChildren
        return nil
}

func (vSelf *ContainerNode) getContainer(pCached bool) (*types.Container,error) {

	vFilters:=filters.NewArgs()
        vFilters.Add("id", vSelf.containerId)

	if pCached == false || vSelf.container == nil {
		vContainers,vContainersError:= vSelf.dockerClient.ContainerList(context.Background(),types.ContainerListOptions{Filters:vFilters,All:true} )
		if vContainersError!=nil {
			return nil,vContainersError
		}	
		if len(vContainers) != 1 {
			return nil, diagnostic.NewError("Invalid number of containers with id %s found: %d", nil, vSelf.containerId, len(vContainers))
		}
		vSelf.container = &vContainers[0]
	}

	return vSelf.container,nil
}

func (vSelf *ContainerNode) GetDirentType() (fuse.DirentType) {
	return fuse.DT_Dir
}

func (vSelf *ContainerNode) nameFunc() ([]byte,error) {
	vContainer,vContainerError:=vSelf.getContainer(true)
	if vContainerError != nil {
		return nil,diagnostic.NewError("failed to get container", vContainerError)
	}
	return []byte(strings.Join(vContainer.Names,"/")), nil
}

func (vSelf *ContainerNode) jsonFunc() ([]byte,error) {
	vContainer,vContainerError:=vSelf.getContainer(false)
	if vContainerError != nil {
		return nil,diagnostic.NewError("failed to get container", vContainerError)
	}
	return json.Marshal(vContainer)
}

func (vSelf *ContainerNode) stateFunc() ([]byte,error) {
	vContainer,vContainerError:=vSelf.getContainer(false)
	if vContainerError != nil {
		return nil,diagnostic.NewError("failed to get container", vContainerError)
	}
	return []byte(vContainer.State), nil
}

func (vSelf *ContainerNode) commandFunc() ([]byte,error) {
	vContainer,vContainerError:=vSelf.getContainer(true)
	if vContainerError != nil {
		return nil,diagnostic.NewError("failed to get container", vContainerError)
	}
	return []byte(vContainer.Command), nil
}

func (vSelf *ContainerNode) stdOutFunc() ([]byte,error) {
	return vSelf.getContainerLogs(types.ContainerLogsOptions{ShowStdout:true, Timestamps:true})
}

func (vSelf *ContainerNode) stdErrFunc() ([]byte,error) {
	return vSelf.getContainerLogs(types.ContainerLogsOptions{ShowStderr:true, Timestamps:true})
}

func (vSelf *ContainerNode) getContainerLogs(pOptions types.ContainerLogsOptions) ([]byte,error) {
	vLogsReader, vLogsError :=vSelf.dockerClient.ContainerLogs(context.Background(), vSelf.containerId, pOptions)

	if vLogsError != nil {
		return nil,diagnostic.NewError("Error getting container logs",vLogsError)
	}
	defer vLogsReader.Close()
	return ioutil.ReadAll(vLogsReader)
}
