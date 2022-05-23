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
			stderrors.New("ye error")
		}},
		{"runtime.Caller", func() {
			runtime.Caller(2)
		}},
		{"runtime.Callers", func() {
			var pcs [DefaultDepth]uintptr
			runtime.Callers(3, pcs[:])
		}},
		{"pkg.New", func() {
			pkgerrs.New("ye error")
		}}, //nolint
		{"pkg.WithStack", func() {
			pkgerrs.WithStack(stderrors.New("ye error"))
		}},
		{"lxt.New", func() {
			New("ye error")
		}},
		{"lxt.NewErr", func() {
			NewErr(-1, "ye error")
		}},
		{"lxt.buildStack", func() {
			buildStack(1 + 1)
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
