package image


import (
	"strings"
	"golang.org/x/net/context"
	"bazil.org/fuse/fs"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/gocommons/mystrings"
	"github.com/mysinmyc/dockfs/utils"
)


type ImagesByTagNode struct {
	fs.Tree
	dockerClient  *client.Client
}


func NewImagesByTagNode( pDockerClient *client.Client) (*ImagesByTagNode, error) {
	vRis:=&ImagesByTagNode {dockerClient: pDockerClient}

	vInitError:=vRis.init(context.Background())

	if vInitError != nil {
		return nil,diagnostic.NewError("initialization failed",vInitError)
	}		
	return vRis,nil
}


func (vSelf *ImagesByTagNode) init(pContext context.Context) (error) {

	vImages,vError:=vSelf.dockerClient.ImageList(pContext,types.ImageListOptions{All:true})
	if vError!=nil {
		return diagnostic.NewError("failed to list images",vError)
	}

	for _,vCurImage := range vImages {
		for _, vCurImageTag := range vCurImage.RepoTags {
			if vCurImageTag == "<none>:<none>" {
				continue
			}
			vAdditionalLevel := ""
			if strings.Index(vCurImageTag,"/") > -1 {
				vAdditionalLevel="../"
			}
			vCurSymLink,_:=utils.NewSymLinkNode(vAdditionalLevel+"../../byId/"+vCurImage.ID)
			vCurImagePath:= mystrings.ReplaceAt(vCurImageTag,strings.LastIndex(vCurImageTag,":"),"/")
			vSelf.Add(vCurImagePath, vCurSymLink)

		}	
	}
	return nil
}

