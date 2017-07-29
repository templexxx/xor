#include "textflag.h"

#define POS R10

// 需要判断两个循环是否结束（1： vect循环完了  2：长度循环完了）
// 或者可以外部传入，这样简单一点
// func updateSSE2(dst, src []byte)
TEXT ·updateSSE2(SB), NOSPLIT, $0
	MOVQ  dst+0(FP), AX
	MOVQ  src+24(FP), BX
	MOVQ  len+8(FP), CX
	TESTQ $63, CX
	JNZ   not_aligned    // check if aligned to 64b

aligned:
	MOVQ  $0, POS

loop64b:
	MOVOU (AX), X0    // read dst
	MOVOU 16(AX), X1
	MOVOU 32(AX), X2
	MOVOU 48(AX), X3

	MOVOU  (BX), X4
	PXOR   X4, X0
	MOVOU  16(BX), X5
	PXOR   X5, X1
	MOVOU  32(BX), X6
	PXOR   X6, X2
	MOVOU  48(BX), X7
	PXOR   X7, X3

	MOVOU X0, (AX)    // update dst
	MOVOU X1, 16(AX)
	MOVOU X2, 32(AX)
	MOVOU X3, 48(AX)


	ADDQ    $64, POS
	ADDQ    $64, AX
	ADDQ    $64, BX
	CMPQ    CX,   POS
	JNE     loop64b
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
	TESTQ  $63, CX
	JZ     aligned

not_aligned:
	TESTQ   $7, CX   //
	JNE     loop_1b
	MOVQ    CX, DX
	ANDQ    $63, DX  // deal with <64 part

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

	CMPQ   CX, $64
	JGE    aligned
	RET

ret:
	RET

// func bytesSSE2(dst, src1, src2 []byte)
TEXT ·bytesSSE2(SB), NOSPLIT, $0
	MOVQ  dst+0(FP), AX
	MOVQ  src1+24(FP), BX
	MOVQ  src2+48(FP), R13
	MOVQ  len+8(FP), CX
	TESTQ $63, CX
	JNZ   not_aligned    // check if aligned to 64b

aligned:
	MOVQ  $0, POS

loop64b:
	MOVOU (R13), X0    // read src2
	MOVOU 16(R13), X1
	MOVOU 32(R13), X2
	MOVOU 48(R13), X3

	MOVOU  (BX), X4
	PXOR   X4, X0
	MOVOU  16(BX), X5
	PXOR   X5, X1
	MOVOU  32(BX), X6
	PXOR   X6, X2
	MOVOU  48(BX), X7
	PXOR   X7, X3

	MOVOU X0, (AX)    // update dst
	MOVOU X1, 16(AX)
	MOVOU X2, 32(AX)
	MOVOU X3, 48(AX)


	ADDQ    $64, POS
	ADDQ    $64, AX
	ADDQ    $64, BX
	ADDQ    $64, R13
	CMPQ    CX,   POS
	JNE     loop64b
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
	TESTQ  $63, CX
	JZ     aligned

not_aligned:
	TESTQ   $7, CX   //
	JNE     loop_1b
	MOVQ    CX, DX
	ANDQ    $63, DX  // deal with <64 part

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

	CMPQ   CX, $64
	JGE    aligned
	RET

ret:
	RET

TEXT ·hasSSE2(SB), NOSPLIT, $0
	XORQ AX, AX
	INCL AX
	CPUID
	SHRQ $26, DX
	ANDQ $1, DX
	MOVB DX, ret+0(FP)
	RET
