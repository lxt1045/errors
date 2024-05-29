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
	"fmt"
	"strconv"
	_ "unsafe" //nolint:bgolint

	"github.com/lxt1045/errors/g"
)

/*
SetJMP 还需要备份 g,_defer 链，在LongJMP的时候恢复，主打一个debug和release一致性
*/

func Setjmp() (jump, error) //nolint:bgolint

func longjmp(jmp jump, err error) uintptr

var defer_offset = g.G__defer_offset

type jump struct {
	pc     uintptr //nolint:unused
	parent uintptr //nolint:unused
	_defer uintptr //nolint:unused

	// noCopy noCopy //nolint:unused
}

//go:noinline
func (t jump) Longjmp(err error) {
	if err == nil {
		return
	}
	//还是要加上检查，否则报错信息太难看
	// 但是检查时只要检查 更上一级的 PC 是否相等即可，不需要复杂的 map 存储了！！！
	parent := longjmp(t, err)
	if parent != t.parent {
		cs := toCallers([]uintptr{parent, t.parent, GetPC()})
		e := fmt.Errorf("handler.Check() should be called in [%s] not in [%s]; file:%s",
			cs[1].Func, cs[0].Func, cs[2].File+":"+strconv.Itoa(cs[2].Line))

		if tryHandlerErr != nil {
			tryHandlerErr(e)
			return
		}
		panic(e)
	}
}

// func GetDefer() unsafe.Pointer
// func getgi() interface{}
// func getdeferi() interface{}

// var g__defer_offset uintptr = func() uintptr {
// 	g := getgi()
// 	if f, ok := reflect.TypeOf(g).FieldByName("_defer"); ok {
// 		return f.Offset
// 	}
// 	panic("can not find g.goid field")
// }()

// var _defer_link_offset uintptr = func() uintptr {
// 	_defer := getdeferi()
// 	if f, ok := reflect.TypeOf(_defer).FieldByName("link"); ok {
// 		return f.Offset
// 	}
// 	panic("can not find g.goid field")
// }()

// type _defer struct {
// 	started   bool
// 	heap      bool
// 	openDefer bool
// 	sp        uintptr // sp at time of defer
// 	pc        uintptr // pc at time of defer
// 	fn        func()  // can be nil for open-coded defers
// 	_panic    uintptr
// 	link      *_defer // next defer on G; can point to either heap or stack!

// 	fd      unsafe.Pointer
// 	varp    uintptr
// 	framepc uintptr
// }

// func (d *_defer) Next() *_defer {
// 	return d.link
// }
// func (d *_defer) PC() uintptr {
// 	return d.pc
// }
