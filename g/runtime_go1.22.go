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

//go:build gc && go1.22
// +build gc,go1.22

package g

import (
	"sync/atomic"
)

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
	stack        stack
	stackguard0  uintptr
	stackguard1  uintptr
	_panic       uintptr
	_defer       uintptr // *_defer
	m            uintptr
	sched        gobuf
	syscallsp    uintptr
	syscallpc    uintptr
	stktopsp     uintptr
	param        uintptr
	atomicstatus atomic.Uint32
	stackLock    uint32
	goid         uint64
}

// type _defer struct {
// 	heap      bool
// 	rangefunc bool
// 	sp        uintptr
// 	pc        uintptr
// 	fn        func()
// 	link      *_defer

// 	head *atomic.Pointer[_defer]
// }

// func (d *_defer) GetLink() *_defer {
// 	return d.link
// }

// func (d *_defer) SetLink(l *_defer) {
// 	d.link = l
// }

// func (d *_defer) GetHead() *_defer {
// 	return d.head.Load()
// }

// func (d *_defer) SetHead(l *_defer) {
// 	d.head.Store(l)
// }
