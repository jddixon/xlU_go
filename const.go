package xlU_go

// The version MUST consist of three parts separated by dots,
// with each part being one or two digits.  It is converted
// into a uint32 in in_handler.go init()
const (
	// the version number tracked in CHANGES
	VERSION      = "0.5.0"
	VERSION_DATE = "2014-05-10"
)

const (
	SHA1_LEN = 40 // length of hex version
	SHA3_LEN = 64

	//           ....x....1....x....2....x....3....x....4
	SHA1_NONE = "0000000000000000000000000000000000000000"

	//          ....x....1....x....2....x....3....x....4....x....5....x....6....
	SHA3_NONE           = "0000000000000000000000000000000000000000000000000000000000000000"
	DEFAULT_BUFFER_SIZE = 128 * 1024
)

type DirStruc int

const (
	DIR_FLAT DirStruc = iota
	DIR16x16
	DIR256x256
)
