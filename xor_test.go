package xor

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestVerifyBytesNoSIMD(t *testing.T) {
	dst := make([]byte, 10)
	src1 := []byte{0, 4, 22, 3, 45, 7, 8, 9, 10, 11}
	src2 := []byte{9, 4, 221, 32, 145, 17, 18, 19, 110, 117}
	for i := 0; i < 10; i++ {
		dst[i] = src1[i] ^ src2[i]
	}
	expect := []byte{9, 0, 203, 35, 188, 22, 26, 26, 100, 126}
	bytesNoSIMD(dst, src1, src2)
	if !bytes.Equal(expect, dst) {
		t.Fatal("xor fault")
	}
}

func TestVerifyBytes(t *testing.T) {
	size := 9999
	dst := make([]byte, size)
	src1 := make([]byte, size)
	src2 := make([]byte, size)
	expect := make([]byte, size)
	rand.Seed(7)
	fillRandom(src1)
	rand.Seed(8)
	fillRandom(src2)
	for i := 0; i < size; i++ {
		expect[i] = src1[i] ^ src2[i]
	}
	xorBytes(dst, src1, src2)
	if !bytes.Equal(expect, dst) {
		t.Fatal("xor fault")
	}
}

func TestMatrixNoSIMD(t *testing.T) {
	size := 9999
	numSRC := 10
	dst := make([]byte, size)
	expect := make([]byte, size)
	src := make([][]byte, numSRC)
	for i := 0; i < numSRC; i++ {
		src[i] = make([]byte, size)
		rand.Seed(int64(i))
		fillRandom(src[i])
	}
	for i := 0; i < size; i++ {
		expect[i] = src[0][i] ^ src[1][i]
	}
	for i := 2; i < numSRC; i++ {
		for j := 0; j < size; j++ {
			expect[j] ^= src[i][j]
		}
	}
	matrixNoSIMD(dst, src)
	if !bytes.Equal(expect, dst) {
		t.Fatal("xor fault")
	}
}

func TestMatrix(t *testing.T) {
	size := 9999
	numSRC := 10
	dst := make([]byte, size)
	expect := make([]byte, size)
	src := make([][]byte, numSRC)
	for i := 0; i < numSRC; i++ {
		src[i] = make([]byte, size)
		rand.Seed(int64(i))
		fillRandom(src[i])
	}
	for i := 0; i < size; i++ {
		expect[i] = src[0][i] ^ src[1][i]
	}
	for i := 2; i < numSRC; i++ {
		for j := 0; j < size; j++ {
			expect[j] ^= src[i][j]
		}
	}
	mAVX2(dst, src)
	if !bytes.Equal(expect, dst) {
		t.Fatal("xor fault")
	}
}

func BenchmarkNewBytes1K(b *testing.B) {
	benchmarkNew(b, 2, 1024)
}
func BenchmarkNewBytes1400B(b *testing.B) {
	benchmarkNew(b, 2, 1400)
}
func BenchmarkNewBytes16M(b *testing.B) {
	benchmarkNew(b, 2, 16*1024*1024)
}
func BenchmarkNewBytes10x1K(b *testing.B) {
	benchmarkNew(b, 10, 1024)
}
func BenchmarkNewBytes10x1400B(b *testing.B) {
	benchmarkNew(b, 10, 1400)
}
func BenchmarkNewBytes10x16M(b *testing.B) {
	benchmarkNew(b, 10, 16*1024*1024)
}

func benchmarkNew(b *testing.B, numSRC, size int) {
	src := make([][]byte, numSRC)
	dst := make([]byte, size)
	for i := 0; i < numSRC; i++ {
		rand.Seed(int64(i))
		src[i] = make([]byte, size)
		fillRandom(src[i])
	}
	mAVX2(dst, src)
	b.SetBytes(int64(size * numSRC))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mAVX2(dst, src)
	}
}

