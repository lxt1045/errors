package errors

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"unsafe"

	"github.com/petermattis/goid"
	"github.com/stretchr/testify/assert"
)

func Test_OKx0(t *testing.T) {
	defer Catcher(NewGuard(), nil)
	defer func() {
		e := recover()
		if e != nil {
			t.Log("cache:", e)
		}
	}()

	err := NewCode(0, 0, "test")
	TryEscape(err)
}
func Test_OKx2(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		assert.Panics(t, func() {
			err := NewCode(0, 0, "test")
			TryEscape(err)
		})
	})
	t.Run("panic", func(t *testing.T) {
		defer Catcher(NewGuard(), nil)
		defer func() {
			e := recover()
			if e != nil {
				t.Log("cache:", e)
			}
		}()

		err := NewCode(0, 0, "test")
		TryEscape(err)
	})
	t.Run("panic", func(t *testing.T) {
		defer Catcher(NewGuard(), func(interface{}) bool {
			t.Log("sssssss")
			return true
		})
		defer func() {
			e := recover()
			if e != nil {
				t.Log("cache:", e)
			}
		}()

		err := NewCode(0, 0, "test")
		TryEscape(err)
	})
}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkTryx)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
go test -benchmem -run=^$ -bench "^(BenchmarkTryx)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
go tool pprof ./errors.test mem.prof
*/
func BenchmarkTryx(b *testing.B) {
	errCodeNotNil := NewCode(0, errCode, errMsg)
	type run struct {
		funcName string       //函数名字
		f        func() error //调用方法
	}
	var gid int64
	_ = gid

	runs := []run{
		{"goid.Get", func() (err error) {
			goid.Get()
			return
		}},
		{"TryEscape(nil)", func() (err error) {
			defer Catcher(NewGuard(), func(e interface{}) (ok bool) {
				err, ok = e.(*Code)
				return true
			})
			TryEscape(nil)
			return
		}},
		{"TryEscape(errCodeNotNil)", func() (err error) {
			defer Catcher(NewGuard(), func(e interface{}) (ok bool) {
				err, ok = e.(*Code)
				return true
			})
			TryEscape(errCodeNotNil)
			return
		}},
	}

	for _, r := range runs[:] {
		b.Run(r.funcName, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f()
			}
			b.StopTimer()
		})
	}
}

func deepCall(depth int, f func()) {
	if depth <= 0 {
		f()
		return
	}
	deepCall(depth-1, f)
}

func GetPc(f func()) uintptr {
	pc := reflect.ValueOf(f).Pointer()
	lockPCs.RLock()
	_, ok := mPCs[pc]
	lockPCs.RUnlock()
	if !ok {
		mPCs[pc] = "100"
	}
	return pc
}

var mPCs = func() map[uintptr]string {
	m := make(map[uintptr]string)
	for i := 0; i < 10000; i++ {
		m[uintptr(i)] = fmt.Sprintf("%d", i)
	}
	return m
}()
var lockPCs sync.RWMutex

var pMPCs unsafe.Pointer = func() unsafe.Pointer {
	m := make(map[uintptr]string)
	for i := 0; i < 10000; i++ {
		m[uintptr(i)] = fmt.Sprintf("%d", i)
	}
	return unsafe.Pointer(&m)
}()

func GetPc0(f func()) uintptr {
	pc := *(*uintptr)(unsafe.Pointer(&f))
	_, ok := mPCs[pc]
	if !ok {
		mPCs[pc] = "100"
	}
	return pc
}
func GetPc01(f func()) uintptr {
	pc := *(*uintptr)(unsafe.Pointer(&f))

	mPCs := *(*map[uintptr]string)(atomic.LoadPointer(&pMPCs))
	if _, ok := mPCs[pc]; !ok {
		mPCs2 := make(map[uintptr]string, len(mPCs)+10)
		mPCs2[pc] = "100"
		for {
			p := atomic.LoadPointer(&pMPCs)
			mPCs = *(*map[uintptr]string)(p)
			for k, v := range mPCs {
				mPCs2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&pMPCs, p, unsafe.Pointer(&mPCs2))
			if swapped {
				break
			}
		}
	}
	return pc
}
func GetPc1(f func(), t *testing.T) int {
	pc := reflect.ValueOf(f).Pointer()
	t.Log(runtime.FuncForPC(pc).FileLine(pc))

	p := (*uintptr)(unsafe.Pointer(&f))
	t.Logf("%d:%d", uintptr(*p), pc)
	return 0
}
func BenchmarkPC(b *testing.B) {

	b.Run("GetPc", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			GetPc(func() {})
		}
		b.StopTimer()
	})
	b.Run("GetPc0", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			GetPc0(func() {})
		}
		b.StopTimer()
	})
	b.Run("GetPc01", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			GetPc01(func() {})
		}
		b.StopTimer()
	})
	b.Run("GetPC", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			goid.GetPC()
		}
		b.StopTimer()
	})
	b.Run("FuncForPC", func(b *testing.B) {
		b.ReportAllocs()
		pc, _ := goid.GetPC(), GetPc(func() {})
		for i := 0; i < b.N; i++ {
			runtime.FuncForPC(pc).FileLine(pc)
		}
		b.StopTimer()
	})

	// debug.PrintStack()

}
