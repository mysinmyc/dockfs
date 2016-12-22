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
	"github.com/mysinmyc/dockfs/utils"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type ContainerNode struct {
	dockerClient  *client.Client
	container     types.Container	
	children      map[string] utils.DirentTyped
}


func NewContainerNode(pDockerClient *client.Client, pContainer types.Container) (*ContainerNode, error) {
	vRis:= &ContainerNode {dockerClient: pDockerClient, container: pContainer}
        vInitError:=vRis.init(context.Background())

        if vInitError != nil {
                return nil,diagnostic.NewError("initialization failed",vInitError)
        }

	
	return vRis,nil
}

func (vSelf *ContainerNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {
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

func (vSelf *ContainerNode) init(pContext context.Context) (error) {

        vChildren:= make(map[string]utils.DirentTyped)
	vChildren["name"],_ = utils.NewDynamicFileNode(vSelf.nameFunc)
	vChildren["json"],_ = utils.NewDynamicFileNode(vSelf.jsonFunc)
	vChildren["state"],_ = utils.NewDynamicFileNode(vSelf.stateFunc)
	vChildren["stdout"],_ = utils.NewDynamicFileNode(vSelf.stdOutFunc)
	vChildren["stderr"],_ = utils.NewDynamicFileNode(vSelf.stdErrFunc)
	vChildren["command"],_ = utils.NewDynamicFileNode(vSelf.commandFunc)
	vChildren["image"],_ = utils.NewSymLinkNode("../../../images/byId/"+vSelf.container.ImageID)
        vSelf.children = vChildren
        return nil
}


func (vSelf *ContainerNode) GetDirentType() (fuse.DirentType) {
	return fuse.DT_Dir
}

func (vSelf *ContainerNode) nameFunc() ([]byte,error) {
	return []byte(strings.Join(vSelf.container.Names,"/")), nil
}

func (vSelf *ContainerNode) jsonFunc() ([]byte,error) {
	return json.Marshal(vSelf.container)
}

func (vSelf *ContainerNode) stateFunc() ([]byte,error) {
	return []byte(vSelf.container.State), nil
}

func (vSelf *ContainerNode) commandFunc() ([]byte,error) {
	return []byte(vSelf.container.Command), nil
}

func (vSelf *ContainerNode) stdOutFunc() ([]byte,error) {
	return vSelf.getContainerLogs(types.ContainerLogsOptions{ShowStdout:true, Timestamps:true})
}

func (vSelf *ContainerNode) stdErrFunc() ([]byte,error) {
	return vSelf.getContainerLogs(types.ContainerLogsOptions{ShowStderr:true, Timestamps:true})
}

func (vSelf *ContainerNode) getContainerLogs(pOptions types.ContainerLogsOptions) ([]byte,error) {
	vLogsReader, vLogsError :=vSelf.dockerClient.ContainerLogs(context.Background(), vSelf.container.ID, pOptions)

	if vLogsError != nil {
		return nil,diagnostic.NewError("Error getting container logs",vLogsError)
	}
	defer vLogsReader.Close()
	return ioutil.ReadAll(vLogsReader)
}
