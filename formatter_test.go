package errors

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	pkgerrs "github.com/pkg/errors"
)

/*
go test -benchmem -run=^$ -bench "^(BenchmarkMarshal)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
go test -benchmem -run=^$ -bench "^(BenchmarkMarshal)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
go tool pprof ./errors.test mem.prof
web
*/
func BenchmarkMarshal(b *testing.B) {
	var err error = NewErrSkip(0, errCode, errMsg)
	for i := 0; i < 0; i++ {
		err = Wrap(err, errTrace)
	}
	b.Run("MarshalJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalJSON(err)
		}
		b.StopTimer()
	})
	b.Run("MarshalText", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			MarshalText(err)
		}
		b.StopTimer()
	})
}
func Benchmark_JSON(b *testing.B) {
	for _, depth := range []int{1, 10} {
		name := fmt.Sprintf("%s-%d", "MarshalJSON", depth)
		b.Run(name, func(b *testing.B) {
			var err error = NewErrSkip(0, errCode, errMsg)
			for i := 0; i < depth; i++ {
				err = Wrap(err, errTrace)
			}
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MarshalJSON(err)
			}
			b.StopTimer()
		})
	}
}

func BenchmarkFormatting(b *testing.B) {
	depths := []int{1, 10} //嵌套深度
	std, pkg, lxt := "std", "pkg", "lxt"
	mErrCache := map[string]map[int][]error{std: {}, pkg: {}, lxt: {}}
	for _, depth := range depths {
		nErrs := 1024
		stdErrs, pkgErrs, lxtErrs := make([]error, nErrs), make([]error, nErrs), make([]error, nErrs)
		for i := range stdErrs {
			stdErrs[i] = errors.New(errMsg)
			pkgErrs[i] = pkgerrs.New(errMsg)
			lxtErrs[i] = NewErrSkip(0, errCode, errMsg)
			for j := 0; j < depth; j++ {
				stdErrs[i] = fmt.Errorf("%w; %s", stdErrs[i], errTrace)
				pkgErrs[i] = pkgerrs.Wrap(pkgErrs[i], errTrace)
				lxtErrs[i] = Wrap(lxtErrs[i], errTrace)
			}
		}
		mErrCache[std][depth] = stdErrs
		mErrCache[pkg][depth] = pkgErrs
		mErrCache[lxt][depth] = lxtErrs
	}

	stdText := func(err error) []byte {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		for ; err != nil; err = errors.Unwrap(err) {
			buf.WriteString(err.Error())
		}
		return buf.Bytes()
	}

	//log
	{
		depth := depths[0]
		b.Logf("std.text:\n%s\n", string(stdText(mErrCache[std][depth][0])))
		b.Logf("lxt.text:\n%s\n", string(MarshalText(mErrCache[lxt][depth][0])))
		b.Logf("lxt.jsob:\n%s\n", string(MarshalJSON(mErrCache[lxt][depth][0])))
		b.Logf("pkg.text:\n%+v\n", (mErrCache[pkg][depth][0]))
	}

	runs := []struct {
		t    string          //函数名字
		name string          //函数名字
		f    func(err error) //调用方法
	}{
		{std, "text", func(err error) {
			stdText(err)
		}},
		{lxt, "text", func(err error) {
			MarshalText(err)
		}},
		{lxt, "json", func(err error) {
			MarshalJSON(err)
		}},
		{pkg, "text.%+v", func(err error) {
			_ = fmt.Sprintf("%+v", err)
		}},
		{pkg, "text.%v", func(err error) {
			_ = fmt.Sprintf("%v", err)
		}},
	}

	for _, run := range runs {
		for _, depth := range depths {
			name := fmt.Sprintf("%s.%s-%d", run.t, run.name, depth)
			b.Run(name, func(b *testing.B) {
				errs := mErrCache[run.t][depth]
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					run.f(errs[i%len(errs)])
				}
				b.StopTimer()
			})
		}
	}
}

func BenchmarkNewAndFormatting(b *testing.B) {
	depths := []int{1, 10} //嵌套深度
	std, pkg, lxt := "std", "pkg", "lxt"

	stdText := func(err error) []byte {
		buf := bytes.NewBuffer(make([]byte, 0, 1024))
		for ; err != nil; err = errors.Unwrap(err) {
			buf.WriteString(err.Error())
		}
		return buf.Bytes()
	}

	runs := []struct {
		t    string          //函数名字
		name string          //函数名字
		f    func(depth int) //调用方法
	}{
		{std, "text", func(depth int) {
			err := errors.New(errMsg)
			for j := 0; j < depth; j++ {
				err = fmt.Errorf("%w; %s", err, errTrace)
			}
			stdText(err)
		}},
		{lxt, "text", func(depth int) {
			var err error = NewErrSkip(0, errCode, errMsg)
			for j := 0; j < depth; j++ {
				err = Wrap(err, errTrace)
			}
			MarshalText(err)
		}},
		{lxt, "json", func(depth int) {
			var err error = NewErrSkip(0, errCode, errMsg)
			for j := 0; j < depth; j++ {
				err = Wrap(err, errTrace)
			}
			MarshalJSON(err)
		}},
		{pkg, "text.%+v", func(depth int) {
			err := pkgerrs.New(errMsg)
			for j := 0; j < depth; j++ {
				err = pkgerrs.Wrap(err, errTrace)
			}
			_ = fmt.Sprintf("%+v", err)
		}},
		{pkg, "text.%v", func(depth int) {
			err := pkgerrs.New(errMsg)
			for j := 0; j < depth; j++ {
				err = pkgerrs.Wrap(err, errTrace)
			}
			_ = fmt.Sprintf("%v", err)
		}},
	}

	for _, run := range runs {
		for _, depth := range depths {
			name := fmt.Sprintf("%s.%s-%d", run.t, run.name, depth)
			b.Run(name, func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					run.f(depth)
				}
				b.StopTimer()
			})
		}
	}
}
