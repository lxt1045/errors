package errors

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
	err := NewErrSkip(0, errCode, errMsg)
	notRecoverPanicHook = func(ctx context.Context, err error) {
		t.Fatal(err)
	}
	t.Run("!OKx.2", func(t *testing.T) {
		notRecoverPanicHook = func(ctx context.Context, err error) {
			t.Fatal(err)
		}
		for i := 0; i < 2; i++ {
			if i == 1 {
				OKx(false, err)
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
							OKx(false, err)
						})
					}()
				}()
			}()

		})
	}

	t.Run("!OKx.1", func(t *testing.T) {
		notRecoverPanicHook = func(ctx context.Context, err error) {
			t.Log(err)
		}
		assert.Panics(t, func() {
			OKx(false, nil)
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
	notRecoverPanicHook = func(ctx context.Context, err error) {
	}

	errCauseNotNil := NewErrSkip(0, errCode, errMsg)
	type run struct {
		funcName string       //函数名字
		f        func() error //调用方法
	}
	var pdp *deferPoint
	var gid int64
	_, _ = pdp, gid

	runs := []run{
		{"deferPoint", func() (err error) {
			var dp deferPoint
			runtime.Callers(1+baseSkip, dp[:])
			return
		}},
		{"TryCatchx.true", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				err, ok = e.(*Cause)
				return true
			})()
			OKx(true, nil)
			return
		}},
		{"TryCatchx.false.nil", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			OKx(false, nil)
			return
		}},
		{"TryCatchx.false.not-nil", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			OKx(false, errCauseNotNil)
			return
		}},
		{"TryCatchx.deepCall.3", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(3, func() { OKx(false, errCauseNotNil) })
			return
		}},
		{"TryCatchx.deepCall.10", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(10, func() { OKx(false, errCauseNotNil) })
			return
		}},
		{"TryCatchx.deepCall.30", func() (err error) {
			defer TryCatchx(func(e interface{}) (ok bool) {
				return true
			})()
			deepCall(30, func() { OKx(false, errCauseNotNil) })
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
