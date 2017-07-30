// +build !amd64 noasm

package xor

func xorBytes(dst, src0, src1 []byte) {
	bytesNoSIMD(dst, src0, src1)
}

func xorMatrix(dst []byte, src [][]byte) {
	matrixNoSIMD(dst, src)
}
