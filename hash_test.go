package xlU_go

import (
	"encoding/hex"
	"fmt"
	xu "github.com/jddixon/xlUtil_go"
	xf "github.com/jddixon/xlUtil_go/lfs"
	. "gopkg.in/check.v1"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	dataPath string
	uPath    string
	uInDir   string
	uTmpDir  string
	whichSHA int
)

// SETUP AND TEARDOWN ///////////////////////////////////////////////
func (s *XLSuite) setUpHashTest() {
	var err error
	found, err := xf.PathExists(dataPath)
	if !found {
		// MODE SUSPECT
		if err = os.MkdirAll(dataPath, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", dataPath, err)
		}
	}
	found, err = xf.PathExists(uPath)
	if !found {
		// MODE SUSPECT
		if err = os.MkdirAll(uPath, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uPath, err)
		}
	}
	found, err = xf.PathExists(uInDir)
	if !found {
		// MODE SUSPECT
		if err = os.MkdirAll(uInDir, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uInDir, err)
		}
	}
	found, err = xf.PathExists(uTmpDir)
	if !found {

		// MODE SUSPECT
		if err = os.MkdirAll(uTmpDir, 0775); err != nil {
			fmt.Printf("error creating %s: %v\n", uTmpDir, err)
		}
	}
}
func (s *XLSuite) teardownHashTest() {
	// arguably should remove the two directories
}

