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
#include "define.h"  // replace MOVQ --> MOVEX

// func TryJmp(pc PC, err error)
TEXT ·TryLong(SB),NOSPLIT, $0-48
    NO_LOCAL_POINTERS
    GO_RESULTS_INITIALIZED

    // checkerr:
    CMPX    err+32(FP), $0 // err.data==nil ;type eface struct { _type *_type; data  unsafe.Pointer }
    JHI    checkparent
    RET


// 需要找到 Set() 函数调用的那个函数。
checkparent:
    MOVEX    pc+16(FP), R13  // get parent 
    MOVEX    BP, BX  // store BP
loop:
    CMPX    8(BP), R13     // // parent 是否相等；不相等则直接返回
    JE    gotojmp

    MOVEX    +0(BP), BP         // last BP; 展开调用栈至上一层
    CMPX    BP, $0             // if (BP) <= 0 { return }
    JA loop                    // 无符号大于就跳转
    MOVEX    BX, BP  // load BP
    RET                     // 找不到，则不处理

gotojmp:
    MOVEX    pc+0(FP), CX // jmp.pc
    MOVEX    pc+8(FP), R15 // jmp.sp
    MOVEX    BP, BX
    SUBX    R15, BX        // SP 最终值
    // MOVEX    pc+16(FP), R13 // jmp.parent
    MOVEX    pc+24(FP), DX // jmp._defer
    MOVEX    pc+32(FP), AX // err.type
    MOVEX    pc+40(FP), R14 // err.data

    MOVEX    BX, SP  // 恢复 SP 物理寄存器
    MOVEX    CX, retaddr-8(FP)  // 恢复 ret addr
    MOVEX    CX, pc+0(FP) // jmp.pc
    MOVEX    R15, pc+8(FP) // jmp.sp
    MOVEX    R13, pc+16(FP) // jmp.parent
    MOVEX    DX, pc+24(FP) // jmp._defer
    MOVEX    AX, pc+32(FP) // err.type
    MOVEX    R14, pc+40(FP) // err.data

    // 以下重置 PC 变量，实现多次调用; Set()函数和Try()函数参数一样，所以可以不处理
    // MOVEX    CX, 16(BP)  // Setjmp.pc
    // MOVEX    parent+8(FP), CX // jmp.parent
    // MOVEX    CX, 24(BP)  // Setjmp.parent
    // MOVEX    parent+16(FP), CX // jmp._defer
    // MOVEX    CX, 32(BP)  // Setjmp._defer
    // MOVEX    pc+24(FP), CX 
    // MOVEX    CX, 40(BP)  // err
    // MOVEX    pc+32(FP), CX 
    // MOVEX    CX, 48(BP)  // err

    // 恢复defer链表。
    // 因为debug时使用使用defer链表而release时不使用，会导致两个情境下执行效果不一致。
    // 所以此处重置defer链表，使debug模式下表现和release一致。
    MOVEX (TLS), AX    // runtime.g
    ADDX ·defer_offset(SB), AX  // &g._defer
    // MOVEX parent+24(FP), DX // jmp._defer
    MOVEX DX, (AX)  // g._defer = jmp._defer


    RET


