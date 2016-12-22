package utils

import (
	"time"
	"sync"
	"strings"
	"os"
	"golang.org/x/net/context"
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"github.com/mysinmyc/gocommons/diagnostic"
)

type ChildrenPopulatorFunc func() (map[string] DirentTyped, error)

type DynamicDir struct {
	lock	    		sync.Mutex
	expiry      		time.Time
	lastUpdate     		time.Time
	CacheInterval 		time.Duration
	Children    		map[string] DirentTyped
	ChildrenPopulatorFunc   ChildrenPopulatorFunc
}

func (vSelf *DynamicDir) Attr(pContext context.Context, pAttr *fuse.Attr) error {
        pAttr.Valid = vSelf.CacheInterval
	if vSelf.lastUpdate.IsZero() == false {
        	pAttr.Ctime = vSelf.lastUpdate
	}
        pAttr.Mode = os.ModeDir | 0555
        return nil
}

func (vSelf *DynamicDir) populateChildren() error {
	if vSelf.expiry.After(time.Now()) {
		return nil
	}
	
	vSelf.lock.Lock()
	defer vSelf.lock.Unlock()

	vNow := time.Now()
	if vSelf.expiry.After(vNow) {
		return nil
	}
	if diagnostic.IsLogTrace() {
		diagnostic.LogTrace("DynamicDir.populateChildren", "Populating children...")
	}
		
	if vSelf.ChildrenPopulatorFunc != nil {
		vChildren,vChildrenError := vSelf.ChildrenPopulatorFunc()

		if vChildrenError != nil {
			return diagnostic.NewError("Failed to populate children", vChildrenError)
		}
		vSelf.Children = vChildren
	}	

	vSelf.expiry = vNow.Add(vSelf.CacheInterval)
	vSelf.lastUpdate = vNow
	return nil
}

func (vSelf *DynamicDir) Lookup(pContext context.Context, pName string) (fs.Node, error) {
	vSelf.populateChildren()
	if vSelf.Children != nil {
		vRis:=vSelf.Children[pName]

		if vRis != nil {
			return vRis.(fs.Node),nil
        	}
	}
        return nil, fuse.ENOENT
}

func (vSelf *DynamicDir) ReadDirAll(pContext context.Context) ([]fuse.Dirent, error) {
	vSelf.populateChildren()
        vRis:= make([]fuse.Dirent,0)
	if vSelf.Children != nil {
		for vCurChildName, vCurChild := range vSelf.Children {
			vRis=append(vRis, fuse.Dirent{ Type: vCurChild.GetDirentType(), Name:vCurChildName})
		}
	}
        return vRis,nil
}


func (vSelf *DynamicDir) Add(pName string, pChild DirentTyped ) error {
	if vSelf.Children ==nil {
		vSelf.Children = make(map[string] DirentTyped)
	}
	return AddDirentTo(pName,pChild,vSelf.Children)
}

func AddDirentTo(pName string, pChild DirentTyped, pParentMap map[string] DirentTyped) error {

	diagnostic.LogTrace("AddDirentTo", "adding %s to %v", pName, pParentMap)

	if strings.HasPrefix(pName,"/") {
		pName = pName[1:]
	}
	vSlashPosition := strings.Index(pName, "/")

	if vSlashPosition != -1 {
		vParentName:=pName[:vSlashPosition]
		vParent:=pParentMap[vParentName]
		if vParent == nil {
			vParent = &DynamicDir {}
			pParentMap[vParentName] = vParent
		}
		vParentTyped, vParentOk := vParent.(DirentParent)

		if vParentOk == false {
			return diagnostic.NewError("Cannot add child %s to %t",nil, pName, vParent)
		}

		return vParentTyped.Add(pName[vSlashPosition+1:], pChild)
	}

	pParentMap[pName] = pChild
	return nil
}

func (vSelf *DynamicDir) GetDirentType() (fuse.DirentType) {
        return fuse.DT_Dir
}

