package xlU_go

import (
	"bytes"
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	xu "github.com/jddixon/xlUtil_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	"io"
	"io/ioutil"
	"os"
)

func New(path string, ds DirStruc, perm os.FileMode) (UI, error) {

	switch ds {
	case DIR_FLAT:
		return NewUFlat(path, perm)
	case DIR16x16:
		return NewU16x16(path, perm)
	case DIR256x256:
		return NewU256x256(path, perm)
	default:
		return nil, DirStrucNotRecognized
	}
}

// PACKAGE-LEVEL FUNCTIONS //////////////////////////////////////////

func CopyFile(destName, srcName string) (written int64, err error) {
	var (
		src, dest *os.File
	)
	if src, err = os.Open(srcName); err != nil {
		return
	}
	defer src.Close()
	if dest, err = os.Create(destName); err != nil {
		return
	}
	defer dest.Close()
	return io.Copy(dest, src) // returns written, err
}

// - FileHexSHA1 --------------------------------------------------------

// XXX DEPRECATED
func FileSHA1(path string) (string, error) {
	return FileHexSHA1(path)
}
// returns the SHA1 binHash of the contents of a file
func FileBinSHA1(path string) (binHash []byte, err error) {
	var data2 []byte
	binHash = xu.SHA1_BIN_NONE
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = errors.New("IllegalArgument: empty path or non-existent file")
	}
	if err == nil {
		data2, err = ioutil.ReadFile(path)
	}
	if err == nil {
		d2 := sha1.New()
		d2.Write(data2)
		binHash = d2.Sum(nil)
	}
	return
}
func FileHexSHA1(path string) (hash string, err error) {
	binHash, err := FileBinSHA1(path)
	if err == nil {
		if bytes.Equal(binHash, xu.SHA1_BIN_NONE) {
			hash = xu.SHA1_HEX_NONE
		} else {
			hash = hex.EncodeToString(binHash)
		}
	}
	return
}
// - FileHexSHA2 --------------------------------------------------------

// XXX DEPRECATED
func FileSHA2(path string) (string, error) {
	return FileHexSHA2(path)
}

// returns the SHA2 hash of the contents of a file
func FileBinSHA2(path string) (binHash []byte, err error) {
	var data2 []byte

	binHash = xu.SHA2_BIN_NONE
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = errors.New("IllegalArgument: empty path or non-existent file")
	}

	if err == nil {
		data2, err = ioutil.ReadFile(path)
	}
	if err == nil {
		d2 := sha256.New()
		d2.Write(data2)
		binHash = d2.Sum(nil)
	}
	return
}
func FileHexSHA2(path string) (hexHash string, err error) {
	binHash, err := FileBinSHA2(path)
	if err == nil {
		hexHash = hex.EncodeToString(binHash)
	} 
	return
}
// - FileHexSHA3 --------------------------------------------------------

// XXX DEPRECATED
func FileSHA3(path string) (string, error) {
	return FileHexSHA3(path)
}
// returns the SHA3 hash of the contents of a file
func FileBinSHA3(path string) (binHash []byte, err error) {
	var data2 []byte

	binHash = xu.SHA3_BIN_NONE
	found, err := xf.PathExists(path)
	if err == nil && !found {
		err = errors.New("IllegalArgument: empty path or non-existent file")
	}

	if err == nil {
		data2, err = ioutil.ReadFile(path)
	}
	if err == nil {
		d2 := sha3.NewKeccak256()
		d2.Write(data2)
		binHash = d2.Sum(nil)
	}
	return
}
func FileHexSHA3(path string) (hexHash string, err error) {
	binHash, err := FileBinSHA3(path)
	if err == nil {
		hexHash = hex.EncodeToString(binHash)
	}
	return
}
