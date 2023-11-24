package errors

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	t.Run("RCUCache", func(t *testing.T) {
		cache := RCUCache[uintptr, uintptr]{
			New: func(k uintptr) (v uintptr) {
				return k
			},
		}
		k := uintptr(100)
		v := cache.Get(100)
		assert.Equal(t, v, k)
	})

	cache := RCUCache[int, int]{
		New: nil,
	}
	t.Run("RCUCache-Set1", func(t *testing.T) {
		f := func(wg *sync.WaitGroup) {
			defer wg.Done()
			for i := 0; i < 1000; i++ {
				cache.Set(i, i*100)
			}
		}

		wg := &sync.WaitGroup{}
		wg.Add(2)
		go f(wg)
		go f(wg)
		wg.Wait()
	})

	t.Run("RCUCache-Set100", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			cache.Set(i, i*100)
		}

		for i := 0; i < 100; i++ {
			v := cache.Get(i)
			assert.Equal(t, v, i*100)
		}
	})

}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkCache)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
*/

func BenchmarkCache(b *testing.B) {
	for i := 0; i < 3; i++ {
		rcuCache := RCUCache[int, int]{
			New: func(k int) (v int) {
				return k
			},
		}
		N := 10240
		for i := 0; i < N; i++ {
			rcuCache.Get(i)
		}
		b.Run("RCUCache", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				rcuCache.Get(i % N)
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