// UNIT TESTS ///////////////////////////////////////////////////////
func (s *XLSuite) doTestCopyAndPut(
	c *C, u UI, digest hash.Hash) {
	//we are testing uLen, uKey, err = u.CopyAndPut3(path, key)

	// create a random file
	rng := u.GetRNG()
	dLen, dPath := rng.NextDataFile(dataPath, 16*1024, 1) //  maxLen, minLen
	var dKey string
	var err error
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil) // actual, Equals, expected

	// invoke function
	var uLen int64
	var uKey string
	switch whichSHA {
	case xu.USING_SHA1:
		uLen, uKey, err = u.CopyAndPut1(dPath, dKey)
	case xu.USING_SHA2:
		uLen, uKey, err = u.CopyAndPut2(dPath, dKey)
	case xu.USING_SHA3:
		uLen, uKey, err = u.CopyAndPut3(dPath, dKey)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	c.Assert(dKey, Equals, uKey)

	// verify that original and copy both exist
	found, err := xf.PathExists(dPath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	xPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, IsNil)
	found, err = xf.PathExists(xPath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	// HACK - SIMPLEST Keccak TEST VECTOR
	if whichSHA == xu.USING_SHA3 {
		dKey, err = FileSHA3("abc")
		fmt.Printf("SHA3-256 for 'abc' is %s\n", dKey)
	}
	// END HACK
}
func (s *XLSuite) doTestExists(c *C, u UI, digest hash.Hash) {
	//we are testing whether = u.Exists( key) and whether = u.KeyExists()

	rng := u.GetRNG()
	dLen, dPath := rng.NextDataFile(dataPath, 16*1024, 1)
	var dKey string
	var err error
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	var uLen int64
	var uKey string
	switch whichSHA {
	case xu.USING_SHA1:
		uLen, uKey, err = u.CopyAndPut1(dPath, dKey)
	case xu.USING_SHA2:
		uLen, uKey, err = u.CopyAndPut2(dPath, dKey)
	case xu.USING_SHA3:
		uLen, uKey, err = u.CopyAndPut3(dPath, dKey)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	kPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, Equals, nil)
	found, err := xf.PathExists(kPath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	bKey, err := hex.DecodeString(uKey)
	c.Assert(err, IsNil)

	found, err = u.HexKeyExists(uKey) // string version of key
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
	found, err = u.ByteKeyExists(bKey) // binary version of key
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	os.Remove(kPath)

	found, err = xf.PathExists(kPath) // string version
	c.Assert(err, IsNil)
	c.Assert(found, Equals, false)

	found, err = u.ByteKeyExists(bKey) // binary version of key
	c.Assert(err, IsNil)
	c.Assert(found, Equals, false)
}
func (s *XLSuite) doTestFileLen(c *C, u UI, digest hash.Hash) {
	//we are testing len = u.fileLen(key)

	rng := u.GetRNG()
	dLen, dPath := rng.NextDataFile(dataPath, 16*1024, 1)
	var dKey string
	var err error
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	var uLen int64
	var uKey string
	switch whichSHA {
	case xu.USING_SHA1:
		uLen, uKey, err = u.CopyAndPut1(dPath, dKey)
	case xu.USING_SHA2:
		uLen, uKey, err = u.CopyAndPut2(dPath, dKey)
	case xu.USING_SHA3:
		uLen, uKey, err = u.CopyAndPut3(dPath, dKey)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	kPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, IsNil)
	_ = kPath // NOT USED
	length, err := u.HexKeyFileLen(uKey)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, length)
}

func (s *XLSuite) doTestFileHash(c *C, u UI, digest hash.Hash) {
	// we are testing shaXKey = fileSHAX(path)
	rng := u.GetRNG()
	dLen, dPath := rng.NextDataFile(dataPath, 16*1024, 1)
	data, err := ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, int64(len(data)))
	digest.Write(data)
	hash := digest.Sum(nil)
	dKey := hex.EncodeToString(hash) // 'expected'
	var fKey string
	switch len(dKey) {
	case xu.SHA1_HEX_LEN:
		fKey, err = FileSHA1(dPath) // 'actual'
	case xu.SHA2_HEX_LEN:
		fKey, err = FileSHA2(dPath) // 'actual'
		// XXX THIS CASE WON'T WORK
		// case xu.SHA3_HEX_LEN:
		//	fKey, err = FileSHA3(dPath) // 'actual'
		// DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	c.Assert(fKey, Equals, dKey)
}

func (s *XLSuite) doTestGetPathForKey(c *C, u UI, digest hash.Hash) {
	// we are testing path = GetPathForKey(key)

	rng := u.GetRNG()
	dLen, dPath := rng.NextDataFile(dataPath, 16*1024, 1)
	var err error
	var dKey, uKey string
	var uLen int64
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
		c.Assert(err, Equals, nil)
		uLen, uKey, err = u.CopyAndPut1(dPath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
		c.Assert(err, Equals, nil)
		uLen, uKey, err = u.CopyAndPut2(dPath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		c.Assert(err, Equals, nil)
		uLen, uKey, err = u.CopyAndPut3(dPath, dKey)
		c.Assert(err, Equals, nil)
		// XXX DEFAULT = ERROR
	}
	c.Assert(uLen, Equals, dLen)
	kPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, IsNil)

	var expectedPath string
	dirStruc := u.GetDirStruc()
	switch dirStruc {
	case DIR_FLAT:
		expectedPath = fmt.Sprintf("%s/%s",
			u.GetPath(), uKey)
	case DIR16x16:
		expectedPath = fmt.Sprintf("%s/%s/%s/%s",
			u.GetPath(), uKey[0:1], uKey[1:2], uKey[2:])
	case DIR256x256:
		expectedPath = fmt.Sprintf("%s/%s/%s/%s",
			u.GetPath(), uKey[0:2], uKey[2:4], uKey[4:])
	}
	c.Assert(expectedPath, Equals, kPath)
}

func (s *XLSuite) doTestPut(c *C, u UI, digest hash.Hash) {
	//we are testing (len,hash)  = put(inFile, key)

	var dLen, uLen int64
	var dPath, dKey, uKey string
	var err error
	rng := u.GetRNG()
	dLen, dPath = rng.NextDataFile(dataPath, 16*1024, 1)
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		c.Assert(err, Equals, nil)
		// XXX DEFAULT = ERROR
	}
	data, err := ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	dupePath := filepath.Join(dataPath, dKey)
	err = ioutil.WriteFile(dupePath, data, 0664)
	c.Assert(err, Equals, nil)
	switch whichSHA {
	case xu.USING_SHA1:
		uLen, uKey, err = u.Put1(dPath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA2:
		uLen, uKey, err = u.Put2(dPath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA3:
		uLen, uKey, err = u.Put3(dPath, dKey)
		c.Assert(err, Equals, nil)
		// XXX DEFAULT = ERROR
	}
	c.Assert(dLen, Equals, uLen)
	kPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, IsNil)
	_ = kPath // NOT USED

	// inFile is renamed
	found, err := xf.PathExists(dPath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, false)

	found, err = u.HexKeyExists(uKey)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)

	var dupeLen int64
	var dupeKey string
	switch whichSHA {
	case xu.USING_SHA1:
		dupeLen, dupeKey, err = u.Put1(dupePath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA2:
		dupeLen, dupeKey, err = u.Put2(dupePath, dKey)
		c.Assert(err, Equals, nil)
	case xu.USING_SHA3:
		dupeLen, dupeKey, err = u.Put3(dupePath, dKey)
		c.Assert(err, Equals, nil)
	}
	c.Assert(uLen, Equals, dupeLen)
	// dupe file is deleted'
	c.Assert(uKey, Equals, dupeKey)
	found, err = xf.PathExists(dupePath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, false)

	found, err = u.HexKeyExists(uKey)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
}

func (s *XLSuite) doTestPutData(c *C, u UI, digest hash.Hash) {
	// we are testing (len,hash)  = putData3(data, key)

	var dPath, dKey, uKey string
	var dLen, uLen int64
	var err error

	rng := u.GetRNG()
	dLen, dPath = rng.NextDataFile(dataPath, 16*1024, 1)
	switch whichSHA {
	case xu.USING_SHA1:
		dKey, err = FileSHA1(dPath)
	case xu.USING_SHA2:
		dKey, err = FileSHA2(dPath)
	case xu.USING_SHA3:
		dKey, err = FileSHA3(dPath)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	data, err := ioutil.ReadFile(dPath)
	c.Assert(err, Equals, nil)
	c.Assert(int64(len(data)), Equals, dLen)

	switch whichSHA {
	case xu.USING_SHA1:
		uLen, uKey, err = u.PutData1(data, dKey)
	case xu.USING_SHA2:
		uLen, uKey, err = u.PutData2(data, dKey)
	case xu.USING_SHA3:
		uLen, uKey, err = u.PutData3(data, dKey)
		// XXX DEFAULT = ERROR
	}
	c.Assert(err, Equals, nil)
	c.Assert(dLen, Equals, uLen)
	c.Assert(dKey, Equals, uKey)

	found, err := u.HexKeyExists(uKey)
	c.Assert(err, Equals, nil)
	c.Assert(found, Equals, true)

	xPath, err := u.GetPathForHexKey(uKey)
	c.Assert(err, IsNil)
	found, err = xf.PathExists(xPath)
	c.Assert(err, IsNil)
	c.Assert(found, Equals, true)
}
