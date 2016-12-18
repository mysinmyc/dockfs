package image


import  (
	"golang.org/x/net/context"
	"os"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
)


type ImagesByIdNode struct {
	dockerClient  *client.Client
	imagesNodes map[string]*ImageNode
}


func NewImagesByIdNode( pDockerClient *client.Client) (*ImagesByIdNode, error) {
	vRis:=&ImagesByIdNode {dockerClient: pDockerClient}

	vInitError:=vRis.init(context.Background())

	if vInitError != nil {
		return nil,diagnostic.NewError("initialization failed",vInitError)
	}		
	return vRis,nil
}

func (vSelf *ImagesByIdNode)  Attr(pContext context.Context, pAttr *fuse.Attr) error {
	pAttr.Mode = os.ModeDir | 0555
	return nil
}

func (vSelf *ImagesByIdNode) Lookup(pContext context.Context, pName string) (fs.Node, error) {
	vRis:=vSelf.imagesNodes[pName]

	if vRis != nil {
		return vRis,nil	
	}
	return nil, fuse.ENOENT
}

func (vSelf *ImagesByIdNode) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
	vRis:= make([]fuse.Dirent,0)
	for vCurImageId, _ := range vSelf.imagesNodes {
		vRis=append(vRis, fuse.Dirent{ Type: fuse.DT_Dir, Name:vCurImageId})     	
	}
	return vRis,nil
}

func (vSelf *ImagesByIdNode) init(pContext context.Context) (error) {

	vImages,vError:=vSelf.dockerClient.ImageList(pContext,types.ImageListOptions{All:true})
	if vError!=nil {
		return diagnostic.NewError("failed to list images",vError)
	}

	vImagesNodes:= make(map[string]*ImageNode,len(vImages))
	for _,vCurImage := range vImages {
		vCurImageNode,vCurImageNodeError:=NewImageNode(vSelf.dockerClient,vCurImage)
		if vCurImageNodeError!= nil {
			return diagnostic.NewError("Failed to create image Node for container %#v", vCurImageNodeError, vCurImage)
		}
		
		vImagesNodes[vCurImage.ID] = vCurImageNode
	}
	vSelf.imagesNodes = vImagesNodes
	return nil
}