// xor bytes
func BenchmarkXorBytes1K(b *testing.B) {
	benchmarkBytes(b, 1024)
}
func BenchmarkXorBytes1400B(b *testing.B) {
	benchmarkBytes(b, 1400)
}
func BenchmarkXorBytes16M(b *testing.B) {
	benchmarkBytes(b, 16*1024*1024)
}

// xor bytes no simd
func BenchmarkXorBytesNoSIMD1K(b *testing.B) {
	benchmarkBytesNoSIMD(b, 1024)
}
func BenchmarkXorBytesNoSIMD1400B(b *testing.B) {
	benchmarkBytesNoSIMD(b, 1400)
}
func BenchmarkXorBytesNoSIMD16M(b *testing.B) {
	benchmarkBytesNoSIMD(b, 16*1024*1024)
}

// xor matrix
func BenchmarkXorMatrix10x1K(b *testing.B) {
	benchmarkMatrix(b, 10, 1024)
}
func BenchmarkXorMatrix10x1400B(b *testing.B) {
	benchmarkMatrix(b, 10, 1400)
}
func BenchmarkXorMatrix10x16M(b *testing.B) {
	benchmarkMatrix(b, 10, 16*1024*1024)
}

// TODO other size bench test see what happend
// TODO  why 16MB so slow
// TODO drop split
// TODO drop 128 rest
// xor matrix no simd
func BenchmarkXorMatrixNoSIMD10x1K(b *testing.B) {
	benchmarkMatrixNoSIMD(b, 10, 1024)
}
func BenchmarkXorMatrixNoSIMD10x1400B(b *testing.B) {
	benchmarkMatrixNoSIMD(b, 10, 1400)
}
func BenchmarkXorMatrixNoSIMD10x16M(b *testing.B) {
	benchmarkMatrixNoSIMD(b, 10, 16*1024*1024)
}

func benchmarkBytes(b *testing.B, size int) {
	src1 := make([]byte, size)
	src2 := make([]byte, size)
	dst := make([]byte, size)
	rand.Seed(int64(0))
	fillRandom(src1)
	rand.Seed(int64(1))
	fillRandom(src2)
	//src := make([][]byte, 2)
	//src[0] = src1
	//src[1] = src2
	// warm up
	//matrixAVX2(dst, src)
	xorBytes(dst, src1, src2)
	b.SetBytes(int64(size) * 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//matrixAVX2(dst, src)
		xorBytes(dst, src1, src2)
	}
}

func benchmarkBytesNoSIMD(b *testing.B, size int) {
	src1 := make([]byte, size)
	src2 := make([]byte, size)
	dst := make([]byte, size)
	rand.Seed(int64(0))
	fillRandom(src1)
	rand.Seed(int64(1))
	fillRandom(src2)
	// warm up
	bytesNoSIMD(dst, src1, src2)
	b.SetBytes(int64(size) * 2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bytesNoSIMD(dst, src1, src2)
	}
}

func benchmarkMatrix(b *testing.B, numSRC, size int) {
	src := make([][]byte, numSRC)
	dst := make([]byte, size)
	for i := 0; i < numSRC; i++ {
		rand.Seed(int64(i))
		src[i] = make([]byte, size)
		fillRandom(src[i])
	}
	xorMatrix(dst, src)
	b.SetBytes(int64(size * numSRC))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		xorMatrix(dst, src)
	}
}

func benchmarkMatrixNoSIMD(b *testing.B, numSRC, size int) {
	src := make([][]byte, numSRC)
	dst := make([]byte, size)
	for i := 0; i < numSRC; i++ {
		rand.Seed(int64(i))
		src[i] = make([]byte, size)
		fillRandom(src[i])
	}
	matrixNoSIMD(dst, src)
	b.SetBytes(int64(size * numSRC))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		matrixNoSIMD(dst, src)
	}
}

func fillRandom(p []byte) {
	for i := 0; i < len(p); i += 7 {
		val := rand.Int63()
		for j := 0; i+j < len(p) && j < 7; j++ {
			p[i+j] = byte(val)
			val >>= 8
		}
	}
}
