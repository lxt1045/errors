//go:build amd64
// +build amd64

package errors

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"
)

func NewCause(skip, code int, format string, a ...interface{}) (c *Cause) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Cause{code: code, msg: format}
	pcs := pool.Get().(*[DefaultDepth]uintptr)
	_ = buildStack(pcs[:])

	//
	cache := *(*stackCacheType)(atomic.LoadPointer(&mStackCache))
	cs, ok := cache[*pcs]
	if !ok {
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

type stackCacheType1 = map[string]*callers

var (
	mStackCache1 unsafe.Pointer = func() unsafe.Pointer {
		m := make(stackCacheType1)
		return unsafe.Pointer(&m)
	}()
)

func toString(p []uintptr) string {
	bs := (*[DefaultDepth * 8]byte)(unsafe.Pointer(&p[0]))[:len(p)*8]
	return *(*string)(unsafe.Pointer(&bs))
}

func NewCause2(skip, code int, format string, a ...interface{}) (c *Cause) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Cause{code: code, msg: format}
	pcs := pool.Get().(*[DefaultDepth]uintptr)
	n := buildStack(pcs[:])
	key := toString(pcs[:n])

	//
	cache := *(*stackCacheType1)(atomic.LoadPointer(&mStackCache1))
	cs, ok := cache[key]
	if !ok {
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

		// 加入缓存中
		cache2 := make(stackCacheType1, len(cache)+10)
		bsKey := []byte(key)
		cache2[string(bsKey)] = cs
		for {
			p := atomic.LoadPointer(&mStackCache1)
			cache = *(*stackCacheType1)(p)
			for k, v := range cache {
				cache2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mStackCache1, p, unsafe.Pointer(&cache2))
			if swapped {
				break
			}
		}
	}
	pool.Put(pcs)
	c.cache = cs
	return
}
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
