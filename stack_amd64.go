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

package errors

import (
	"reflect"
	"unsafe"
	_ "unsafe" //nolint:bgolint
)

func buildStack(s []uintptr) int

func Getg() int64

func getgi() interface{}

var gGoidOffset uintptr = func() uintptr { //nolint
	g := getgi()
	if f, ok := reflect.TypeOf(g).FieldByName("goid"); ok {
		return f.Offset
	}
	panic("can not find g.goid field")
}()

// runtime_g_type 变量由汇编初始化值
var runtime_g_type uint64

var gGoidOffset2 uintptr = func() uintptr { //nolint
	var iface interface{}
	type eface struct {
		_type uint64
		data  unsafe.Pointer
	}
	(*eface)(unsafe.Pointer(&iface))._type = runtime_g_type
	if f, ok := reflect.TypeOf(iface).FieldByName("goid"); ok {
		return f.Offset
	}
	panic("can not find g.goid field")
}()
