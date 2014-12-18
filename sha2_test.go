package xlU_go

import (
	"crypto/sha256"
	xu "github.com/jddixon/xlUtil_go"
	. "gopkg.in/check.v1"
)

func (s *XLSuite) setUp2() {
	dataPath = "myData"
	uPath = "myU2"
	uInDir = "myU2/in"
	uTmpDir = "myU2/tmp"
	s.setUpHashTest()
	whichSHA = xu.USING_SHA2
}

func (s *XLSuite) doTestCopyAndPut2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestCopyAndPut(c, myU, sha256.New())
}
func (s *XLSuite) TestCopyAndPut2(c *C) {
	s.setUp2()
	s.doTestCopyAndPut2(c, DIR_FLAT)
	s.doTestCopyAndPut2(c, DIR16x16)
	s.doTestCopyAndPut2(c, DIR256x256)
}

func (s *XLSuite) doTestExists2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestExists(c, myU, sha256.New())
}
func (s *XLSuite) TestExists2(c *C) {
	s.setUp2()
	s.doTestExists2(c, DIR_FLAT)
	s.doTestExists2(c, DIR16x16)
	s.doTestExists2(c, DIR256x256)
}

func (s *XLSuite) doTestFileLen2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileLen(c, myU, sha256.New())
}
func (s *XLSuite) TestFileLen2(c *C) {
	s.setUp2()
	s.doTestFileLen2(c, DIR_FLAT)
	s.doTestFileLen2(c, DIR16x16)
	s.doTestFileLen2(c, DIR256x256)
}

func (s *XLSuite) doTestFileHash2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestFileHash(c, myU, sha256.New())
}
func (s *XLSuite) TestFileHash2(c *C) {
	s.setUp2()
	s.doTestFileHash2(c, DIR_FLAT)
	s.doTestFileHash2(c, DIR16x16)
	s.doTestFileHash2(c, DIR256x256)
}

func (s *XLSuite) doTestGetPathForKey2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestGetPathForKey(c, myU, sha256.New())
}
func (s *XLSuite) TestGetPathForKey2(c *C) {
	s.setUp2()
	s.doTestGetPathForKey2(c, DIR_FLAT)
	s.doTestGetPathForKey2(c, DIR16x16)
	s.doTestGetPathForKey2(c, DIR256x256)
}

func (s *XLSuite) doTestPut2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPut(c, myU, sha256.New())
}
func (s *XLSuite) TestPut2(c *C) {
	s.setUp2()
	s.doTestPut2(c, DIR_FLAT)
	s.doTestPut2(c, DIR16x16)
	s.doTestPut2(c, DIR256x256)
}

func (s *XLSuite) doTestPutData2(c *C, ds DirStruc) {
	s.setUp2()
	myU, err := New(uPath, ds, 0)
	c.Assert(err, IsNil)
	s.doTestPutData(c, myU, sha256.New())
}
func (s *XLSuite) TestPutData2(c *C) {
	s.setUp2()
	s.doTestPutData2(c, DIR_FLAT)
	s.doTestPutData2(c, DIR16x16)
	s.doTestPutData2(c, DIR256x256)
}
