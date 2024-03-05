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

func Test_stackToNewStr(t *testing.T) {
	mStr := map[string]int{}
	const N, LK = 20, 32
	key := [DefaultDepth]uintptr{}
	for i := 0; i < N; i++ {
		key[0] = uintptr(i)
		mStr[stackToNewStr(&key, LK)] = i + 88
	}
	t.Run("stackToStr", func(t *testing.T) {
		for i := 0; i < N; i++ {
			key[0] = uintptr(i)
			x, ok := mStr[stackToStr(&key, LK)]
			if !ok || x != i+88 {
				t.Fatal(i)
			}
		}
	})

}

func BenchmarkCacheDefaultDepth(b *testing.B) {
	for i := 0; i < 2; i++ {
		rcuCache := RCUCache[[DefaultDepth]uintptr, int]{
			New: func(k [DefaultDepth]uintptr) (v int) {
				return int(k[0]) + 1
			},
		}
		rcuCacheStr := RCUCache[string, int]{
			New: func(k string) (v int) {
				return int(k[0]) + 1
			},
		}
		stackCache := StackCache[int]{
			RCUCache[string, int]{
				New: func(k string) (v int) {
					return int(k[0]) + 1
				},
			},
		}
		stackCache1 := NewStackCache[int](func(k *[DefaultDepth]uintptr, l int) (v int) { return int(k[0]) + 1 })
		N, LK := 1024*4, 8
		key := [DefaultDepth]uintptr{}
		for i := 0; i < N; i++ {
			key[0] = uintptr(i)
			rcuCache.Get(key)
			stackCache.Get(&key, LK)
			stackCache1.Get(&key, LK)
			strKey := stackToNewStr(&key, LK)
			str := rcuCacheStr.Get(strKey)
			_ = str
		}
		b.Run("RCUCache", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				rcuCache.Get(key)
			}
			b.StopTimer()
		})
		b.Run("RCUCache-str", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				x := rcuCacheStr.Get(stackToStr(&key, LK))
				if x == 0 {
					b.Fatal(i)
				}
			}
			b.StopTimer()
		})
		b.Run("stackCache", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				x := stackCache.Get(&key, LK)
				if x == 0 {
					b.Fatal(i)
				}
			}
			b.StopTimer()
		})
		b.Run("stackCache1", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				x := stackCache1.Get(&key, LK)
				if x == 0 {
					b.Fatal(i)
				}
			}
			b.StopTimer()
		})

		m := map[[DefaultDepth]uintptr]int{}
		mStr := map[string]int{}
		for i := 0; i < N; i++ {
			key[0] = uintptr(i % N)
			m[key] = i
			mStr[stackToNewStr(&key, LK)] = i
		}
		b.Run("map", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				_, ok := m[key]
				if !ok {
					b.Fatal(i)
				}
			}
			b.StopTimer()
		})
		b.Run("map-str", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				strKey := stackToStr(&key, LK)
				_, ok := mStr[strKey]
				if !ok {
					b.Fatal(i)
				}
			}
			b.StopTimer()
		})
		var lock sync.RWMutex
		b.Run("map+RWMutex", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				lock.RLock()
				_, ok := m[key]
				lock.RUnlock()
				if !ok {
					lock.Lock()
					m[key] = i
					lock.Unlock()
				}
			}
			b.StopTimer()
		})
		b.Run("map-str+RWMutex", func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				key[0] = uintptr(i % N)
				strKey := stackToStr(&key, LK)
				lock.RLock()
				_, ok := mStr[strKey]
				lock.RUnlock()
				if !ok {
					lock.Lock()
					mStr[strKey] = i
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
