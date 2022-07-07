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

//go:build (amd64 || amd64p32) && gc && go1.5
// +build amd64 amd64p32
// +build gc
// +build go1.5

package errors

import (
	"runtime"
	"strconv"
	"unsafe"
	_ "unsafe"
)

type _func struct {
	entryoff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
	nameoff  int32  // function name

	args        int32  // in/out args size
	deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

	pcsp      uint32
	pcfile    uint32
	pcln      uint32
	npcdata   uint32
	cuOffset  uint32 // runtime.cutab offset of this function's CU
	funcID    funcID // set for certain special runtime functions
	flag      funcFlag
	_         [1]byte // pad
	nfuncdata uint8   // must be last, must end on a uint32-aligned boundary
}
type funcID uint8
type funcFlag uint8

func raw(f *runtime.Func) *_func {
	return (*_func)(unsafe.Pointer(f))
}

type Caller struct {
	pc   uintptr
	line string
	name string
}

func toCallers(pcs []uintptr) (callers []Caller) {
	for f, i := runtime.CallersFrames(pcs), 0; ; i++ {
		ff, more := f.Next()
		line := ff.File + ":" + strconv.Itoa(ff.Line)
		callers = append(callers, Caller{
			pc:   pcs[i],
			line: line,
			name: ff.Function,
		})
		if !more {
			break
		}
	}
	return
}

func GetFuncInfo(pc uintptr) *_func {
	funcInfo := runtime.FuncForPC(pc)
	return raw(funcInfo)
}
