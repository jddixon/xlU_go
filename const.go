package xlU_go

// The version MUST consist of three parts separated by dots,
// with each part being one or two digits.  It is converted
// into a uint32 in in_handler.go init()
const (
	// the version number tracked in CHANGES
	VERSION      = "0.6.0"
	VERSION_DATE = "2014-10-07"
)

const (
	SHA1_HEX_LEN = 40 // length of hex version of digest
	SHA3_HEX_LEN = 64

	//               ....x....1....x....2....x....3....x....4
	SHA1_HEX_NONE = "0000000000000000000000000000000000000000"

	//               ....x....1....x....2....x....3....x....4....x....5....x....6....
	SHA3_HEX_NONE       = "0000000000000000000000000000000000000000000000000000000000000000"
	DEFAULT_BUFFER_SIZE = 256 * 1024
)

type DirStruc int

const (
	DIR_FLAT DirStruc = iota
	DIR16x16
	DIR256x256
)
