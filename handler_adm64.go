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
	_ "unsafe" //nolint:bgolint
)

var tryHandlerErr func(error)

func NewHandler() (handler, error) //nolint:bgolint

func tryJump(pc, parent uintptr, err error) uintptr

type handler struct {
	pc     uintptr
	parent uintptr
}

//go:noinline
func (t handler) Check(err error) {
	//还是要加上检查，否则报错信息太难看
	// 但是检查时只要检查 更上一级的 PC 是否相等即可，不需要复杂的 map 存储了！！！
	parent := tryJump(t.pc, t.parent, err)
	if parent != t.parent {
		cs := toCallers([]uintptr{parent, t.parent, GetPC()})
		e := fmt.Errorf("handler.Check() should be called in [%s] not in [%s]; file:%s",
			cs[1].Func, cs[0].Func, cs[2].File)
		if err != nil {
			e = fmt.Errorf("%w; %+v", err, e)
		}
		if tryHandlerErr != nil {
			tryHandlerErr(e)
			return
		}
		panic(e)
	}
}
