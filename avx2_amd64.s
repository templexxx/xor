#include "textflag.h"

#define POS R10

// 需要判断两个循环是否结束（1： vect循环完了  2：长度循环完了）
// 或者可以外部传入，这样简单一点
// func updateAVX2(dst, src []byte)
TEXT ·updateAVX2(SB), NOSPLIT, $0
	MOVQ  dst+0(FP), AX
	MOVQ  src+24(FP), BX
	MOVQ  len+8(FP), CX
	TESTQ $127, CX
	JNZ   not_aligned    // check if aligned to 128b

aligned:
	MOVQ  $0, POS

loop128b:
	VMOVDQU (AX), Y0    // read dst
	VMOVDQU 32(AX), Y1
	VMOVDQU 64(AX), Y2
	VMOVDQU 96(AX), Y3

	VPXOR   (BX), Y0, Y0
	VPXOR   32(BX), Y1, Y1
	VPXOR   64(BX), Y2, Y2
	VPXOR   96(BX), Y3, Y3

	VMOVDQU Y0, (AX)    // update dst
	VMOVDQU Y1, 32(AX)
	VMOVDQU Y2, 64(AX)
	VMOVDQU Y3, 96(AX)


	ADDQ    $128, POS
	ADDQ    $128, AX
	ADDQ    $128, BX
	CMPQ    CX,   POS
	JNE     loop128b
	RET

loop_1b:
	MOVB   (AX), R8
	MOVB   (BX), R9
	XORB   R9, R8
	MOVB   R8, (AX)

	SUBQ   $1, CX
	ADDQ   $1, AX
	ADDQ   $1, BX
	TESTQ  $7, CX
	JNZ    loop_1b

	CMPQ   CX, $0
	JE     ret
	TESTQ  $127, CX
	JZ     aligned

not_aligned:
	TESTQ   $7, CX   //
	JNE     loop_1b
	MOVQ    CX, DX
	ANDQ    $127, DX  // deal with <128 part

loop_8b:
	MOVQ   (AX), R11
	MOVQ   (BX), R12
	XORQ   R12, R11
	MOVQ   R11, (AX)
	ADDQ   $8, AX
    ADDQ   $8, BX
	SUBQ   $8, CX
	SUBQ   $8, DX
	JNE    loop_8b

	CMPQ   CX, $128
	JGE    aligned
	RET

ret:
	RET

// func bytesAVX2(dst, src1, src2 []byte)
TEXT ·bytesAVX2(SB), NOSPLIT, $0
	MOVQ  dst+0(FP), AX
	MOVQ  src1+24(FP), BX
	MOVQ  src2+48(FP), R13
	MOVQ  len+8(FP), CX
	TESTQ $127, CX
	JNZ   not_aligned    // check if aligned to 128b

aligned:
	MOVQ  $0, POS

loop128b:
	VMOVDQU (R13), Y0    // read src2
	// TODO find way to get index
	// maybe use macro like sha1block_amd64.s or memmove.s
	VMOVDQU 32(R13), Y1
	VMOVDQU 64(R13), Y2
	VMOVDQU 96(R13), Y3

	VPXOR   (BX), Y0, Y0
	VPXOR   32(BX), Y1, Y1
	VPXOR   64(BX), Y2, Y2
	VPXOR   96(BX), Y3, Y3

	VMOVDQU Y0, (AX)    // update dst
	VMOVDQU Y1, 32(AX)
	VMOVDQU Y2, 64(AX)
	VMOVDQU Y3, 96(AX)

	ADDQ    $128, POS
	ADDQ    $128, AX
	ADDQ    $128, BX
	ADDQ    $128, R13
	CMPQ    CX,   POS
	JNE     loop128b
	RET

loop_1b:
	MOVB   (R13), R8
	MOVB   (BX), R9
	XORB   R9, R8
	MOVB   R8, (AX)

	SUBQ   $1, CX
	ADDQ   $1, AX
	ADDQ   $1, BX
	ADDQ   $1, R13
	TESTQ  $7, CX
	JNZ    loop_1b

	CMPQ   CX, $0
	JE     ret
	TESTQ  $127, CX
	JZ     aligned

not_aligned:
	TESTQ   $7, CX   //
	JNE     loop_1b
	MOVQ    CX, DX
	ANDQ    $127, DX  // deal with <128 part

loop_8b:
	MOVQ   (R13), R11
	MOVQ   (BX), R12
	XORQ   R12, R11
	MOVQ   R11, (AX)
	ADDQ   $8, AX
    ADDQ   $8, BX
    ADDQ   $8, R13
	SUBQ   $8, CX
	SUBQ   $8, DX
	JNE    loop_8b

	CMPQ   CX, $128
	JGE    aligned
	RET

ret:
	RET

TEXT ·hasAVX2(SB), NOSPLIT, $0
	XORQ AX, AX
	XORQ CX, CX
	ADDL $7, AX
	CPUID              // when CPUID excutes with AX set to 07H, feature info is ret in BX
	SHRQ $5, BX        // AVX -> BX[5] = 1
	ANDQ $1, BX
	MOVB BX, ret+0(FP)
	RET
