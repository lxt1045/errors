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


// func getPC() [1]uintptr
TEXT ·getPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS

	// 返回上一层调用栈的 pc; 由于getPC栈帧为0，所以没有压入BP，这里就少一层BP
	// 理论上，不压入BP的平台(非arm64、X64)也可以如此获取，就是要改成 "MOVQ	+0(BP), AX"，因为有BP才需要+8
	MOVQ	+8(BP), AX		
	// SUBQ	$1, AX          // pc-1 才是真正的PC，但是 runtime.Callers 并没有这样做，为了一致性，这里先注释了
	MOVQ	AX, ret+0(FP)
	RET


// func GetPC() uintptr
TEXT ·GetPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS
	MOVQ	+8(BP), AX		// 返回上一层调用栈的 pc
	// SUBQ	$1, AX
	MOVQ	AX, ret+0(FP)
	RET

// func buildStack(s []uintptr) int
TEXT ·buildStack(SB), NOSPLIT, $24-8
	NO_LOCAL_POINTERS
	MOVQ 	cap+16(FP), DX 	// s.cap
	MOVQ 	p+0(FP), AX		// s.ptr
	MOVQ	$0, CX			// loop.i=0

	MOVQ	+0(BP), BP      // skip +1 // 和 GetPC() 不同，此函数有参数，需要入栈 BP，所以这里要跳过一层调用栈
loop:
	CMPQ	CX, DX			// if i >= s.cap { return }
	JAE	return				// 无符号大于等于就跳转

	MOVQ	+8(BP), BX		// last pc -> BX
	// SUBQ	$1, BX          //  pc-1 才是真正的PC，但是 runtime.Callers 并没有这样做
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
	// SUBQ	$1, BX
	MOVQ	BX, 0(AX)(CX*8)		// s[i] = BX
	
	MOVQ	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPQ	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转

return:
	MOVQ	CX,n+24(FP) 	// ret n
	RET

