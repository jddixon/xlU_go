package xlU_go

const (
	DEFAULT_BUFFER_SIZE = 256 * 1024
)

type DirStruc int

const (
	DIR_FLAT DirStruc = iota
	DIR16x16
	DIR256x256
)
