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
	"fmt"
	"sync"
)

func init() {
	// log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile) //log.Llongfile
}

var (
	mRoutineLastDefer = map[int64]struct{}{}
	lockRoutineDefer  sync.RWMutex

	tryEscapeErr func(error) // { runtime.Goexit() }
)

type guard struct {
	gid int64
	own bool
	noCopy
}

//go:noinline
func NewGuard() guard { //nolint:bgolint
	gid := Getg()
	lockRoutineDefer.Lock()
	_, ok := mRoutineLastDefer[gid]
	if !ok {
		mRoutineLastDefer[gid] = struct{}{}
	}
	lockRoutineDefer.Unlock()

	return guard{
		gid: gid,
		own: !ok,
	}
}

func Catch(g guard, f func(err interface{}) bool) { //nolint:govet
	if g.own {
		lockRoutineDefer.Lock()
		delete(mRoutineLastDefer, g.gid)
		lockRoutineDefer.Unlock()
	}
	e := recover()
	if e == nil {
		return
	}

	if f != nil && f(e) {
		return
	}
	panic(e)
}

func Check(err *Code) {
	gid := Getg()
	lockRoutineDefer.Lock()
	_, ok := mRoutineLastDefer[gid]
	lockRoutineDefer.Unlock()
	if !ok {
		cs := toCallers([]uintptr{GetPC()})
		e := fmt.Errorf("should call defer Catch(NewGuard(),func()bool before call Check(err)); file:%s",
			cs[0].File)
		if err != nil {
			e = fmt.Errorf("%w; %+v", err, e)
		}
		if tryEscapeErr != nil {
			tryEscapeErr(e)
			return
		}
		panic(e)
	}
	if err != nil {
		panic(err)
	}

	return
}
