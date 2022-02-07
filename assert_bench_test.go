package errors

import (
	"errors"
	"fmt"
	"testing"
)

// go test -benchmem -run=^$ -bench "^(BenchmarkTryCatch)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
// go tool pprof ./errors.test cpu.prof
// web
func BenchmarkTryCatch(b *testing.B) {
	bizErr1 := NewErr(11, "msg1", "trace1")
	bizErr := NewErr(88, "msg", "trace")
	fNilErr := func() error {
		return nil
	}
	fSimpleErr := func() error {
		return errors.New("error")
	}
	fBizErr := func() error {
		return bizErr1
	}
	fErrors := []struct {
		funcName string       //函数名字
		f        func() error //调用方法
	}{
		{"fNilErr", fNilErr},
		{"fSimpleErr", fSimpleErr},
		{"fBizErr", fBizErr},
	}
	// fErrors=fErrors[2:]
	ifErr := func(fErr func() error) error {
		err := fErr()
		if err == nil {
			return bizErr
		}
		return nil
	}
	tryCatch := func(fErr func() error) (errRet error) {
		defer TryCatch(func(err BizErr) {
			errRet = bizErr
		})()
		err := fErr()
		MustNil(err, bizErr)
		return
	}
	tryCatchErr := func(fErr func() error) (errRet error) {
		defer TryCatchErr(&errRet)()
		err := fErr()
		MustNil(err, bizErr)
		return
	}

	type run struct {
		funcName string                        //函数名字
		f        func(fErr func() error) error //调用方法
	}

	runs := []run{
		{"ifErr", ifErr},
		{"tryCatch", tryCatch},
		{"tryCatchErr", tryCatchErr},
	}
	// runs = runs[2:]
	for _, fError := range fErrors {
		for _, r := range runs {
			name := fmt.Sprintf("%s-%s", fError.funcName, r.funcName)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					r.f(fError.f)
				}
				b.StopTimer()
			})
		}
	}
}
