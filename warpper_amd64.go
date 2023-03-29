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

//go:noinline
//Wrap 需要禁止内联，因为内联后，将无法通过 BP 获取正确的PC。
func Wrap(err error, format string, ifaces ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(ifaces) > 0 {
		format = fmt.Sprintf(format, ifaces...)
	}
	return &wrapper{
		pc:  GetPC(),
		err: err,
		msg: format,
	}
}

//go:noinline
func NewLine(format string, ifaces ...interface{}) error {
	if len(ifaces) > 0 {
		format = fmt.Sprintf(format, ifaces...)
	}
	return &wrapper{
		pc:  GetPC(),
		err: nil,
		msg: format,
	}
}

//go:noinline
func Line() string {
	w := &wrapper{
		pc: GetPC(),
	}
	f := w.parse()
	return f.stack
}
