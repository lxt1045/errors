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

package jmp

import (
	_ "unsafe" //nolint:bgolint
)

// 类似 C 语言的 setjmp.h 里的 setjmp() 函数
func Set() (PC, error) //nolint:bgolint

// 类似 C 语言的 setjmp.h 里的 longjmp() 函数
// 注意 Try() 必须和生成 PC 的 Set() 函数在同一个函数内，否则会无效
func Try(pc PC, err error)

type PC struct {
	pc     uintptr //nolint:unused
	sp     uintptr //nolint:unused
	parent uintptr //nolint:unused
	_defer uintptr //nolint:unused

	// noCopy noCopy //nolint:unused
}
