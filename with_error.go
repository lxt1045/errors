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

package errors

import (
	"reflect"
	"sync"
)

var (
	mErrFunc sync.Map

	_ = func() error {
		Register(&Code{}, func(err error) string {
			return err.(*Code).msg
		})
		return nil
	}()
)

func Register(e error, f func(err error) string) (err error) {
	fOld := getErrorFunc(e)
	if fOld != nil {
		err = NewCode(1, 0, "error type already registered")
		return
	}
	mErrFunc.Store(errKey(e), f)
	return
}

func errKey(err error) string {
	typ := reflect.TypeOf(err)
	k := typ.String() + "/" + typ.Name()
	return k
}

func getErrorFunc(err error) (f func(err error) string) {
	v, ok := mErrFunc.Load(errKey(err))
	if !ok {
		return
	}
	f = v.(func(err error) string)
	return
}
