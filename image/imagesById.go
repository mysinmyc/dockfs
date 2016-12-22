package image


import  (
	"time"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
)


type ImagesByIdNode struct {
	utils.DynamicDir
	dockerClient  *client.Client
}


func NewImagesByIdNode( pDockerClient *client.Client) (*ImagesByIdNode, error) {
	vRis:=&ImagesByIdNode {dockerClient: pDockerClient}
        vRis.DynamicDir.CacheInterval= time.Second*5
        vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ImagesByIdNode) populateChildren() (map[string] utils.DirentTyped,error) {


	vImages,vError:=vSelf.dockerClient.ImageList(context.Background(),types.ImageListOptions{All:true})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list images",vError)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vImages))
	for _,vCurImage := range vImages {
		vCurImageNode,vCurImageNodeError:=NewImageNode(vSelf.dockerClient,vCurImage)
		if vCurImageNodeError!= nil {
			return nil, diagnostic.NewError("Failed to create image Node for image %#v", vCurImageNodeError, vCurImage)
		}
		
		vRis[vCurImage.ID] = vCurImageNode
	}
	return vRis,nil
}

