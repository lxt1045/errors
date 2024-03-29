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


// func NewHandler2() (handler, error)
// func tryJump(pc, parent uintptr, err error) uintptr
TEXT ·tryJump(SB),NOSPLIT, $0-40
	NO_LOCAL_POINTERS

	MOVQ	(BP), R14	 // get parent
	MOVQ	+8(R14), R13	
	MOVQ	R13, ret+32(FP)
	
	// CMPQ	24(BP), R13
	CMPQ	pc+8(FP), R13  // parent 是否相等
	JE	checkerr
	RET
checkerr:
	CMPQ	err+24(FP), $0 //type eface struct { _type *_type; data  unsafe.Pointer }
	JHI	gotohandler
	RET
gotohandler:
	MOVQ	pc+0(FP), CX // retpc -> return addr
	MOVQ	CX, 8(BP)  //ret addr
	MOVQ	CX, 16(BP)  //t.pc
	MOVQ	pc+8(FP), CX  
	MOVQ	CX, 24(BP)  //t.parent
	MOVQ	pc+16(FP), CX 
	MOVQ	CX, 32(BP)  // err
	MOVQ	pc+24(FP), CX 
	MOVQ	CX, 40(BP)  // err
	RET



// func NewHandler() (handler, error)
TEXT ·NewHandler(SB),NOSPLIT,$32-24
	NO_LOCAL_POINTERS
	MOVQ	$0, ret+0(FP)  // 返回值清零
	MOVQ	$0, ret+8(FP)
	MOVQ	$0, ret+16(FP)
	GO_RESULTS_INITIALIZED
	MOVQ	pc-8(FP), R13  // pc
	MOVQ	R13, ret+0(FP)

	MOVQ	(BP), R14	   // parent_pc
	MOVQ	+8(R14), R13
	MOVQ	R13, ret+8(FP)

	RET
