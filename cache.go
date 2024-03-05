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
	"sync"
	"sync/atomic"
	"unsafe"
)

/*
	此 Cache和AtomicCache 主要用于只生成一次，永不变更/过期的缓存，二者性能相当
*/

type RCUCache[K comparable, V any] struct {
	noCopy noCopy //nolint:unused

	cache unsafe.Pointer
	New   func(K) V
}

func (c *RCUCache[K, V]) JustGet(key K) (v V, ok bool) {
	var cache map[K]V

	if p := atomic.LoadPointer(&c.cache); p != nil {
		cache = *(*map[K]V)(p)
		v, ok = cache[key]
	}

	return
}

func (c *RCUCache[K, V]) Get(key K) (v V) {
	var ok bool
	var cache map[K]V

	if p := atomic.LoadPointer(&c.cache); p != nil {
		cache = *(*map[K]V)(p)
		v, ok = cache[key]
		if ok {
			return v
		}
	}
	if c.New == nil {
		return
	}
	v = c.New(key)
	cacheNew := make(map[K]V, len(cache)+8)
	cacheNew[key] = v
	for {
		p := atomic.LoadPointer(&c.cache)
		if p != nil {
			cache = *(*map[K]V)(p)
		}
		for k, v := range cache {
			cacheNew[k] = v
		}
		swapped := atomic.CompareAndSwapPointer(&c.cache, p, unsafe.Pointer(&cacheNew))
		if swapped {
			break
		}
	}

	return v
}

func (c *RCUCache[K, V]) Set(key K, value V) {
	var cache map[K]V

	p := atomic.LoadPointer(&c.cache)
	if p != nil {
		cache = *(*map[K]V)(p)
	}

	cacheNew := make(map[K]V, len(cache)+8)
	cacheNew[key] = value
	for {
		for k, v := range cache {
			cacheNew[k] = v
		}
		swapped := atomic.CompareAndSwapPointer(&c.cache, p, unsafe.Pointer(&cacheNew))
		if swapped {
			break
		}
		p = atomic.LoadPointer(&c.cache)
		if p != nil {
			cache = *(*map[K]V)(p)
		}
	}
}

type StackCache[V any] struct {
	RCUCache[string, V]
}

func (c *StackCache[V]) Get(s *[DefaultDepth]uintptr, l int) (v V) {
	v, ok := c.JustGet(stackToStr(s, l))
	if ok || c.New == nil {
		return
	}
	return c.RCUCache.Get(stackToNewStr(s, l))
}

func (c *StackCache[V]) Set(s *[DefaultDepth]uintptr, l int, v V) {
	c.RCUCache.Set(stackToNewStr(s, l), v)
}

type stackCache1[V any] struct {
	noCopy noCopy //nolint:unused
	sync.RWMutex
	m    map[string]V
	fNew func(s *[DefaultDepth]uintptr, l int) V
}

func NewStackCache[V any](fNew func(s *[DefaultDepth]uintptr, l int) V) *stackCache1[V] {
	return &stackCache1[V]{
		m:    make(map[string]V),
		fNew: fNew,
	}
}

func (c *stackCache1[V]) Get(s *[DefaultDepth]uintptr, l int) (v V) {
	c.RLock()
	v, ok := c.m[stackToStr(s, l)]
	c.RUnlock()
	if ok {
		return
	}
	v = c.fNew(s, l)
	c.Set(s, l, v)
	return
}

func (c *stackCache1[V]) Set(s *[DefaultDepth]uintptr, l int, v V) {
	c.Lock()
	defer c.Unlock()
	c.m[stackToNewStr(s, l)] = v
}

func stackToNewStr(s *[DefaultDepth]uintptr, l int) string {
	p := (*byte)(unsafe.Pointer(s))
	return string(unsafe.Slice(p, l*8))
}
func stackToStr(s *[DefaultDepth]uintptr, l int) string {
	p := (*byte)(unsafe.Pointer(s))
	return unsafe.String(p, l*8)
}
