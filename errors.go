package xlU_go

import "errors"

var (
	BadKeyLength          = errors.New("bad key length")
	DirStrucNotRecognized = errors.New("DirStruc not recognized")
	EmptyKey              = errors.New("empty key parameter")
	EmptyPath             = errors.New("empty path parameter")
	FileNotFound          = errors.New("file not found")
	NilKey                = errors.New("nil binary key parameter")
)
