package container


import  (
	"golang.org/x/net/context"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mysinmyc/gocommons/diagnostic"
	"github.com/mysinmyc/dockfs/utils"
	"github.com/docker/docker/api/types/filters"
)


type ContainersOfImageNode struct {
	utils.DynamicDir	
	dockerClient  		*client.Client
	imageId			string
}


func NewContainersOfImageNode(pImageId string, pDockerClient *client.Client) (*ContainersOfImageNode, error) {
	vRis:=&ContainersOfImageNode{imageId:pImageId, dockerClient: pDockerClient}
	vRis.DynamicDir.CacheInterval= time.Second*5
	vRis.DynamicDir.ChildrenPopulatorFunc= vRis.populateChildren
	return vRis,nil
}


func (vSelf *ContainersOfImageNode) populateChildren() (map[string] utils.DirentTyped,error) {

	vFilters:=filters.NewArgs()
	vFilters.Add("ancestor", vSelf.imageId)
	vContainers,vError:=vSelf.dockerClient.ContainerList(context.Background(),types.ContainerListOptions{All:true, Filters:vFilters, Limit:-1})
	if vError!=nil {
		return nil, diagnostic.NewError("failed to list containers",vError)
	}

	vRis:= make(map[string]utils.DirentTyped,len(vContainers))
	for _,vCurContainer := range vContainers {

		vSymLink,_:=utils.NewSymLinkNode("../../../../containers/byId/"+vCurContainer.ID)
                vRis[vCurContainer.ID] = vSymLink

	}
	return vRis,nil
}

