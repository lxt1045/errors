//go:build !amd64
// +build !amd64

package errors

import (
	"fmt"
	"runtime"
	_ "runtime"
)

func NewCauseSlow(skip, code int, format string, a ...interface{}) (c *Cause) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Cause{code: code, msg: format}

	pcs := pool.Get().(*[DefaultDepth]uintptr)
	n := runtime.Callers(skip+baseSkip, pcs[:DefaultDepth-skip])

	cache := *(*stackCacheType)(atomic.LoadPointer(&mStackCache))
	cs, ok := cache[*pcs]
	if !ok {
		cs = &callers{}
		cs.stack = parseSlow(pcs[:n])
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

		// 加入缓存中
		cache2 := make(stackCacheType, len(cache)+10)
		cache2[*pcs] = cs
		for {
			p := atomic.LoadPointer(&mStackCache)
			cache = *(*stackCacheType)(p)
			for k, v := range cache {
				cache2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mStackCache, p, unsafe.Pointer(&cache2))
			if swapped {
				break
			}
		}
	}
	pool.Put(pcs)
	c.cache = cs
	return
}
