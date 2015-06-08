package xlU_go

import (
	"code.google.com/p/go.crypto/sha3"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt" // DEBUG
	xr "github.com/jddixon/rnglib_go"
	xu "github.com/jddixon/xlUtil_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	"os"
	"path/filepath"
)

// CLASS, so to speak ///////////////////////////////////////////////

type UFlat struct {
	path   string // all parameters are private
	inDir  string
	tmpDir string
	rng    *xr.PRNG
}

// Create a new Flat file system, ORing perm into the default permissions.
// If perm is 0, the default is to allow user and group access.
// If the root is U, then this creates U/, U/tmp, and U/in.
func NewUFlat(path string, perm os.FileMode) (udir *UFlat, err error) {
	// TODO: validate path
	var (
		inDir, tmpDir string
	)
	err = os.MkdirAll(path, 0750|perm)
	if err == nil {
		inDir = filepath.Join(path, "in")
		err = os.MkdirAll(inDir, 0770|perm)
		if err == nil {
			tmpDir = filepath.Join(path, "tmp")
			err = os.MkdirAll(tmpDir, 0700)
		}
	}
	udir = &UFlat{
		path:   path,
		rng:    xr.MakeSimpleRNG(),
		inDir:  inDir,
		tmpDir: tmpDir,
	}
	return
}

func (u *UFlat) GetDirStruc() DirStruc { return DIR_FLAT }
func (u *UFlat) GetPath() string       { return u.path }
func (u *UFlat) GetRNG() *xr.PRNG      { return u.rng }

// HEX KEY FUNCTIONS ================================================

// DEPRECATED
func (u *UFlat) Exists(key string) (found bool, err error) {
	return u.HexKeyExists(key)
}

func (u *UFlat) HexKeyExists(key string) (found bool, err error) {
	path, err := u.GetPathForHexKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	return
}

// DEPRECATED
func (u *UFlat) FileLen(key string) (length int64, err error) {
	return u.HexKeyFileLen(key)
}
func (u *UFlat) HexKeyFileLen(key string) (length int64, err error) {
	path, err := u.GetPathForHexKey(key)
	if err == nil {
		var info os.FileInfo
		info, err = os.Stat(path) // ERRORS IGNORED
		if err == nil {
			length = info.Size()
		}
	}
	return
}

// DEPRECATED
func (u *UFlat) GetPathForKey(key string) (path string, err error) {
	return u.GetPathForHexKey(key)
}

// Returns a path to a file with the content key passed.
func (u *UFlat) GetPathForHexKey(key string) (path string, err error) {
	if key == "" {
		err = EmptyKey
	} else {
		path = filepath.Join(u.path, key)
	}
	return
}

// BINARY KEY FUNCTIONS =============================================

// DEPRECATED
func (u *UFlat) KeyExists(key []byte) (found bool, err error) {
	return u.ByteKeyExists(key)
}

func (u *UFlat) ByteKeyExists(key []byte) (found bool, err error) {
	path, err := u.GetPathForByteKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	return
}

// DEPRECATED
func (u *UFlat) KeyFileLen(key []byte) (length int64, err error) {
	return u.ByteKeyFileLen(key)
}

func (u *UFlat) ByteKeyFileLen(key []byte) (length int64, err error) {
	path, err := u.GetPathForByteKey(key)
	if err == nil {
		var info os.FileInfo
		info, err = os.Stat(path)
		if err == nil {
			length = info.Size()
		}
	}
	return
}

// DEPRECATED
func (u *UFlat) GetPathForBinaryKey(key []byte) (path string, err error) {
	return u.GetPathForByteKey(key)
}

// Returns a path to a file with the content key passed.
//
func (u *UFlat) GetPathForByteKey(key []byte) (path string, err error) {
	if key == nil {
		err = NilKey
	} else {
		strKey := hex.EncodeToString(key)
		path, err = u.GetPathForHexKey(strKey)
	}
	return
}

// SHA1/SHA3 NEUTRAL FUNCTIONS ======================================

