package utils

import  (
        "bazil.org/fuse"
)


type DirentTyped interface {
	GetDirentType() fuse.DirentType
}
