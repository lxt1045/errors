// MIT License
//
// Copyright (c) 2021 Xiantu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//go:build amd64 || amd64p32 || arm64
// +build amd64 amd64p32 arm64

#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"


GLOBL ·runtime_g_type(SB),NOPTR,$8
DATA ·runtime_g_type+0(SB)/8,$type·runtime·g(SB) // # 汇编初始化 go 声明的变量


// func getPC() [1]uintptr
TEXT ·getPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS
	MOVQ	+8(BP), AX		// 上一层调用栈的返回 pc
	SUBQ	$1, AX
	MOVQ	AX, ret+0(FP)
	RET


// func GetPC() uintptr
TEXT ·GetPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS
	MOVQ	+8(BP), AX		// 上一层调用栈的返回 pc
	SUBQ	$1, AX
	MOVQ	AX, ret+0(FP)
	RET

// func buildStack(s []uintptr) int
TEXT ·buildStack(SB), NOSPLIT, $24-8
	NO_LOCAL_POINTERS
	MOVQ 	cap+16(FP), DX 	// s.cap
	MOVQ 	p+0(FP), AX		// s.ptr
	MOVQ	$0, CX			// loop.i=0
loop:
	CMPQ	CX, DX			// if i >= s.cap { return }
	JAE	return				// 无符号大于等于就跳转

	MOVQ	+8(BP), BX		// last pc -> BX
	SUBQ	$1, BX
	MOVQ	BX, 0(AX)(CX*8)		// s[i] = BX
	
	ADDQ	$1, CX			// CX++ / i++

	MOVQ	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPQ	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转

return:
	MOVQ	CX,n+24(FP) 	// ret n
	RET



// func buildStack(s []uintptr) int
TEXT ·buildStack2(SB), NOSPLIT, $24-8
	NO_LOCAL_POINTERS
	MOVQ 	cap+16(FP), DX 	// s.cap
	MOVQ 	p+0(FP), AX		// s.ptr
	MOVQ	$0, CX			// loop.i

	CMPQ	DX, $1			// if s.cap<=0 { return }
	JL	return				// 有符号大于等于就跳转
	MOVQ	pc-8(FP),BX
	MOVQ	BX, 0(AX)(CX*8)	
loop:
	ADDQ	$1, CX			// CX++ / i++
	CMPQ	CX, DX			// if s.len >= s.cap { return }
	JAE	return				// 无符号大于等于就跳转

	MOVQ	+8(BP), BX		// last pc -> BX
	SUBQ	$1, BX
	MOVQ	BX, 0(AX)(CX*8)		// s[i] = BX
	
	MOVQ	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPQ	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转

return:
	MOVQ	CX,n+24(FP) 	// ret n
	RET



// func getg() unsafe.Pointer
TEXT ·Getg(SB), NOSPLIT, $0-8
    MOVQ (TLS), AX
	ADDQ ·gGoidOffset(SB),AX
    MOVQ (AX), BX
    MOVQ BX, ret+0(FP)
    RET

// func getgi() interface{}
TEXT ·getgi(SB), NOSPLIT, $32-16
    NO_LOCAL_POINTERS

    MOVQ $0, ret_type+0(FP)
    MOVQ $0, ret_data+8(FP)
    GO_RESULTS_INITIALIZED

    // get runtime.g
    MOVQ (TLS), AX

    // get runtime.g type
    MOVQ $type·runtime·g(SB), BX
    // return interface{}
    MOVQ BX, ret_type+0(FP)
    MOVQ AX, ret_data+8(FP)
    RET
