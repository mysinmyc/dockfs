package utils

import ( 
        "golang.org/x/net/context"
        "bazil.org/fuse"
	"github.com/mysinmyc/gocommons/diagnostic"
)

type ReadFunc func()([]byte,error)  

type DynamicFileNode struct {
	readFunc ReadFunc 
	content []byte
}


func NewDynamicFileNode(pReadFunc ReadFunc)(*DynamicFileNode, error) {
	return &DynamicFileNode{readFunc:pReadFunc},nil
}

func (vSelf *DynamicFileNode) Attr(pContext context.Context, pAttr *fuse.Attr) error {

	_,vError:=vSelf.GetContent()
	if vError != nil {
		return diagnostic.NewError("Error getting size",vError)
	}
	pAttr.Size = uint64(len(vSelf.content))
	pAttr.Mode = 0444
	return nil
}


func (vSelf *DynamicFileNode) GetContent() ([]byte, error) {
	if vSelf.content !=nil {
		return vSelf.content,nil
	}
	vContent,vError:=vSelf.readFunc()
	if vError != nil {
		return nil,diagnostic.NewError("Error getting content",vError)
	}
	vSelf.content = vContent
	return vContent,nil
}

func (vSelf *DynamicFileNode) ReadAll(pContext context.Context) ([]byte, error) {
	return vSelf.GetContent()
}

func (vSelf *DynamicFileNode) GetDirentType() (fuse.DirentType) {
	return fuse.DT_File 
}
