package errors

import (
	stderrs "errors"
	"fmt"
	"testing"

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
		fFormatting func(err error)
	}{
		{"stdNew", stdNew, func(err error) {

		}},
		{"lxtNew", lxtNew, func(err error) {

		}},
		{"pkgNew", pkgNew, func(err error) {
			return
		}},
	}
	for _, depth := range []int{1, 10, 100} {
		for _, r := range runs {
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
