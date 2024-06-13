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

//go:build (386 || amd64 || amd64p32 || arm64) && gc && go1.5

#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"


// func TryJmp(pc PC, err error)
TEXT ·TryLong(SB),NOSPLIT, $0-48
    NO_LOCAL_POINTERS
    GO_RESULTS_INITIALIZED

    // checkerr:
    CMPQ    err+32(FP), $0 // err.data==nil ;type eface struct { _type *_type; data  unsafe.Pointer }
    JHI    checkparent
    RET


// 需要找到 Set() 函数调用的那个函数。
checkparent:
    MOVQ    pc+16(FP), R13  // get parent 
    MOVQ    BP, BX  // store BP
loop:
    CMPQ    8(BP), R13     // // parent 是否相等；不相等则直接返回
    JE    gotojmp

	MOVQ	+0(BP), BP 		// last BP; 展开调用栈至上一层
	CMPQ	BP, $0 			// if (BP) <= 0 { return }
	JA loop					// 无符号大于就跳转
    MOVQ    BX, BP  // load BP
    RET                     // 找不到，则不处理

gotojmp:
    MOVQ    pc+0(FP), CX // jmp.pc
    MOVQ    pc+8(FP), R15 // jmp.sp
    MOVQ    BP, BX
    SUBQ    R15, BX        // SP 最终值
    // MOVQ    pc+16(FP), R13 // jmp.parent
    MOVQ    pc+24(FP), DX // jmp._defer
    MOVQ    pc+32(FP), AX // err.type
    MOVQ    pc+40(FP), R14 // err.data

    MOVQ    BX, SP  // 恢复 SP 物理寄存器
    MOVQ    CX, retaddr-8(FP)  // 恢复 ret addr
    MOVQ    CX, pc+0(FP) // jmp.pc
    MOVQ    R15, pc+8(FP) // jmp.sp
    MOVQ    R13, pc+16(FP) // jmp.parent
    MOVQ    DX, pc+24(FP) // jmp._defer
    MOVQ    AX, pc+32(FP) // err.type
    MOVQ    R14, pc+40(FP) // err.data

    // 以下重置 PC 变量，实现多次调用; Set()函数和Try()函数参数一样，所以可以不处理
    // MOVQ    CX, 16(BP)  // Setjmp.pc
    // MOVQ    parent+8(FP), CX // jmp.parent
    // MOVQ    CX, 24(BP)  // Setjmp.parent
    // MOVQ    parent+16(FP), CX // jmp._defer
    // MOVQ    CX, 32(BP)  // Setjmp._defer
    // MOVQ    pc+24(FP), CX 
    // MOVQ    CX, 40(BP)  // err
    // MOVQ    pc+32(FP), CX 
    // MOVQ    CX, 48(BP)  // err

    // 恢复defer链表。
    // 因为debug时使用使用defer链表而release时不使用，会导致两个情境下执行效果不一致。
    // 所以此处重置defer链表，使debug模式下表现和release一致。
    MOVQ (TLS), AX    // runtime.g
    ADDQ ·defer_offset(SB), AX  // &g._defer
    // MOVQ parent+24(FP), DX // jmp._defer
    MOVQ DX, (AX)  // g._defer = jmp._defer


    RET


