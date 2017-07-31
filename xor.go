package xor

// SIMD Extensions
const (
	none = iota
	avx2
	sse2
)

var extension = none

// split slice for cache-friendly
const unitSize = 16 * 1024

// The arguments are assumed to be of equal length && != 0
func Bytes(dst, src0, src1 []byte) {
	xorBytes(dst, src0, src1)
}

// The arguments are assumed to be of equal length && != 0
// len(src) must >= 2
func Matrix(dst []byte, src [][]byte) {
	xorMatrix(dst, src)
}

// src1's len is 16bytes
// dst >= src0, src0 >= src1
func XorAESBlock(dst, src0, src1 []byte) {
	aesBlock(dst, src0, src1)
}
