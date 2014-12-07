# xlU_go

Content-keyed directory structures for xlattice_go.

2013-07-06

The code in this directory is a Go implementation of a file system
based on content keys.  Files are stored under two directory levels.
At the upper level, directory names are the hex value of the first
8 bits of the content key, so there are 256 such directories.  
Below each is a set of directories whose names correspond to the
second byte of the content key.  Files within the lower directories
have names based on the hex value of the rest of the content key.

So if a file has a content key like 0xABCDEF01234... then it will
be found at AB/CD/EF01234...

The system supports either SHA1 or SHA3, where the latter means
the 256-bit version of Keccak.

The actual Go code is in u.go.  Test code is in three files:
hash_test.go contains common tests, and sha1_test.go and
sha3_test.go contain driver code specifc to SHA1 and SHA3 respectively.

myData, myU1, and myU3 are scratch directories used in testing.
The file named abc contains just that string and is used in an
extremely simple check that SHA3 is functionng correctly.  
The Keccak distribution (KeccakKat-3.zip) contains a very large collection 
of test vectors.

## On-line Documentation

More information on the **xlU_go** project can be found [here](https://jddixon.github.io/xlU_go)