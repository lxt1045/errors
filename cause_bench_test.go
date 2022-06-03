package errors

import (
	stderrors "errors"
	"runtime"
	"testing"

	pkgerrs "github.com/pkg/errors"
)

//
func BenchmarkNew(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"std.New", func() {
			_ = stderrors.New("ye error")
		}},
		{"runtime.Caller", func() {
			runtime.Caller(2)
		}},
		{"runtime.Callers", func() {
			var pcs [DefaultDepth]uintptr
			runtime.Callers(3, pcs[:])
		}},
		{"pkg.New", func() {
			_ = pkgerrs.New("ye error")
		}},
		{"pkg.WithStack", func() {
			_ = pkgerrs.WithStack(stderrors.New("ye error"))
		}},
		{"lxt.New", func() {
			_ = New("ye error")
		}},
		{"lxt.NewErr", func() {
			_ = NewErr(-1, "ye error")
		}},
	}
	for _, r := range runs {
		name := r.funcName
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f() //nolint
			}
			b.StopTimer()
		})
	}
}
