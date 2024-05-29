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

//go:build amd64
// +build amd64

#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"


// func longjmp(jmp jump, err error) uintptr
TEXT ·longjmp(SB),NOSPLIT, $0-48
	NO_LOCAL_POINTERS

	MOVQ	(BP), R14	 // get parent
	MOVQ	+8(R14), R13	
	MOVQ	R13, ret+40(FP)  // return parent
	
	CMPQ	pc+8(FP), R13  // parent 是否相等
	JE	checkerr
	MOVQ	R13, ret+40(FP)  // return parent
	RET
checkerr:
	CMPQ	err+32(FP), $0 // err.data==nil ;type eface struct { _type *_type; data  unsafe.Pointer }
	JHI	gotohandler
	RET
gotohandler:
	MOVQ	pc+0(FP), CX // jmp.pc
	MOVQ	CX, 8(BP)  // ret addr

	// 以下 jump 变量重置，理论上也可以不处理？ 
	MOVQ	CX, 16(BP)  // Setjmp.pc
	MOVQ	parent+8(FP), CX // jmp.parent
	MOVQ	CX, 24(BP)  // Setjmp.parent
	MOVQ	parent+16(FP), CX // jmp._defer
	MOVQ	CX, 32(BP)  // Setjmp._defer

    // 恢复defer链表。
    // 因为debug时使用使用defer链表而release时不使用，会导致两个情境下执行效果不一致。
	// 所以此处重置defer链表，使debug模式下表现和release一致。
    MOVQ (TLS), AX    // runtime.g
	ADDQ ·defer_offset(SB), AX  // &g._defer
    MOVQ CX, (AX)  // g._defer = jmp._defer


	MOVQ	pc+24(FP), CX 
	MOVQ	CX, 40(BP)  // err
	MOVQ	pc+32(FP), CX 
	MOVQ	CX, 48(BP)  // err
	RET



// func Setjmp() (handler, error)
TEXT ·Setjmp(SB),NOSPLIT,$0-40
	NO_LOCAL_POINTERS
	MOVQ	$0, ret+0(FP)  // 返回值清零, pc
	MOVQ	$0, ret+8(FP)  // parent
	MOVQ	$0, ret+16(FP) // _defer
	MOVQ	$0, ret+24(FP) // err
	MOVQ	$0, ret+32(FP) // err
	GO_RESULTS_INITIALIZED
	MOVQ	pc-8(FP), R13  // pc
	MOVQ	R13, ret+0(FP)

	// MOVQ	(BP), R14	   // parent_pc
	// MOVQ	+8(R14), R13
	// 函数栈帧大小(本地变量占用空间大小)为0时，BP未入栈
	MOVQ	8(BP), R13	   // parent_pc
	MOVQ	R13, parent+8(FP)

    MOVQ (TLS), AX    // runtime.g
	ADDQ ·defer_offset(SB),AX
    MOVQ (AX), BX
    MOVQ BX, _defer+16(FP)

	RET


// // func GetDefer() *_defer
// TEXT ·GetDefer(SB), NOSPLIT, $0-8
//     MOVQ (TLS), AX
//     ADDQ ·g__defer_offset(SB),AX
//     MOVQ (AX), BX
//     MOVQ BX, ret+0(FP)  
//     RET

// // func getgi() interface{}
// TEXT ·getgi(SB), NOSPLIT, $32-16
//     NO_LOCAL_POINTERS

//     MOVQ $0, ret_type+0(FP)
//     MOVQ $0, ret_data+8(FP)
//     GO_RESULTS_INITIALIZED

//     // get runtime.g
//     MOVQ (TLS), AX

//     // get runtime.g type
//     MOVQ $type·runtime·g(SB), BX

//     // return interface{}
//     MOVQ BX, ret_type+0(FP)
//     MOVQ AX, ret_data+8(FP)
//     RET

// // func getdeferi() interface{}
// TEXT ·getdeferi(SB), NOSPLIT, $32-16
//     NO_LOCAL_POINTERS

//     MOVQ $0, ret_type+0(FP)
//     MOVQ $0, ret_data+8(FP)
//     GO_RESULTS_INITIALIZED

//     // get runtime._defer
//     MOVQ 0, AX

//     // get runtime._defer type
//     MOVQ $type·runtime·_defer(SB), BX

//     // return interface{}
//     MOVQ BX, ret_type+0(FP)
//     MOVQ AX, ret_data+8(FP)
//     RET
