package errors

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"

	"github.com/petermattis/goid"
	"github.com/stretchr/testify/assert"
)

// "github.com/petermattis/goid"
// https://github.com/cch123/goroutineid

func Test_Assertx_3(t *testing.T) {
	defer func() {
		err := recover()
		t.Log("1:", err)
	}()
	for i := 0; i < 2; i++ {
		if i > 0 {
			panic("xxx")
		}
		defer func() {
			err := recover()
			t.Log("1:", err)
		}()
	}

	t.Log("NumGoroutine:", runtime.NumGoroutine())
}

func Test_Assertx_2(t *testing.T) {
	c := make(chan int)
	go func() {
		defer func() {
			t.Log(New("runtime.Goexit()").Error())
			<-c
		}()
		c <- 1
		runtime.Goexit()
		panic("xxx")
	}()
	<-c
	t.Log("NumGoroutine:", runtime.NumGoroutine())
	c <- 1
	time.Sleep(time.Second + 3)
	t.Log("NumGoroutine:", runtime.NumGoroutine())

}

func Test_Assertx_1(t *testing.T) {
	err := NewCause(0, errCode, errMsg)
	hookBeforePanic = func(ctx context.Context, err error) {
		t.Fatal(err)
	}
	t.Run("!OKx.2", func(t *testing.T) {
		hookBeforePanic = func(ctx context.Context, err error) {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			if i == 1 {
				OKx(context.Background(), false, err)
			}
			defer TryCatchx(func(e interface{}) (ok bool) {
				err1 := e.(*Cause)
				assert.Equal(t, err.code, err1.code)
				assert.Equal(t, err.msg, err1.msg)
				return true
			})()
		}
	})
	for i := 0; i < 2; i++ {
		t.Run("!OKx.0", func(t *testing.T) {
			func() {
				func() {
					func() {
						defer TryCatchx(func(e interface{}) (ok bool) {
							err1 := e.(*Cause)
							assert.Equal(t, err.code, err1.code)
							assert.Equal(t, err.msg, err1.msg)
							return true
						})()
						deepFunc(0, func() {
							OKx(context.Background(), false, err)
						})
					}()
				}()
			}()

		})
	}

	t.Run("!OKx.1", func(t *testing.T) {
		hookBeforePanic = func(ctx context.Context, err error) {
			t.Log(err)
		}
		assert.Panics(t, func() {
			OKx(context.Background(), false, nil)
			defer TryCatchx(func(e interface{}) (ok bool) {
				err1 := e.(*Cause)
				assert.Equal(t, err.code, err1.code)
				assert.Equal(t, err.msg, err1.msg)
				return true
			})()
		})
	})

}
func deepFunc(deep int, f func()) {
	if deep > 0 {
		deepFunc(deep-1, f)
	}
	f()
}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkTryx)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
go test -benchmem -run=^$ -bench "^(BenchmarkTryx)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
go tool pprof ./errors.test mem.prof
*/
func BenchmarkTryx(b *testing.B) {
	hookBeforePanic = func(ctx context.Context, err error) {}

	errCauseNotNil := NewCause(0, errCode, errMsg)
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
		{"goid.GetDefer", func() (err error) {
			goid.GetDefer()
			return
		}},
		{"TryCatchx.true", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				err, ok = e.(*Cause)
				return true
			})()
			OKx(context.Background(), true, nil)
			return
		}},
		{"TryCatchx.false.nil", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			OKx(context.Background(), false, nil)
			return
		}},
		{"TryCatchx.false.not-nil", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			OKx(context.Background(), false, errCauseNotNil)
			return
		}},
		{"TryCatchx.deepCall.3", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(3, func() { OKx(context.Background(), false, errCauseNotNil) })
			return
		}},
		{"TryCatchx.deepCall.10", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(10, func() { OKx(context.Background(), false, errCauseNotNil) })
			return
		}},
		{"TryCatchx.deepCall.30", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(30, func() { OKx(context.Background(), false, errCauseNotNil) })
			return
		}},
		{"TryByFunc.true", func() (err error) {
			defer TryByFunc(func(e interface{}) (ok bool) {
				return true
			})
			OK(true, errCauseNotNil)
			return
		}},
		{"TryByFunc.false.not-nil", func() (err error) {
			defer TryByFunc(func(e interface{}) (ok bool) {
				return true
			})
			OK(false, errCauseNotNil)
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

func Test_PC(t *testing.T) {
	func() {
		func() {
			var pc uintptr
			for i := 0; i < 10; i++ {
				func() {
					func() {
						pc, _ = goid.GetPC(), GetPc1(func() {}, t)
					}()
				}()

				t.Log(runtime.FuncForPC(pc).FileLine(pc))
				t.Log(runtime.FuncForPC(pc).Name())
				t.Log("pc:", pc)

				pc2, bp := goid.GetPCBP()
				t.Logf("pc:%d, bp:%d\n", pc2, bp)
				f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
				t.Logf("%+v\n", f)
			}
			goid.Ret(nil)
		}()
		t.Log("yyy")
	}()
	t.Log("xxx")
}

func Test_PC1(t *testing.T) {
	err := func() (err error) {
		var y uintptr
		err = func() (err error) {
			goid.Print()
			// t.Log(goid.Getcallerpc())
			goid.Ret(err)
			err = errors.New("test")
			goid.Ret(err)
			t.Log("111")
			return
		}()
		atomic.LoadUintptr(&y)
		t.Log("yyy")
		return
	}()
	t.Log("xxx")
	_ = err
	t.Log(err)
}
func Test_NewLog(t *testing.T) {
	l := goid.NewLog(0, "")

	t.Logf("xxx:%+v", l)

	pc := goid.GetPC()
	t.Log(runtime.FuncForPC(pc).FileLine(pc))
	t.Log(runtime.FuncForPC(pc).Name())
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

func TestGetDeferPC(t *testing.T) {
	GetDeferPC()
}

func GetDeferPC() {
	defer func() {

	}()
	var pc [1]uintptr
	runtime.Callers(0, pc[:])

	f := runtime.FuncForPC(pc[0])
	fRaw := (*_func)(unsafe.Pointer(f))
	fmt.Printf("%+v\n", fRaw)
}

type funcFlag uint8
type funcID uint8
type _func struct {
	entryoff uint32 // start pc, as offset from moduledata.text/pcHeader.textStart
	nameoff  int32  // function name

	args        int32  // in/out args size
	deferreturn uint32 // offset of start of a deferreturn call instruction from entry, if any.

	pcsp      uint32
	pcfile    uint32
	pcln      uint32
	npcdata   uint32
	cuOffset  uint32 // runtime.cutab offset of this function's CU
	funcID    funcID // set for certain special runtime functions
	flag      funcFlag
	_         [1]byte // pad
	nfuncdata uint8   // must be last, must end on a uint32-aligned boundary
}
