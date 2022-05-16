package errors

import (
	"encoding/json"
	stderrors "errors"
	"fmt"
	"runtime"
	"strconv"
	"testing"

	pkgerrors "github.com/pkg/errors"
)

func Test_pkgWithStack(t *testing.T) {
	err := pkgerrors.WithStack(stderrors.New("ye error"))
	fmt.Printf("err:%+v", err)
}

// go test -benchmem -run=^$ -bench "^(BenchmarkLxtNew)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
// go tool pprof ./errors.test cpu.prof
// go test -benchmem -run=^$ -bench "^(BenchmarkLxtNew)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
// go tool pprof ./errors.test mem.prof
// web
func BenchmarkLxtNew(b *testing.B) {
	b.Run("NewErrSkip", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			NewErr(-1, "ye error") //nolint
		}
		b.StopTimer()
	})
}

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
			pkgerrors.New("ye error")
		}}, //nolint
		{"pkg.WithStack", func() {
			pkgerrors.WithStack(stderrors.New("ye error"))
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

func stdNew(depth int) error {
	if depth <= 0 {
		return stderrors.New("no error")
	}
	return stderrors.New(stdNew(depth - 1).Error())
}

func pkgNew(depth int) error {
	if depth <= 0 {
		return pkgerrors.New("ye error")
	}
	// pkgerrors.Wrap(pkgNew(depth - 1))
	return pkgerrors.WithStack(pkgNew(depth - 1))
}

func lxtNew(depth int) error {
	if depth <= 0 {
		return New("ye error")
	}
	return Wrap(lxtNew(depth-1), "")
}

func BenchmarkWarp(b *testing.B) {
	type run struct {
		funcName string                //函数名字
		f        func(depth int) error //调用方法
	}

	runs := []run{
		{"stdNew", stdNew},
		{"pkgNew", pkgNew},
		{"lxtNew", lxtNew},
	}
	depths := []int{0, 1, 5, 10} //嵌套深度
	for _, depth := range depths {
		for _, r := range runs {
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

func BenchmarkErrorFormat(b *testing.B) {
	b.Run("FE3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := FE3().(*Err)
			_ = err
		}
		b.StopTimer()
	})

	b.Run("Error", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			err := FE3().(*Err)
			err.Error() //nolint
		}
		b.StopTimer()
	})

}

func BenchmarkFormatting(b *testing.B) {
	N := 1024
	depths := []int{1, 5, 10} //嵌套深度
	mErrCache := make(map[string][]error)
	for _, depth := range depths {
		pkgErrs := make([]error, N)
		lxtErrs := make([]error, N)
		for i := range pkgErrs {
			pkgErrs[i] = pkgNew(depth)
			lxtErrs[i] = lxtNew(depth)
		}
		str := strconv.Itoa(depth)
		mErrCache["pkg."+str] = pkgErrs
		mErrCache["lxt."+str] = pkgErrs
	}

	for _, depth := range depths {
		name := fmt.Sprintf("%s-%d", "pkgNew.%v", depth)
		b.Run(name, func(b *testing.B) {
			errs := mErrCache["pkg."+strconv.Itoa(depth)]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fmt.Sprintf("%v", errs[i%len(errs)])
			}
			b.StopTimer()
		})

		name = fmt.Sprintf("%s-%d", "pkgNew.%+v", depth)
		b.Run(name, func(b *testing.B) {
			errs := mErrCache["pkg."+strconv.Itoa(depth)]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fmt.Sprintf("%+v", errs[i%len(errs)])
			}
			b.StopTimer()
		})

		name = fmt.Sprintf("%s-%d", "lxtNew.text", depth)
		b.Run(name, func(b *testing.B) {
			errs := mErrCache["lxt."+strconv.Itoa(depth)]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fmt.Sprintf("%v", errs[i%len(errs)])
			}
			b.StopTimer()
		})

		name = fmt.Sprintf("%s-%d", "lxtNew.json", depth)
		b.Run(name, func(b *testing.B) {
			errs := mErrCache["lxt."+strconv.Itoa(depth)]
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				json.Marshal(errs[i%len(errs)]) //nolint
			}
			b.StopTimer()
		})
	}
}

func BenchmarkNewAndFormatting(b *testing.B) {
	type run struct {
		funcName string                //函数名字
		f        func(depth int) error //调用方法
	}

	depths := []int{1, 5, 10} //嵌套深度
	formats := []string{
		"%s", "%+v",
	}
	for _, format := range formats {
		for _, depth := range depths {
			name := fmt.Sprintf("%s-%d-%s", "pkgNew", depth, format)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					fmt.Sprintf(format, pkgNew(depth))
				}
				b.StopTimer()
			})

			name = fmt.Sprintf("%s-%d-%s", "lxtNew", depth, format)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					fmt.Sprintf(format, lxtNew(depth))
				}
				b.StopTimer()
			})
		}
	}
}

func FE1() error {
	err := NewErr(1600002, "card status error 2")
	return err
}

func FE2() error {
	err := FE1()
	err = Wrap(err, "error log 2")
	return err
}

func FE3() error {
	err := FE2()
	err = Wrap(err, "error log 3")
	return err
}

func FE4() error {
	err := fmt.Errorf("std error")
	return err
}

func FE5() error {
	err := FE4()
	err = Wrap(err, "error log 5")
	return err
}
