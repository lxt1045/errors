package errors

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Run("Cache", func(t *testing.T) {
		cache := AtomicCache[uintptr, uintptr]{
			New: func(k uintptr) (v uintptr) {
				return k
			},
		}
		k := uintptr(100)
		v := cache.Get(100)
		assert.Equal(t, v, k)
	})
	t.Run("AtomicCache", func(t *testing.T) {
		cache := AtomicCache[uintptr, uintptr]{
			New: func(k uintptr) (v uintptr) {
				return k
			},
		}
		k := uintptr(100)
		v := cache.Get(100)
		assert.Equal(t, v, k)
	})
}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkCache)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
*/

func BenchmarkCache(b *testing.B) {
	for i := 0; i < 100; i++ {
		cache := AtomicCache[int, int]{
			New: func(k int) (v int) {
				return k
			},
		}
		atomicCache := AtomicCache[int, int]{
			New: func(k int) (v int) {
				return k
			},
		}
		N := 10240
		for i := 0; i < N; i++ {
			cache.Get(i)
			atomicCache.Get(i)
		}
		b.Run("cache", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				cache.Get(i % N)
			}
			b.StopTimer()
		})
		b.Run("AtomicCache", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				atomicCache.Get(i % N)
			}
			b.StopTimer()
		})

		m := map[int]int{}
		for i := 0; i < N; i++ {
			m[i] = i
		}
		b.Run("map", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_, ok := m[i%N]
				if !ok {
					m[i%N] = i
				}
			}
			b.StopTimer()
		})
		var lock sync.RWMutex
		b.Run("map+RWMutex", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				lock.RLock()
				_, ok := m[i%N]
				lock.RUnlock()
				if !ok {
					lock.Lock()
					m[i%N] = i
					lock.Unlock()
				}
			}
			b.StopTimer()
		})
		b.Run("RWMutex", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				lock.RLock()
				lock.RUnlock()
			}
			b.StopTimer()
		})
	}
}
