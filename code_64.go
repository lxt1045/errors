//go:build amd64 || amd64p32 || arm64
// +build amd64 amd64p32 arm64

package errors

import (
	"fmt"
	"runtime"
)

func NewCode(skip, code int, format string, a ...interface{}) (c *Code) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Code{code: code, msg: format, skip: skip}
	pcs := pool.Get().(*[DefaultDepth]uintptr)
	for i := range pcs {
		pcs[i] = 0
	}
	n := buildStack(pcs[:])
	// for i := n; i < DefaultDepth; i++ {
	// 	pcs[i] = 0
	// }

	//
	cs := cacheStack.Get(pcs, n)
	if cs == nil {
		pcs1 := make([]uintptr, DefaultDepth)
		npc1 := runtime.Callers(baseSkip, pcs1[:DefaultDepth])
		cs = &callers{}
		for _, c := range parseSlow(pcs1[:npc1]) {
			cs.stack = append(cs.stack, c.String())
		}
		l := 0
		for i, str := range cs.stack {
			// 检查是否需要转换 JSON 特殊字符串
			lStack, yes := countEscape(str)
			l += lStack
			if yes {
				cs.attr |= 1 << i
			}
		}
		cs.attr |= uint64(l) << 32

		cacheStack.Set(pcs, n, cs)
	}
	pool.Put(pcs)
	c.cache = cs
	return
}
