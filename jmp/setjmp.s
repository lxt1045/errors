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

//go:build (386 || amd64 || amd64p32 || arm || arm64) && gc && go1.5

#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"


// func Set() (PC, error)
TEXT ·Set(SB),NOSPLIT,$0-48
    NO_LOCAL_POINTERS
    // MOVQ    $0, ret+0(FP)  // 返回值清零, pc
    // MOVQ    $0, ret+8(FP)  // 返回值清零, pc
    // MOVQ    $0, ret+16(FP)  // parent
    // MOVQ    $0, ret+24(FP) // _defer
    MOVQ    $0, ret+32(FP) // err
    MOVQ    $0, ret+40(FP) // err
    GO_RESULTS_INITIALIZED
    MOVQ    pc-8(FP), R13  // pc
    MOVQ    R13, ret+0(FP)
    MOVQ    BP, AX
    SUBQ    SP, AX
    MOVQ    AX, ret+8(FP)   // 因为是拷贝栈，所以SP不能直接存，只能存SP和BP的差值！！！

    // MOVQ    (BP), R14       // parent_pc
    // MOVQ    +8(R14), R13
    // 函数栈帧大小(本地变量占用空间大小)为0时，BP未入栈
    MOVQ    8(BP), R13       // parent_pc
    MOVQ    R13, parent+16(FP)

    MOVQ (TLS), AX    // runtime.g
    ADDQ ·defer_offset(SB),AX
    MOVQ (AX), BX
    MOVQ BX, _defer+24(FP)

    RET


// func Try(pc PC, err error)
TEXT ·Try(SB),NOSPLIT, $0-48
    NO_LOCAL_POINTERS
    GO_RESULTS_INITIALIZED

// checkerr:
    CMPQ    err+32(FP), $0 // err.data==nil ;type eface struct { _type *_type; data  unsafe.Pointer }
    JHI    checkparent
    RET

checkparent:
    MOVQ    8(BP), R13     // get parent    
    CMPQ    pc+16(FP), R13  // parent 是否相等；不相等则直接返回
    JE    gotojmp
    RET

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



