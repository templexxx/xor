package xor

import "errors"

// SIMD Extensions
const (
	none = iota
	avx2
	sse2
)

var extension = none

// split slice for cache-friendly
const unitSize = 16 * 1024

func Bytes(dst, src1, src2 []byte) (err error) {
	err = checkSize(dst, src1, src2)
	if err != nil {
		return
	}
	xorBytes(dst, src1, src2)
	return
}

var ErrSrcNum = errors.New("xor: num of src must >= 2")

func Matrix(dst []byte, src [][]byte) (err error) {
	if len(src) < 2 {
		return ErrSrcNum
	}
	err = checkSize(dst, src...)
	if err != nil {
		return
	}
	xorMatrix(dst, src)
	return
}

var ErrShardSize = errors.New("xor: shards size equal 0 or not match")

func checkSize(dst []byte, src ...[]byte) (err error) {
	size := len(dst)
	if size == 0 {
		return ErrShardSize
	}
	for _, s := range src {
		if len(s) != size {
			return ErrShardSize
		}
	}
	return nil
}
