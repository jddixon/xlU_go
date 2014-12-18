package xlU_go

import (
	xr "github.com/jddixon/rnglib_go"
)

type UI interface {
	CopyAndPut(path string, key []byte) (int64, error)
	CopyAndPut1(path, key string) (int64, string, error)
	CopyAndPut2(path, key string) (int64, string, error)
	CopyAndPut3(path, key string) (int64, string, error)

	GetData(key []byte) ([]byte, error)
	GetData1(key string) (data []byte, err error)
	GetData2(key string) (data []byte, err error)
	GetData3(key string) (data []byte, err error)

	Put(inFile string, key []byte) (length int64, err error)
	Put1(inFile, key string) (length int64, hash string, err error)
	Put2(inFile, key string) (length int64, hash string, err error)
	Put3(inFile, key string) (length int64, hash string, err error)

	PutData(data []byte, key []byte) (length int64, hash []byte, err error)
	PutData1(data []byte, key string) (length int64, hash string, err error)
	PutData2(data []byte, key string) (length int64, hash string, err error)
	PutData3(data []byte, key string) (length int64, hash string, err error)

	// presumably temporary
	GetRNG() *xr.PRNG

	GetDirStruc() DirStruc
	GetPath() string // utility?

	// These exist in two forms, for convenience
	HexKeyExists(key string) (bool, error)
	ByteKeyExists(key []byte) (bool, error)

	HexKeyFileLen(key string) (length int64, err error)
	ByteKeyFileLen(key []byte) (length int64, err error)

	GetPathForHexKey(key string) (string, error)
	GetPathForByteKey(key []byte) (string, error)

	// DEPRECATED
	//Exists(key string) (bool, error)
	//KeyExists(key []byte) (bool, error)
	//FileLen(key string) (length int64, err error)
	//KeyFileLen(key []byte) (length int64, err error)
	//GetPathForBinaryKey(key []byte) (string, error)
	//GetPathForKey(key string) (string, error)
}
