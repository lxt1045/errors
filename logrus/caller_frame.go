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

//go:build amd64
// +build amd64

package logrus

import (
	"runtime"
	"strconv"
	"sync/atomic"
	"unsafe"
	_ "unsafe" //nolint:bgolint
)

type caller struct {
	Func string
	File string
}

var mapCaller unsafe.Pointer = func() unsafe.Pointer {
	m := make(map[uintptr]*caller)
	return unsafe.Pointer(&m)
}()

// CallerFrame 使用 Read-copy update(RCU) 缓存提高性能
func CallerFrame(l uintptr) (cf *caller) {
	mPCs := *(*map[uintptr]*caller)(atomic.LoadPointer(&mapCaller))
	cf, ok := mPCs[l]
	if !ok {
		c, _ := runtime.CallersFrames([]uintptr{l}).Next()
		cf = &caller{
			Func: c.Function,
			File: c.File + ":" + strconv.Itoa(c.Line),
		}
		mPCs2 := make(map[uintptr]*caller, len(mPCs)+10)
		mPCs2[l] = cf
		for {
			p := atomic.LoadPointer(&mapCaller)
			mPCs = *(*map[uintptr]*caller)(p)
			for k, v := range mPCs {
				mPCs2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mapCaller, p, unsafe.Pointer(&mPCs2))
			if swapped {
				break
			}
		}
	}
	return
}
