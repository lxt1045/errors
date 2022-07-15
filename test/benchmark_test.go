package errors

import (
	stderrs "errors"
	"fmt"
	"testing"
	"unsafe"

	lxterrs "github.com/lxt1045/errors"
	pkgerrs "github.com/pkg/errors"
)

func stdNew(at, depth int) error {
	if at >= depth {
		return stderrs.New("no error")
	}
	return stdNew(at+1, depth)
}

func pkgNew(at, depth int) error {
	if at >= depth {
		return pkgerrs.New("ye error")
	}
	return pkgNew(at+1, depth)
}

func lxtNew(at, depth int) error {
	if at >= depth {
		return lxterrs.New("ye error")
	}
	return lxtNew(at+1, depth)
}

func NewCode(at, depth int) error {
	if at >= depth {
		return lxterrs.NewCode(0, 0, "ye error")
	}
	return NewCode(at+1, depth)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE string

//
func BenchmarkNew(b *testing.B) {
	runs := []struct {
		funcName string
		f        func(at, depth int) error
	}{
		{"stdNew", func(at, depth int) error {
			return stdNew(at, depth)
		}},
		{"lxtNew", func(at, depth int) error {
			return lxtNew(at, depth)
		}},
		{"NewCode", func(at, depth int) error {
			return NewCode(at, depth)
		}},
		{"pkgNew", func(at, depth int) error {
			return pkgNew(at, depth)
		}},
	}
	for _, depth := range []int{1, 10, 100} {
		for _, r := range runs {
			err := stderrs.New("")
			name := fmt.Sprintf("%s-%d", r.funcName, depth)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r.f(0, depth)
				}
				b.StopTimer()
				GlobalE = err.Error()
			})
		}
	}
}

func BenchmarkFormatting(b *testing.B) {
	runs := []struct {
		funcName    string
		fNew        func(at, depth int) error
		fFormatting func(err error) string
	}{
		{"std.%+v", stdNew, func(err error) string {
			return fmt.Sprintf("%+v", err)
		}},
		{"lxt.%+v", lxtNew, func(err error) string {
			return fmt.Sprintf("%+v", err)
		}},
		{"lxt.Json", lxtNew, func(err error) string {
			bs := lxterrs.MarshalJSON(err)
			return *(*string)(unsafe.Pointer(&bs))
		}},
		{"pkg.%+v", pkgNew, func(err error) string {
			return fmt.Sprintf("%+v", err)
		}},
	}
	for _, r := range runs {
		for _, depth := range []int{1, 10, 100} {
			err := r.fNew(0, depth)
			name := fmt.Sprintf("%s-%d", r.funcName, depth)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r.fFormatting(err)
				}
				b.StopTimer()
				GlobalE = err.Error()
			})
		}
	}
}

//
func deepCall(depth int, f func()) {
	if depth <= 0 {
		f()
		return
	}
	deepCall(depth-1, f)
}
func BenchmarkPkg(b *testing.B) {
	b.Run("pkg/errors", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		deepCall(10, func() {
			for i := 0; i < b.N; i++ {
				err = pkgerrs.New("ye error")
				GlobalE = fmt.Sprintf("%+v", err)
			}
			b.StopTimer()
		})
	})
	b.Run("stderrors-Sprintf", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		deepCall(10, func() {
			for i := 0; i < b.N; i++ {
				err = stderrs.New("ye error")
				GlobalE = fmt.Sprintf("%+v", err)
			}
			b.StopTimer()
		})
	})
	b.Run("stderrors", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		deepCall(10, func() {
			for i := 0; i < b.N; i++ {
				err = stderrs.New("ye error")
				GlobalE = err.Error()
			}
			b.StopTimer()
		})
	})
}
