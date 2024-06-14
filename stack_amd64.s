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
#include "jmp/define.h"  // replace MOVQ --> MOVEX


// func getPC() [1]uintptr
TEXT ·getPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS

	// 返回上一层调用栈的 pc; 由于getPC栈帧为0，所以没有压入BP，这里就少一层BP
	// 理论上，不压入BP的平台(非arm64、X64)也可以如此获取，就是要改成 "MOVEX	+0(BP), AX"，因为有BP才需要+8
	MOVEX	+8(BP), AX		
	// SUBX	$1, AX          // pc-1 才是真正的PC，但是 runtime.Callers 并没有这样做，为了一致性，这里先注释了
	MOVEX	AX, ret+0(FP)
	RET


// func GetPC() uintptr
TEXT ·GetPC(SB),NOSPLIT,$0-8
	NO_LOCAL_POINTERS
	MOVEX	+8(BP), AX		// 返回上一层调用栈的 pc
	// SUBX	$1, AX
	MOVEX	AX, ret+0(FP)
	RET

// func buildStack(s []uintptr) int
TEXT ·buildStack2(SB), NOSPLIT, $24-8
	NO_LOCAL_POINTERS
	MOVEX 	cap+16(FP), DX 	// s.cap
	MOVEX 	p+0(FP), AX		// s.ptr
	MOVEX	$0, CX			// loop.i=0

	MOVEX	BP, R13			// store BP
	MOVEX	+0(BP), BP      // skip +1 // 和 GetPC() 不同，此函数有参数，需要入栈 BP，所以这里要跳过一层调用栈
loop:
	CMPX	CX, DX			// if i >= s.cap { return }
	JAE	return				// 无符号大于等于就跳转

	MOVEX	+8(BP), BX		// last pc -> BX
	// SUBX	$1, BX          //  pc-1 才是真正的PC，但是 runtime.Callers 并没有这样做
	MOVEX	BX, 0(AX)(CX*8)		// s[i] = BX
	
	ADDX	$1, CX			// CX++ / i++

	MOVEX	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPX	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转

return:
	MOVEX	CX,n+24(FP) 	// ret n
	MOVEX	R13, BP			// load BP
	RET



// func buildStack(s []uintptr) int
TEXT ·buildStack(SB), NOSPLIT, $0-32
	NO_LOCAL_POINTERS
	MOVEX 	cap+16(FP), DX 	// s.cap
	MOVEX 	p+0(FP), AX		// s.ptr
	MOVEX	$0, CX			// loop.i=0
	MOVEX	BP, R13			// store BP

	// MOVEX	+0(BP), BP      // skip +1 // 和 GetPC() 不同，此函数有参数，需要入栈 BP，所以这里要跳过一层调用栈
	// CMPX	BP, CX			// if i >= s.cap { return }
	// JE	return				// 无符号大于等于就跳转
loop:
	CMPX	CX, DX			// if i >= s.cap { return }
	JAE	return				// 无符号大于等于就跳转

	MOVEX	+8(BP), BX		// last pc -> BX
	// SUBX	$1, BX          //  pc-1 才是真正的PC，但是 runtime.Callers 并没有这样做
	MOVEX	BX, 0(AX)(CX*8)		// s[i] = BX
	
	ADDX	$1, CX			// CX++ / i++

	MOVEX	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPX	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转

return:
	MOVEX	CX,n+24(FP) 	// ret n
	MOVEX	R13, BP			// load BP
	RET

