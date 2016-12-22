package image


import (
	"time"
	"strings"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/gocommons/mystrings"
	"github.com/mysinmyc/dockfs/utils"
)


type ImagesByTagNode struct {
	utils.DynamicDir
	dockerClient  *client.Client
}


func NewImagesByTagNode( pDockerClient *client.Client) (*ImagesByTagNode, error) {
	vRis:=&ImagesByTagNode {dockerClient: pDockerClient}
       	vRis.DynamicDir.CacheInterval= time.Second*5
        vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ImagesByTagNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vImages,vError:=vSelf.dockerClient.ImageList(context.Background(),types.ImageListOptions{All:true})
	if vError!=nil {
		return nil,diagnostic.NewError("failed to list images",vError)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vImages))
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
			vAddError:=utils.AddDirentTo(vCurImagePath, vCurSymLink, vRis)
			if vAddError!=nil {
				return nil, diagnostic.NewError("failed to add child ",vAddError)
			}

		}	
	}
	return vRis,nil
}

