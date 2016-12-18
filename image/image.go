package image


import  (
	"encoding/json" 
	"os"
	"golang.org/x/net/context"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/mysinmyc/dockfs/utils"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type ImageNode struct {
	dockerClient  *client.Client
	image     types.ImageSummary
	children      map[string] utils.DirentTyped
}


func NewImageNode(pDockerClient *client.Client, pImage types.ImageSummary) (*ImageNode, error) {
	vRis:= &ImageNode {dockerClient: pDockerClient, image: pImage}
        vInitError:=vRis.init(context.Background())

        if vInitError != nil {
                return nil,diagnostic.NewError("initialization failed",vInitError)
        }

	
	return vRis,nil
}

func (vSelf *ImageNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {
	pAttr.Mode = os.ModeDir | 0555
	return nil
}

func (vSelf *ImageNode) Lookup(pContext context.Context, pName string) (fs.Node, error) {
        vRis:=vSelf.children[pName]

        if vRis != nil {
                return vRis.(fs.Node),nil
        }
        return nil, fuse.ENOENT
}

func (vSelf *ImageNode) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
        vRis:= make([]fuse.Dirent,0)
        for vCurChildName, vCurChild := range vSelf.children {
                vRis=append(vRis, fuse.Dirent{ Type: vCurChild.GetDirentType(), Name:vCurChildName})   
        }
        return vRis,nil
}

func (vSelf *ImageNode) init(pContext context.Context) (error) {

        vChildren:= make(map[string]utils.DirentTyped)
	vChildren["json"],_ = utils.NewDynamicFileNode(vSelf.jsonFunc)
	if vSelf.image.ParentID != "" {
		vChildren["parent"],_ = utils.NewSymLinkNode("../"+vSelf.image.ParentID)
	}
        vSelf.children = vChildren
        return nil
}



func (vSelf *ImageNode) jsonFunc() ([]byte,error) {
		return json.Marshal(vSelf.image)
}

