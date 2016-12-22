package network


import  (
	"encoding/json" 
	"os"
	"golang.org/x/net/context"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/mysinmyc/dockfs/utils"
	"github.com/mysinmyc/dockfs/container"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type NetworkNode struct {
	dockerClient  *client.Client
	network     types.NetworkResource	
	children      map[string] utils.DirentTyped
}


func NewNetworkNode(pDockerClient *client.Client, pNetwork types.NetworkResource) (*NetworkNode, error) {
	vRis:= &NetworkNode {dockerClient: pDockerClient, network: pNetwork}
        vInitError:=vRis.init(context.Background())

        if vInitError != nil {
                return nil,diagnostic.NewError("initialization failed",vInitError)
        }

	
	return vRis,nil
}

func (vSelf *NetworkNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {
	pAttr.Mode = os.ModeDir | 0555
	return nil
}

func (vSelf *NetworkNode) Lookup(pContext context.Context, pName string) (fs.Node, error) {
        vRis:=vSelf.children[pName]

        if vRis != nil {
                return vRis.(fs.Node),nil
        }
        return nil, fuse.ENOENT
}

func (vSelf *NetworkNode) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
        vRis:= make([]fuse.Dirent,0)
        for vCurChildName, vCurChild := range vSelf.children {
                vRis=append(vRis, fuse.Dirent{ Type: vCurChild.GetDirentType(), Name:vCurChildName})   
        }
        return vRis,nil
}

func (vSelf *NetworkNode) init(pContext context.Context) (error) {

        vChildren:= make(map[string]utils.DirentTyped)
	vChildren["name"],_ = utils.NewDynamicFileNode(vSelf.nameFunc)
	vChildren["json"],_ = utils.NewDynamicFileNode(vSelf.jsonFunc)
	vChildren["containers"],_ = container.NewContainersOfNetworkNode(vSelf.network.ID,vSelf.dockerClient)
        vSelf.children = vChildren
        return nil
}


func (vSelf *NetworkNode) GetDirentType() (fuse.DirentType) {
	return fuse.DT_Dir
}

func (vSelf *NetworkNode) nameFunc() ([]byte,error) {
	return []byte(vSelf.network.Name), nil
}

func (vSelf *NetworkNode) jsonFunc() ([]byte,error) {
	return json.Marshal(vSelf.network)
}

