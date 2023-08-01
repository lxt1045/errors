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
	"runtime"
	_ "unsafe" //nolint:bgolint

	"github.com/rs/zerolog"
)

func getPCSlow() (pcs [1]uintptr) {
	runtime.Callers(3, pcs[:])
	return
}

func buildStackSlow(s []uintptr) int {
	return runtime.Callers(3, s[:])
}

func CallersSkip(skip int) (cs []caller) {
	pcs := pool.Get().(*[DefaultDepth]uintptr)
	for i := range pcs {
		pcs[i] = 0
	}
	n := buildStack(pcs[:]) //仅当特征码使用，有点大材小用
	for i := n; i < DefaultDepth; i++ {
		pcs[i] = 0
	}

	//
	cs = cacheCallers.Get(*pcs)
	if cs == nil {
		pcs1 := make([]uintptr, DefaultDepth)
		npc1 := runtime.Callers(baseSkip, pcs1[:DefaultDepth])
		cs = parseSlow(pcs1[:npc1])

		cacheCallers.Set(*pcs, cs)
	}
	pool.Put(pcs)
	cs = cs[skip:]

	return
}

// CallerFrame 使用 Read-copy update(RCU) 缓存提高性能
func CallerFrame(l uintptr) (c *caller) {
	c = cacheCaller.Get(l)
	if c != nil {
		return
	}

	cs := parseSlow([]uintptr{l})
	if len(cs) > 0 {
		c = &cs[0]
		cacheCaller.Set(l, c)
	}
	return
}

type zeroStack struct {
	stack []caller
	skip  int
}

func (e zeroStack) MarshalZerologArray(a *zerolog.Array) {
	if len(e.stack) > e.skip {
		for _, c := range e.stack[e.skip:] {
			a.Str(c.String())
		}
	}
}

func ZerologStack(skip int) zeroStack {
	cs := CallersSkip(skip + 1)
	return zeroStack{stack: cs, skip: skip}
}
