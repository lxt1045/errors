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


#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"


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