// - CopyAndPut -----------------------------------------------------

// Copy the file at path to a randomly-named temporary file under U/tmp.
// If that operation succeeds, we then attempt to rename the file into
// the appropriate U data subdirectory.  If the file is already present,
// we silently discard the copy.  Returns the length of the file in bytes
// and any error.
//
// CopyAndPut1 and 3 return the actual content hash; this doesn't.
//
func (u *UFlat) CopyAndPut(pathToFile string, key []byte) (
	length int64, err error) {

	if pathToFile == "" {
		err = EmptyPath
	} else if key == nil {
		err = NilKey
	} else {
		strKey := hex.EncodeToString(key)
		switch len(strKey) {
		case xu.SHA1_HEX_LEN:
			length, _, err = u.CopyAndPut1(pathToFile, strKey)
		case xu.SHA3_HEX_LEN:
			length, _, err = u.CopyAndPut3(pathToFile, strKey)
		default:
			err = BadKeyLength
		}
	}
	return
}

// - GetData --------------------------------------------------------

// Retrieves file contents using a binary key.  The key is the SHA1
// or SHA3 hash of the file contents.
//
func (u *UFlat) GetData(key []byte) (data []byte, err error) {
	if key == nil {
		err = NilKey
	} else {
		strKey := hex.EncodeToString(key)
		switch len(strKey) {
		case xu.SHA1_HEX_LEN:
			data, err = u.GetData1(strKey)
		case xu.SHA3_HEX_LEN:
			data, err = u.GetData3(strKey)
		default:
			err = BadKeyLength
		}
	}
	return
}

// - Put ------------------------------------------------------------

// Given a local temporary file, either rename it into U or just silently
// delete it if the data is already present in U.  Returns the length
// of the file and any error.
//
// Put1 and 3 return the actual content hash; this doesn't.
//
func (u *UFlat) Put(tmpFile string, key []byte) (length int64, err error) {

	if tmpFile == "" {
		err = EmptyPath
	} else if key == nil {
		err = NilKey
	} else {
		strKey := hex.EncodeToString(key)
		switch len(strKey) {
		case xu.SHA1_HEX_LEN:
			length, _, err = u.Put1(tmpFile, strKey)
		case xu.SHA3_HEX_LEN:
			length, _, err = u.Put3(tmpFile, strKey)
		default:
			err = BadKeyLength
		}
	}
	return
}

// - PutData --------------------------------------------------------

// Write data into the store using a binary key.
//
func (u *UFlat) PutData(data []byte, key []byte) (
	length int64, hash []byte, err error) {

	if key == nil {
		err = NilKey
	} else {
		var strHash string
		strKey := hex.EncodeToString(key)
		switch len(strKey) {
		case xu.SHA1_HEX_LEN:
			length, strHash, err = u.PutData1(data, strKey)
		case xu.SHA3_HEX_LEN:
			length, strHash, err = u.PutData3(data, strKey)
		default:
			err = BadKeyLength
		}
		if err == nil {
			hash, err = hex.DecodeString(strHash)
		}
	}
	return
}

// SHA1 CODE ========================================================

// CopyAndPut1 ------------------------------------------------------

func (u *UFlat) CopyAndPut1(path, key string) (
	written int64, hash string, err error) {

	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	found, err := xf.PathExists(tmpFileName)
	for found {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
		found, err = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u.Put1(tmpFileName, key)
	}
	return
}

// - GetData1 --------------------------------------------------------
func (u *UFlat) GetData1(key string) (data []byte, err error) {

	var (
		found bool
		path  string
		src   *os.File
	)
	path, err = u.GetPathForHexKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		src, err = os.Open(path)
	}
	if err == nil {
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 128 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
	}
	return
}

// - Put1 ------------------------------------------------------------

// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized _FLAT
// key is an sha1 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// XXX We don't do much checking.
//
func (u *UFlat) Put1(inFile, key string) (
	length int64, hash string, err error) {

	var (
		found       bool
		fullishPath string
	)
	hash, err = FileHexSHA1(inFile)
	if err != nil {
		fmt.Printf("DEBUG: FileHexSHA1 returned error %v\n", err)
		return
	}
	if hash != key {
		fmt.Printf("expected %s to have key %s, but the content key is %s\n",
			inFile, key, hash)
		err = errors.New("IllegalArgument: Put1: key does not match content")
		return
	}
	info, err := os.Stat(inFile)
	if err == nil {
		length = info.Size()
		fullishPath = filepath.Join(u.path, key)
		found, err = xf.PathExists(fullishPath)
	}
	if err == nil {
		if found {
			// drop the temporary input file
			err = os.Remove(inFile)
		} else {
			// rename the temporary file into U
			err = os.Rename(inFile, fullishPath)
			if err == nil {
				err = os.Chmod(fullishPath, 0444)
			}
		}
	}
	return
}

// PutData1 ---------------------------------------------------------
func (u *UFlat) PutData1(data []byte, key string) (
	length int64, hash string, err error) {

	var fullishPath string
	var found bool

	s := sha1.New()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length = int64(len(data))

	if err == nil {
		fullishPath = filepath.Join(u.path, key)
		found, err = xf.PathExists(fullishPath)
		if err == nil && !found {
			var dest *os.File
			dest, err = os.Create(fullishPath)
			if err == nil {
				var count int
				defer dest.Close()
				count, err = dest.Write(data)
				if err == nil {
					length = int64(count)
				}
			}
		}
	}
	return
}

// SHA3 CODE ========================================================

//- CopyAndPut3 -----------------------------------------------------

// Copy the file at path to a randomly-named temporary file under U/tmp.
// If that operation succeeds, we then attempt to rename the file into
// the appropriate U data subdirectory.  If the file is already present,
// we silently discard the copy.  Returns the length of the file in bytes,
// its actual content hash, and any error.
//
func (u *UFlat) CopyAndPut3(path, key string) (
	written int64, hash string, err error) {

	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	found, _ := xf.PathExists(tmpFileName) // XXX error ignored
	for found {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
		found, _ = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u.Put3(tmpFileName, key)
	}
	return
}

// - GetData3 --------------------------------------------------------

func (u *UFlat) GetData3(key string) (data []byte, err error) {
	var (
		found bool
		path  string
	)
	path, err = u.GetPathForHexKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		var src *os.File
		if src, err = os.Open(path); err != nil {
			return
		}
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 128 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
	}
	return
}

// - Put3 ------------------------------------------------------------

// inFile is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized _FLAT
// key is an sha3 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// We don't do much checking.
//
func (u *UFlat) Put3(inFile, key string) (
	length int64, hash string, err error) {

	var fullishPath string

	hash, err = FileHexSHA3(inFile)
	if err != nil {
		fmt.Printf("DEBUG: FileHexSHA3 returned error %v\n", err)
		return
	}
	if hash != key {
		fmt.Printf("expected %s to have key %s, but the content key is %s\n",
			inFile, key, hash)
		err = errors.New("IllegalArgument: Put3: key does not match content")
		return
	}
	info, err := os.Stat(inFile)
	if err != nil {
		return
	}
	length = info.Size()

	if err == nil {
		var found bool

		fullishPath = filepath.Join(u.path, key)
		found, err = xf.PathExists(fullishPath)
		if err == nil {
			if found {
				// drop the temporary input file
				err = os.Remove(inFile)
			} else {
				// rename the temporary file into U
				err = os.Rename(inFile, fullishPath)
			}
		}
	}
	if err == nil {
		err = os.Chmod(fullishPath, 0444)
	}
	return
}

// - PutData3 --------------------------------------------------------

func (u *UFlat) PutData3(data []byte, key string) (length int64, hash string, err error) {
	s := sha3.NewKeccak256()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length = int64(len(data))
	fullishPath := filepath.Join(u.path, key)
	found, err := xf.PathExists(fullishPath)
	if !found {
		var dest *os.File
		dest, err = os.Create(fullishPath)
		if err == nil {
			var count int
			defer dest.Close()
			count, err = dest.Write(data)
			if err == nil {
				length = int64(count)
			}
		}
	}
	return
}

