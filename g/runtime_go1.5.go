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

//go:build go1.5 && !go1.6
// +build go1.5,!go1.6

package g

// Just enough of the structs from runtime/runtime2.go to get the offset to goid.
// See https://github.com/golang/go/blob/release-branch.go1.5/src/runtime/runtime2.go

type stack struct {
	lo uintptr
	hi uintptr
}

type gobuf struct {
	sp   uintptr
	pc   uintptr
	g    uintptr
	ctxt uintptr
	ret  uintptr
	lr   uintptr
	bp   uintptr
}

type g struct {
	stack       stack
	stackguard0 uintptr
	stackguard1 uintptr

	_panic       uintptr
	_defer       uintptr
	m            uintptr
	stackAlloc   uintptr
	sched        gobuf
	syscallsp    uintptr
	syscallpc    uintptr
	stkbar       []uintptr
	stkbarPos    uintptr
	param        uintptr
	atomicstatus uint32
	stackLock    uint32
	goid         int64 // Here it is!
}

type _defer struct {
	siz     int32
	started bool
	sp      uintptr // sp at time of defer
	pc      uintptr
	fn      uintptr
	_panic  uintptr // panic that is running defer
	link    *_defer
}
