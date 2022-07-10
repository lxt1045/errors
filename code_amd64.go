//go:build amd64
// +build amd64

package errors

import (
	"fmt"
	"runtime"
)

func NewCode(skip, code int, format string, a ...interface{}) (c *Code) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Code{code: code, msg: format}
	pcs := pool.Get().(*[DefaultDepth]uintptr)
	n := buildStack(pcs[:])
	key := toString(pcs[:n])

	//
	cs := cacheStack.Get(key)
	if cs == nil {
		pcs1 := make([]uintptr, DefaultDepth)
		npc1 := runtime.Callers(skip+baseSkip, pcs1[:DefaultDepth-skip])
		cs = &callers{}
		cs.stack = parseSlow(pcs1[:npc1])
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

		cacheStack.Set(key, cs)
	}
	pool.Put(pcs)
	c.cache = cs
	return
}
