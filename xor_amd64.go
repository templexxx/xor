package xor

func init() {
	getEXT()
}

func xorMatrix(dst []byte, src [][]byte) {
	switch extension {
	case avx2:
		matrixAVX2(dst, src)
	case sse2:
		matrixSSE2(dst, src)
	default:
		matrixNoSIMD(dst, src)
	}
}

func matrixAVX2(dst []byte, src [][]byte) {
	size := len(dst)
	start := 0
	do := unitSize
	for start < size {
		end := start + do
		if end+do <= size {
			partAVX2(start, end, dst, src)
			start = start + do
		} else {
			partAVX2(start, size, dst, src)
			start = size
		}
	}
}

func partAVX2(start, end int, dst []byte, src [][]byte) {
	bytesAVX2(dst[start:end], src[0][start:end], src[1][start:end])
	for i := 2; i < len(src); i++ {
		updateAVX2(dst[start:end], src[i][start:end])
	}
}

// TODO bench no split
func matrixSSE2(dst []byte, src [][]byte) {
	size := len(dst)
	start := 0
	do := unitSize
	for start < size {
		end := start + do
		if end <= size {
			partSSE2(start, end, dst, src)
			start = start + do
		} else {
			partSSE2(start, size, dst, src)
			start = size
		}
	}
}

func partSSE2(start, end int, dst []byte, src [][]byte) {
	bytesSSE2(dst[start:end], src[0][start:end], src[1][start:end])
	for i := 2; i < len(src); i++ {
		updateSSE2(dst[start:end], src[i][start:end])
	}
}

func xorBytes(dst, src1, src2 []byte) {
	switch extension {
	case avx2:
		bytesAVX2(dst, src1, src2)
	case sse2:
		bytesSSE2(dst, src1, src2)
	default:
		bytesNoSIMD(dst, src1, src2)
	}
}
func getEXT() {
	if hasAVX2() {
		extension = avx2
	} else if hasSSE2() {
		extension = sse2
	} else {
		extension = none
	}
	return
}

//go:noescape
func mAVX2(dst []byte, src [][]byte)

//go:noescape
func bytesAVX2(dst, src1, src2 []byte)

//go:noescape
func updateAVX2(dst, src []byte)

//go:noescape
func bytesSSE2(dst, src1, src2 []byte)

//go:noescape
func updateSSE2(dst, src []byte)

//go:noescape
func hasAVX2() bool

//go:noescape
func hasSSE2() bool
