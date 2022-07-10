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
	"sync/atomic"
	"unsafe"
)

/*
	此 Cache和AtomicCache 主要用于只生成一次，永不变更/过期的缓存，二者性能相当
*/

//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin()

type Cache[K comparable, V any] struct {
	noCopy noCopy

	local     unsafe.Pointer // local fixed-size per-P pool, actual type is [P]poolLocal
	localSize uint32         // size of the local array
	New       func(K) V
}

func indexLocal[K comparable, V any](l unsafe.Pointer, i int) *map[K]V {
	size := unsafe.Sizeof([1]V{})
	lp := unsafe.Pointer(uintptr(l) + uintptr(i)*(size))
	return (*map[K]V)(lp)
}

func (c *Cache[K, V]) Get(key K) V {
	pid := runtime_procPin()
again:
	size := atomic.LoadUint32(&c.localSize)
	l := atomic.LoadPointer(&c.local) // load-consume
	if uintptr(pid) < uintptr(size) {
		m := indexLocal[K, V](l, pid) // pid 按顺序增长
		v, ok := (*m)[key]
		if !ok && c.New != nil {
			v = c.New(key)
			(*m)[key] = v
		}
		runtime_procUnpin()
		return v
	}

	// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
	sizeNew := uint32(runtime.GOMAXPROCS(0))
	localNew := make([]map[K]V, sizeNew)
	if size > 0 {
		copy(localNew[:size], (*[1 << 31]map[K]V)(l)[:])
	}
	for i := size; i < sizeNew; i++ {
		localNew[i] = make(map[K]V)
	}

	swapped := atomic.CompareAndSwapPointer(&c.local, l, unsafe.Pointer(&localNew[0]))
	if swapped {
		atomic.StoreUint32(&c.localSize, sizeNew)
	}
	goto again
}

func (c *Cache[K, V]) Set(key K, value V) {
	pid := runtime_procPin()
again:
	size := atomic.LoadUint32(&c.localSize)
	l := atomic.LoadPointer(&c.local) // load-consume
	if uintptr(pid) < uintptr(size) {
		m := indexLocal[K, V](l, pid) // pid 按顺序增长
		(*m)[key] = value

		runtime_procUnpin()
		return
	}

	// If GOMAXPROCS changes between GCs, we re-allocate the array and lose the old one.
	sizeNew := uint32(runtime.GOMAXPROCS(0))
	localNew := make([]map[K]V, sizeNew)
	if size > 0 {
		copy(localNew[:size], (*[1 << 31]map[K]V)(l)[:])
	}
	for i := size; i < sizeNew; i++ {
		localNew[i] = make(map[K]V)
	}

	swapped := atomic.CompareAndSwapPointer(&c.local, l, unsafe.Pointer(&localNew[0]))
	if swapped {
		atomic.StoreUint32(&c.localSize, sizeNew)
	}
	goto again
}

type AtomicCache[K comparable, V any] struct {
	noCopy noCopy

	cache unsafe.Pointer
	New   func(K) V
}

func (c *AtomicCache[K, V]) Get(key K) (v V) {
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

func (c *AtomicCache[K, V]) Set(key K, value V) {
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
		p := atomic.LoadPointer(&c.cache)
		if p != nil {
			cache = *(*map[K]V)(p)
		}
	}
}
