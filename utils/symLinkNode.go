package utils

import ( 
	"os"
        "golang.org/x/net/context"
        "bazil.org/fuse"
)


type SymLinkNode struct {
	target string
}


func NewSymLinkNode(pTarget string)(*SymLinkNode, error) {
	return &SymLinkNode{target:pTarget},nil
}

func (vSelf *SymLinkNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {

	pAttr.Size = uint64(len(vSelf.target))
	pAttr.Mode = os.ModeSymlink | 0555
	return nil
}

func (vSelf *SymLinkNode) Readlink(pContext context.Context, pRequest *fuse.ReadlinkRequest) (string, error) {
	return vSelf.target,nil
}

func (vSelf *SymLinkNode) GetDirentType() (fuse.DirentType) {
	return fuse.DT_Link 
}
