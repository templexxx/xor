#include "textflag.h"

#define DST BX
#define SRC SI
#define VECT CX
#define LEN DX
#define POS R8
// store vect
#define TMP1 R9
// store vect_pos
#define TMP2 R10
// store addr of data/parity
#define TMP3 R11
// store
#define TMP4 R12

// func mAVX2(dst []byte, src [][]byte)
TEXT ·mAVX2(SB), NOSPLIT, $0
	MOVQ  dst+0(FP), DST
	MOVQ  src+24(FP), SRC
	MOVQ  vec+32(FP), VECT
	MOVQ  len+8(FP), LEN
	TESTQ $127, LEN
	JNZ   not_aligned    // check if aligned to 128b

aligned:
	MOVQ  $0, POS

loop128b:
	MOVQ VECT, TMP1
	MOVQ $0, TMP2
	MOVQ (SRC)(TMP2*1), TMP3  // addr of src[TMP2/24]
	SUBQ $2, TMP1
	VMOVDQU (TMP3)(POS*1), Y0 // read src[TMP2/24] from POS
	VMOVDQU 32(TMP3)(POS*1), Y1
	VMOVDQU 64(TMP3)(POS*1), Y2
	VMOVDQU 96(TMP3)(POS*1), Y3

next_vect:
	ADDQ $24, TMP2
	MOVQ (SRC)(TMP2*1), TMP3
	VMOVDQU (TMP3)(POS*1), Y4 // read src[TMP2/24] from POS
    VMOVDQU 32(TMP3)(POS*1), Y5
   	VMOVDQU 64(TMP3)(POS*1), Y6
   	VMOVDQU 96(TMP3)(POS*1), Y7
   	VPXOR Y4, Y0, Y0
   	VPXOR Y5, Y1, Y1
   	VPXOR Y6, Y2, Y2
   	VPXOR Y7, Y3, Y3
   	SUBQ    $1, TMP1
   	JGE     next_vect

	// TODO Non-temporal
	VMOVDQU Y0, (DST)(POS*1)
	VMOVDQU Y1, 32(DST)(POS*1)
	VMOVDQU Y2, 64(DST)(POS*1)
	VMOVDQU Y3, 96(DST)(POS*1)
	ADDQ    $128, POS
	CMPQ    LEN, POS
	JNE     loop128b
	RET

loop_1b:
	MOVQ VECT, TMP1
	MOVQ $0, TMP2
	MOVQ (SRC)(TMP2*1), TMP3
	SUBQ $2, TMP1
	MOVB -1(TMP3)(LEN*1), R13

next_vect_1b:
	ADDQ $24, TMP2
	MOVQ (SRC)(TMP2*1), TMP3
	MOVB -1(TMP3)(LEN*1), R14
	XORB R14, R13
	SUBQ $1, TMP1
	JGE  next_vect_1b

	MOVB R13, -1(DST)(LEN*1)
	SUBQ $1, LEN
	TESTQ $7, LEN
	JNZ loop_1b

	CMPQ LEN, $0
	JE  ret
	TESTQ $127, LEN
	JZ  aligned

not_aligned:
	TESTQ   $7, LEN   //
	JNE     loop_1b
	MOVQ    LEN, TMP4
	ANDQ    $127, TMP4  // deal with <128 part

loop_8b:
	MOVQ VECT, TMP1
	MOVQ $0, TMP2
	MOVQ (SRC)(TMP2*1), TMP3  // addr of src[TMP2/24]
	SUBQ $2, TMP1
	MOVQ  -8(TMP3)(LEN*1), R13

next_vect_8b:
	ADDQ $24, TMP2
    MOVQ (SRC)(TMP2*1), TMP3
    MOVQ -8(TMP3)(LEN*1), R14
    XORQ R14, R13
    SUBQ    $1, TMP1
    JGE     next_vect_8b

    MOVQ R13, -8(DST)(LEN*1)
    SUBQ $8, LEN
    SUBQ $8, TMP4
    JG   loop_8b

    CMPQ  LEN, $128
    JGE   aligned
    RET

ret:
	RET