// SHA2 CODE ========================================================

// CopyAndPut2 ------------------------------------------------------

func (u *UFlat) CopyAndPut2(path, key string) (
	written int64, hash string, err error) {

	// the temporary file MUST be created on the same device
	// xxx POSSIBLE RACE CONDITION
	tmpFileName := filepath.Join(u.tmpDir, u.rng.NextFileName(16))
	found, err := xf.PathExists(tmpFileName)
	for found {
		tmpFileName = filepath.Join(u.tmpDir, u.rng.NextFileName(16))
		found, err = xf.PathExists(tmpFileName)
	}
	written, err = CopyFile(tmpFileName, path) // dest <== src
	if err == nil {
		written, hash, err = u.Put2(tmpFileName, key)
	}
	return
}

// - GetData2 --------------------------------------------------------
func (u *UFlat) GetData2(key string) (data []byte, err error) {

	var (
		found bool
		path  string
		src   *os.File
	)
	path, err = u.GetPathForHexKey(key)
	if err == nil {
		found, err = xf.PathExists(path)
	}
	if err == nil && !found {
		err = FileNotFound
	}
	if err == nil {
		src, err = os.Open(path)
	}
	if err == nil {
		defer src.Close()
		var count int
		// XXX THIS WILL NOT WORK FOR LARGER FILES!  It will ignore
		//     anything over 128 KB
		data = make([]byte, DEFAULT_BUFFER_SIZE)
		count, err = src.Read(data)
		// XXX COUNT IS IGNORED
		_ = count
	}
	return
}

// - Put2 ------------------------------------------------------------

// tmp is the path to a local file which will be renamed into U (or deleted
// if it is already present in U)
// u.path is an absolute or relative path to a U directory organized _FLAT
// key is an sha256 content hash.
// If the operation succeeds we return the length of the file (which must
// not be zero.  Otherwise we return 0.
// XXX We don't do much checking.
//
func (u *UFlat) Put2(inFile, key string) (
	length int64, hash string, err error) {

	var (
		found       bool
		fullishPath string
	)
	hash, err = FileHexSHA2(inFile)
	if err != nil {
		fmt.Printf("DEBUG: FileHexSHA2 returned error %v\n", err)
		return
	}
	if hash != key {
		fmt.Printf("expected %s to have key %s, but the content key is %s\n",
			inFile, key, hash)
		err = errors.New("IllegalArgument: Put2: key does not match content")
		return
	}
	info, err := os.Stat(inFile)
	if err == nil {
		length = info.Size()
		fullishPath = filepath.Join(u.path, key)
		found, err = xf.PathExists(fullishPath)
	}
	if err == nil {
		if found {
			// drop the temporary input file
			err = os.Remove(inFile)
		} else {
			// rename the temporary file into U
			err = os.Rename(inFile, fullishPath)
			if err == nil {
				err = os.Chmod(fullishPath, 0444)
			}
		}
	}
	return
}

// PutData2 ---------------------------------------------------------
func (u *UFlat) PutData2(data []byte, key string) (
	length int64, hash string, err error) {

	var fullishPath string
	var found bool

	s := sha256.New()
	s.Write(data)
	hash = hex.EncodeToString(s.Sum(nil))
	if hash != key {
		fmt.Printf("expected data to have key %s, but content key is %s",
			key, hash)
		err = errors.New("content/key mismatch")
		return
	}
	length = int64(len(data))

	if err == nil {
		fullishPath = filepath.Join(u.path, key)
		found, err = xf.PathExists(fullishPath)
		if err == nil && !found {
			var dest *os.File
			dest, err = os.Create(fullishPath)
			if err == nil {
				var count int
				defer dest.Close()
				count, err = dest.Write(data)
				if err == nil {
					length = int64(count)
				}
			}
		}
	}
	return
}
