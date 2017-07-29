// +build !amd64 noasm

package xor

func bytes(dst, src1, src2 []byte) {
	bytesNoSIMD(dst, src1, src2)
}

func matrix(dst, src []byte) {
	matrixNoSIMD(dst, src)
}
