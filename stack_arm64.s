


#include "go_asm.h"
#include "textflag.h"
#include "funcdata.h"

// func getPC() [1]uintptr
TEXT ·getPC(SB), NOSPLIT, $0-8
    NO_LOCAL_POINTERS
    MOVD	8(R29), R0		// load caller's PC
    MOVD	R0, ret+0(FP)
    RET

// func GetPC() uintptr
TEXT ·GetPC(SB), NOSPLIT, $0-8
    NO_LOCAL_POINTERS
    // MOVD	0(R29), R29		// 获取更上一层栈的PC
    MOVD	8(R29), R0		// load caller's PC
    MOVD	R0, ret+0(FP)
    RET

// func buildStack2(s []uintptr) int
TEXT ·buildStack2(SB), NOSPLIT, $0-32
    NO_LOCAL_POINTERS
    
    // Load slice parameters
    MOVD	p+0(FP), R0		// s.ptr
    MOVD	cap+16(FP), R2		// s.cap
    MOVD	$0, R1			// i = 0
    MOVD	R29, R19		// save original FP

    MOVD	(R29), R29		// skip current frame: move to caller's FP

loop:
    CMP	R1, R2			// compare i and cap
    BHS	done			// if i >= cap, break

    // s[i] = PC
    MOVD	8(R29), R3		// R3 = PC
    MOVD	R1, R4			// R4 = i
    LSL	$3, R4, R4		// R4 = i * 8
    ADD	R0, R4, R4		// R4 = base + i*8
    MOVD	R3, (R4)		// store PC

    ADD	$1, R1, R1		// i++

    MOVD	(R29), R29		// move to next frame
    CBNZ	R29, loop		// if FP != 0, continue

done:
    MOVD	R19, R29		// restore original FP
    MOVD	R1, n+24(FP)		// return i
    RET

// func buildStack(s []uintptr) int
TEXT ·buildStack(SB), NOSPLIT, $0-32
	NO_LOCAL_POINTERS
	MOVD	cap+16(FP), R1     // s.cap -> R1 (替代 DX) 
	MOVD	p+0(FP), R0        // s.ptr -> R0 (替代 AX)
	MOVD	$0, R2             // loop.i=0 -> R2 (替代 CX)
	MOVD	R29, R19           // 保存 FP (R29) 到 R19 (替代 R13)

loop:
	CMP	R2, R1              // 比较 i (R2) 和 cap (R1)
	// BHS	return              // 无符号大于等于跳转 (替代 JAE)  ; 这个指令有问题
    BLE return                  // 根据go源码反编译 go tool compile -S main.go 结果, 编译器使用使用BLE ; BHS 和 BLE 分别表示 ≥ 和 ≤

	MOVD	8(R29), R3         // 加载返回地址: [FP+8] -> R3 (替代 BX)
	MOVD	R3, (R0)(R2<<3)    // 存储到 s[i]: R3 -> R0 + R2*8 (替代 0(AX)(CX*8)) 

	ADD	$1, R2, R2         // i++ (R2 = R2 + 1) (替代 ADDX $1, CX)

	MOVD	(R29), R29        // 加载上级帧指针: [FP] -> FP (替代 MOVEX +0(BP), BP)
	// CMP	R29, $0            // 检查 FP 是否为 0
	// BHI	loop               // 无符号大于 0 则循环 (替代 JA) 
    CBNZ	R29, loop          // FP!=0时循环 [修复点]

return:
	MOVD	R2, n+24(FP)      // 返回值 n (替代 MOVEX CX, n+24(FP))
	MOVD	R19, R29          // 恢复 FP (替代 MOVEX R13, BP)
	RET


// 和 buildStack 等效
TEXT ·buildStack4(SB), NOSPLIT, $0-32
	NO_LOCAL_POINTERS
	MOVD	cap+16(FP), R1     // s.cap -> R1 (替代 DX) 
	MOVD	p+0(FP), R0        // s.ptr -> R0 (替代 AX)
	MOVD	$0, R2             // loop.i=0 -> R2 (替代 CX)
	MOVD	R29, R19           // 保存 FP (R29) 到 R19 (替代 R13)

    JMP check
loop:
    MOVD	8(R29), R3         // 加载返回地址: [FP+8] -> R3 (替代 BX)
	MOVD	R3, (R0)(R2<<3)    // 存储到 s[i]: R3 -> R0 + R2*8 (替代 0(AX)(CX*8)) 

	ADD	$1, R2, R2         // i++ (R2 = R2 + 1) (替代 ADDX $1, CX)

	MOVD	(R29), R29        // 加载上级帧指针: [FP] -> FP (替代 MOVEX +0(BP), BP)
	// CMP	R29, $0            // 检查 FP 是否为 0
	// BHI	loop               // 无符号大于 0 则循环 (替代 JA) 
    CBNZ	R29, check          // FP!=0时循环 [修复点]

    JMP return
check:
	CMP	R2, R1              // 比较 i (R2) 和 cap (R1)
    BGT	loop

return:
	MOVD	R2, n+24(FP)      // 返回值 n (替代 MOVEX CX, n+24(FP))
	MOVD	R19, R29          // 恢复 FP (替代 MOVEX R13, BP)
	RET




