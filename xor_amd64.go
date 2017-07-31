package xor

func init() {
	getEXT()
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

func aesBlock(dst, src0, src1 []byte) {
	bytesSSE2once(dst, src0, src1)
}

func xorBytes(dst, src0, src1 []byte) {
	switch extension {
	case avx2:
		bytesAVX2(dst, src0, src1)
	case sse2:
		bytesSSE2(dst, src0, src1)
	default:
		bytesNoSIMD(dst, src0, src1)
	}
}

func bytesAVX2(dst, src0, src1 []byte) {
	size := len(dst)
	if size < 32 {
		// if has avx2 it must has sse2
		bytesSSE2mini(dst, src0, src1)
	} else if size >= 32 && size < 128 {
		bytesAVX2mini(dst, src0, src1)
	} else if size >= 128 && size <= unitSize {
		bytesAVX2small(dst, src0, src1)
	} else {
		bytesAVX2big(dst, src0, src1)
	}
}

//go:noescape
func bytesSSE2once(dst, src0, src1 []byte)

//go:noescape
func bytesSSE2mini(dst, src0, src1 []byte)

//go:noescape
func bytesAVX2mini(dst, src0, src1 []byte)

//go:noescape
func bytesAVX2big(dst, src0, src1 []byte)

//go:noescape
func bytesAVX2small(dst, src0, src1 []byte)

func bytesSSE2(dst, src0, src1 []byte) {
	size := len(dst)
	if size < 64 {
		bytesSSE2mini(dst, src0, src1)
	} else if size >= 64 && size <= unitSize {
		bytesSSE2small(dst, src0, src1)
	} else {
		bytesSSE2big(dst, src0, src1)
	}
}

//go:noescape
func bytesSSE2big(dst, src0, src1 []byte)

//go:noescape
func bytesSSE2small(dst, src0, src1 []byte)

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
	if size > unitSize {
		// use non-temporal hint
		matrixAVX2big(dst, src)
	} else {
		matrixAVX2small(dst, src)
	}
}

//go:noescape
func matrixAVX2big(dst []byte, src [][]byte)

//go:noescape
func matrixAVX2small(dst []byte, src [][]byte)

func matrixSSE2(dst []byte, src [][]byte) {
	size := len(dst)
	if size > unitSize {
		// use non-temporal hint
		matrixSSE2big(dst, src)
	} else {
		matrixSSE2big(dst, src)
	}
}

//go:noescape
func matrixSSE2big(dst []byte, src [][]byte)

//go:noescape
func matrixSSE2small(dst []byte, src [][]byte)

//go:noescape
func hasAVX2() bool

//go:noescape
func hasSSE2() bool
