package errors

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	pkgerrs "github.com/pkg/errors"
)

func BenchmarkWrap(b *testing.B) {
	runs := []struct {
		funcName string                //函数名字
		f        func(depth int) error //调用方法
	}{
		{"stdWrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = fmt.Errorf("%w", errors.New(errTrace))
			}
			return err
		}},
		{"Wrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = Wrap(err, errTrace)
			}
			return err
		}},
		{"pkg.Wrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = pkgerrs.Wrap(err, errTrace)
			}
			return err
		}},
	}
	depths := []int{1, 10} //嵌套深度
	for _, r := range runs {
		for _, depth := range depths {
			name := fmt.Sprintf("%s-%d", r.funcName, depth)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r.f(depth) //nolint
				}
				b.StopTimer()
			})
		}
	}
}

type M struct {
	sync.RWMutex
	m map[int]int
}

func NewM(n int) (m *M) {
	m = &M{
		m: make(map[int]int, n),
	}
	for i := 0; i < n; i++ {
		m.m[i] = i
	}
	return
}

func (m *M) Get(k int) int {
	m.RLock()
	defer m.RUnlock()
	return m.m[k]
}

type SyncM struct {
	sync.Map
}

func NewSyncM(n int) (m *SyncM) {
	m = &SyncM{}
	for i := 0; i < n; i++ {
		m.Store(i, i)
	}
	return
}

func (m *SyncM) Get(k int) int {
	v, _ := m.Load(k)
	n, _ := v.(int)
	return n
}

func BenchmarkMap(b *testing.B) {
	N := 6
	for x := 0; x < N; x++ {
		n := 10
		for i := 0; i < x; i++ {
			n = 10 * n
		}
		stdMap, syncMap := NewM(n), NewSyncM(n)

		b.Run(fmt.Sprintf("syncMap-%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				k := i % n
				syncMap.Get(k)
			}
			b.StopTimer()
		})
		b.Run(fmt.Sprintf("stdMap-%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				k := i % n
				stdMap.Get(k)
			}
			b.StopTimer()
		})

	}
}

func BenchmarkMapAccess(b *testing.B) {
	mLens := []int{100, 1000, 10000, 100000}
	value := []string{
		"github.com/lxt1045/errors/warpper_test.go",
		"github.com/lxt1045/errors/warpper_test.go",
		"github.com/lxt1045/errors/warpper_test.go",
	}
	type vt struct {
		K [DefaultDepth]uintptr
		V []string
	}
	for _, l := range mLens {
		mStacks := make(map[[DefaultDepth]uintptr][]string)
		mStacks2 := make(map[uintptr]*vt)
		keys := make([][DefaultDepth]uintptr, l)
		for i := 0; i < l; i++ {
			key := [DefaultDepth]uintptr{}
			for j := range key {
				key[j] = uintptr(i)
			}
			keys[i] = key
			mStacks[key] = value
			mStacks2[uintptr(i)] = &vt{K: key, V: value}
		}
		b.Run(fmt.Sprintf("1-%d", l), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v := mStacks[keys[i%l]]
				_ = v
			}
			b.StopTimer()
		})
		b.Run(fmt.Sprintf("2-%d", l), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				key := keys[i%l]
				v := mStacks2[key[0]]
				if v.K == key {
					_ = v.V
				}
			}
			b.StopTimer()
		})
	}
}

func Args(x int, y bool, z int) int {
	if z > 10 {
		return z
	}
	return z
}
func Args2(x int, y bool, zs ...int) int {
	if len(zs) > 0 && zs[0] > 10 {
		return zs[0]
	}
	return x
}

func Args3(zs ...int) int {
	if len(zs) > 0 {
		return zs[0]
	}
	return 1
}

func BenchmarkArgs(b *testing.B) {
	b.Run("args1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100; i++ {
				Args(i, true, i)
			}
		}
		b.StopTimer()
	})
	b.Run("args2", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100; i++ {
				Args2(i, true, i)
			}
		}
		b.StopTimer()
	})
	b.Run("args2.1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100; i++ {
				Args2(i, true)
			}
		}
		b.StopTimer()
	})
	b.Run("args3", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100; i++ {
				Args3(i)
			}
		}
		b.StopTimer()
	})
	b.Run("args3.1", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for i := 0; i < 100; i++ {
				Args3()
			}
		}
		b.StopTimer()
	})
}
