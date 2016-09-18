# xlU_go

Content-keyed directory structures for
[xlattice_go;](https://jddixon.github.io/xlattice_go)
a Go implementation of a file system
based on **content keys**, where a content key is a cryptographically
secure hash either 20 or 32 bytes long.

## Directory Structure

Files are stored under two directory levels.
At the upper level, directory names are the hex value of the first
8 bits of the content key, so there are 256 such directories.
Below each is a set of directories whose names correspond to the
second byte of the content key.  Files within the lower directories
have names based on the hex value of the rest of the content key.
Conventionally the upper directory is called `U`.

So if a file has a content key like `0xabcdef01234...` then it will
usually be found at `U/ab/cd/ef01234...`

There will also be at least two subdirectories at the upper level

* `U/in`, used for assembling files prior to being moved into the
  appropriate lower directory; typically used for files being
  received over the wire
* `U/tmp`, similarly used by the system for temporary storage of files
  being moved into `U/` from the local file system

## Content Keys

Hashes are calculated using
a variant of the
[Secure Hash Algorithm](https://en.wikipedia.org/wiki/Secure_Hash_Algorithm).
Currently three versions of SHA are supported:

* **SHA1** (20 bytes/160 bits), also called SHA-1
* **SHA2** (32 bytes/256 bits), also called SHA-256
* **SHA3** (32 bytes/256 bits), which is the 256-bit variant of **Keccak**,
  the winner of a recent competition run by an agency of the US government

## API

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

	GetDirStruc() DirStruc
	GetPath() string

	// These exist in two forms, for convenience
	HexKeyExists(key string) (bool, error)
	ByteKeyExists(key []byte) (bool, error)

	HexKeyFileLen(key string) (length int64, err error)
	ByteKeyFileLen(key []byte) (length int64, err error)

	GetPathForHexKey(key string) (string, error)
	GetPathForByteKey(key []byte) (string, error)

## Implementation

The actual Go code is in `u.go`.  Test code is in three files:
`hash_test.go` contains common tests, and `sha1_test.go` and
`sha3_test.go` contain driver code specifc to SHA1 and SHA3 respectively.

`myData`, `myU1`, `myU2`, and `myU3` are scratch directories used in testing.
The file named `abc` contains just that string and is used in an
extremely simple check that SHA3 is functionng correctly.
The Keccak distribution (`KeccakKat-3.zip`) contains a very large collection
of test vectors.

## Project Status

A good beta.  The code is stable and well-tested.

## On-line Documentation

More information on the **xlU_go** project can be found
[here](https://jddixon.github.io/xlU_go)
