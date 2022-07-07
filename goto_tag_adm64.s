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

// func Tag() error
TEXT ·Tag(SB),NOSPLIT,$16-16
	NO_LOCAL_POINTERS
	MOVQ	$0, ret+0(FP)
	MOVQ	$0, ret+8(FP)
	GO_RESULTS_INITIALIZED
	MOVQ	pc-8(FP), R13
	MOVQ	R13, ·pc(SB)

	MOVQ	R13, 0(SP)
	CALL ·storeTag(SB)

	RET



// Tag 和 GotoTag 建没有插入调用怎没问题
TEXT ·GotoTag(SB),NOSPLIT, $16-16
	NO_LOCAL_POINTERS
	MOVQ	pc-8(FP), R13
	MOVQ	R13, 0(SP)
	CALL ·loadTag(SB)

	CMPQ	err+8(FP), $0 //type eface struct { _type *_type; data  unsafe.Pointer }
	JHI	gototag
	RET
gototag:
	// MOVQ	err+0(FP),R13  //err
	// MOVQ	err+8(FP),R14  //err

	// MOVQ	·pc(SB),CX  // get pc
	MOVQ	8(SP), CX


	// MOVQ	CX, 0(SP)
	MOVQ	CX, retpc-8(FP)   // 因为 ·GotoTagx 和 ·Tagx 的栈结构是一样的，所以，，，
	// MOVQ	R13, 8(SP)
	// MOVQ	R14, 16(SP)
	RET


TEXT ·Jump1(SB),NOSPLIT, $0-24
	NO_LOCAL_POINTERS

	MOVQ	pc+0(FP), CX // retpc -> return addr

	// MOVQ	CX, 0(SP)
	MOVQ	CX, 8(BP)
	MOVQ	pc+8(FP), CX 
	MOVQ	$0, 16(BP)
	MOVQ	CX, 24(BP)
	MOVQ	pc+16(FP), CX 
	MOVQ	CX, 32(BP)
	RET

// func NewTag2() (tag, error)
// func Jump2(pc, parent uintptr, err error) uintptr
TEXT ·Jump2(SB),NOSPLIT, $0-40
	NO_LOCAL_POINTERS

	MOVQ	(BP), R14	
	MOVQ	+8(R14), R13	
	MOVQ	R13, ret+32(FP)
	
	// CMPQ	24(BP), R13
	CMPQ	pc+8(FP), R13
	JE	checkerr
	RET
checkerr:
	CMPQ	err+24(FP), $0 //type eface struct { _type *_type; data  unsafe.Pointer }
	JHI	gototag
	RET
gototag:
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


//func Jump(pc uintptr, err error)
TEXT ·Jump(SB),NOSPLIT, $0-24
	NO_LOCAL_POINTERS

	MOVQ	pc+0(FP), CX // retpc -> return addr

	// MOVQ	CX, 0(SP)
	MOVQ	CX, 8(BP)
	MOVQ	pc+8(FP), CX 
	MOVQ	CX, 16(BP)
	MOVQ	pc+16(FP), CX 
	MOVQ	CX, 24(BP)
	RET


// func NewTag() (error, func(error))
TEXT ·NewTag(SB),NOSPLIT,$16-24
	NO_LOCAL_POINTERS
	MOVQ	$0, ret+0(FP)
	MOVQ	$0, ret+8(FP)
	MOVQ	$0, ret+16(FP)
	GO_RESULTS_INITIALIZED
	MOVQ	pc-8(FP), R13
	// MOVQ	R13, ·pc(SB)

	MOVQ	R13, 0(SP)
	CALL ·newTag(SB)
	MOVQ	8(SP), R13
	MOVQ	R13, ret+0(FP)

	RET


// func NewTag2() (tag, error)
TEXT ·NewTag2(SB),NOSPLIT,$32-24
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
