package xlU_go

import (
	xu "github.com/jddixon/xlUtil_go"
	"golang.org/x/crypto/sha3"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) setUp3() {
	dataPath = "myData"
	uPath = "myU3"
	uInDir = "myU3/in"
	uTmpDir = "myU3/tmp"
	s.setUpHashTest()
	whichSHA = xu.USING_SHA3
}

func (s *XLSuite) doTestCopyAndPut3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestCopyAndPut(c, myU, sha3.New256())
}
func (s *XLSuite) TestCopyAndPut3(c *C) {
	s.setUp3()
	s.doTestCopyAndPut3(c, DIR_FLAT)
	s.doTestCopyAndPut3(c, DIR16x16)
	s.doTestCopyAndPut3(c, DIR256x256)
}

func (s *XLSuite) doTestExists3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestExists(c, myU, sha3.New256())
}
func (s *XLSuite) TestExists3(c *C) {
	s.setUp3()
	s.doTestExists3(c, DIR_FLAT)
	s.doTestExists3(c, DIR16x16)
	s.doTestExists3(c, DIR256x256)
}

func (s *XLSuite) doTestFileLen3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileLen(c, myU, sha3.New256())
}
func (s *XLSuite) TestFileLen3(c *C) {
	s.setUp3()
	s.doTestFileLen3(c, DIR_FLAT)
	s.doTestFileLen3(c, DIR16x16)
	s.doTestFileLen3(c, DIR256x256)
}

func (s *XLSuite) doTestFileHash3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileHash(c, myU, sha3.New256())
}

// XXX FAILS BECAUSE CAN'T DISTINGUISH HASH TYPE BY KEY LENGTH
//func (s *XLSuite) TestFileHash3(c *C) {
//	s.setUp3()
//	s.doTestFileHash3(c, DIR_FLAT)
//	s.doTestFileHash3(c, DIR16x16)
//	s.doTestFileHash3(c, DIR256x256)
//}

func (s *XLSuite) doTestGetPathForKey3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestGetPathForKey(c, myU, sha3.New256())
}
func (s *XLSuite) TestGetPathForKey3(c *C) {
	s.setUp3()
	s.doTestGetPathForKey3(c, DIR_FLAT)
	s.doTestGetPathForKey3(c, DIR16x16)
	s.doTestGetPathForKey3(c, DIR256x256)
}

func (s *XLSuite) doTestPut3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPut(c, myU, sha3.New256())
}
func (s *XLSuite) TestPut3(c *C) {
	s.setUp3()
	s.doTestPut3(c, DIR_FLAT)
	s.doTestPut3(c, DIR16x16)
	s.doTestPut3(c, DIR256x256)
}

func (s *XLSuite) doTestPutData3(c *C, ds DirStruc) {
	s.setUp3()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPutData(c, myU, sha3.New256())
}
func (s *XLSuite) TestPutData3(c *C) {
	s.setUp3()
	s.doTestPutData3(c, DIR_FLAT)
	s.doTestPutData3(c, DIR16x16)
	s.doTestPutData3(c, DIR256x256)
}
