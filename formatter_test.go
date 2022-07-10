package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	stderrs "errors"
	"fmt"
	"testing"
	"unicode/utf8"

	pkgerrs "github.com/pkg/errors"
)

func TestCountEscape(t *testing.T) {
	safeSet1 := [utf8.RuneSelf]bool{
		' ':      true,
		'!':      true,
		'"':      false,
		'#':      true,
		'$':      true,
		'%':      true,
		'&':      true,
		'\'':     true,
		'(':      true,
		')':      true,
		'*':      true,
		'+':      true,
		',':      true,
		'-':      true,
		'.':      true,
		'/':      true,
		'0':      true,
		'1':      true,
		'2':      true,
		'3':      true,
		'4':      true,
		'5':      true,
		'6':      true,
		'7':      true,
		'8':      true,
		'9':      true,
		':':      true,
		';':      true,
		'<':      true,
		'=':      true,
		'>':      true,
		'?':      true,
		'@':      true,
		'A':      true,
		'B':      true,
		'C':      true,
		'D':      true,
		'E':      true,
		'F':      true,
		'G':      true,
		'H':      true,
		'I':      true,
		'J':      true,
		'K':      true,
		'L':      true,
		'M':      true,
		'N':      true,
		'O':      true,
		'P':      true,
		'Q':      true,
		'R':      true,
		'S':      true,
		'T':      true,
		'U':      true,
		'V':      true,
		'W':      true,
		'X':      true,
		'Y':      true,
		'Z':      true,
		'[':      true,
		'\\':     false,
		']':      true,
		'^':      true,
		'_':      true,
		'`':      true,
		'a':      true,
		'b':      true,
		'c':      true,
		'd':      true,
		'e':      true,
		'f':      true,
		'g':      true,
		'h':      true,
		'i':      true,
		'j':      true,
		'k':      true,
		'l':      true,
		'm':      true,
		'n':      true,
		'o':      true,
		'p':      true,
		'q':      true,
		'r':      true,
		's':      true,
		't':      true,
		'u':      true,
		'v':      true,
		'w':      true,
		'x':      true,
		'y':      true,
		'z':      true,
		'{':      true,
		'|':      true,
		'}':      true,
		'~':      true,
		'\u007f': true,
	}

	for i := 0; i < 128; i++ {
		if safeSet1[i] != safeSet[i] {
			t.Fatal(i)
		}
	}

	bs := make([]byte, 0, 128)
	for i := 0; i < 128; i++ {
		bs = append(bs, byte(i))
	}
	check := func(str string) {
		buf := &writeBuffer{}
		l, escape := countEscape(str)
		if l != 1 && !escape {
			t.Fatalf("c:%x, l:%d, escape:%v", str[0], l, escape)
		}
		buf.WriteEscape(str)
		bs := buf.Bytes()
		if l != len(bs) {
			t.Fatalf("c:%x, l:%d, bs:%s", str[0], l, string(bs))
		}
	}
	for _, c := range bs {
		str := string(c)
		check(str)
	}

	check(string('\u2028'))
	check(string('\u2029'))
	check(string([]byte{0xff, 0x00}))
}

func TestMarshalJSON(t *testing.T) {
	for _, depth := range []int{1, 10} {
		name := fmt.Sprintf("%s-%d", "MarshalJSON.wrapper", depth)
		t.Run(name, func(t *testing.T) {
			var err error = NewCode(0, errCode, errMsg)
			for i := 0; i < depth; i++ {
				err = Wrap(err, fmt.Sprintf("%d", i))
				err = Wrap(err, fmt.Sprintf("%d", i*1000))
			}
			bs, err := json.Marshal(err)
			t.Log("wrapper:", string(bs))
			if err != nil {
				bs := MarshalJSON(err)
				t.Log(string(bs))
				t.Fatal(err)
			}
			m := map[string]interface{}{}
			err = json.Unmarshal(bs, &m)
			if err != nil {
				t.Log(string(bs))
				t.Fatal(err)
			}
		})

		name = fmt.Sprintf("%s-%d", "MarshalJSON.std", depth)
		t.Run(name, func(t *testing.T) {
			err := stderrs.New(errMsg)
			for i := 0; i < depth; i++ {
				err = fmt.Errorf("%d:%w", i, err)
				err = fmt.Errorf("%d:%w", i*1000, err)
			}
			bs := MarshalJSON(err)
			t.Log(string(bs))
			m := map[string]interface{}{}
			err = json.Unmarshal(bs, &m)
			if err != nil {
				t.Log(string(bs))
				t.Fatal(err)
			}
		})
		name = fmt.Sprintf("%s-%d", "MarshalJSON.pkg", depth)
		t.Run(name, func(t *testing.T) {
			err := pkgerrs.New(errMsg)
			for i := 0; i < depth; i++ {
				err = pkgerrs.Wrap(err, fmt.Sprintf("%d", i))
				err = pkgerrs.Wrap(err, fmt.Sprintf("%d", i))
			}
			bs := MarshalJSON(err)
			t.Log(string(bs))
			m := map[string]interface{}{}
			err = json.Unmarshal(bs, &m)
			if err != nil {
				t.Log(string(bs))
				t.Fatal(err)
			}
		})
	}
}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkMarshal)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
go test -benchmem -run=^$ -bench "^(BenchmarkMarshal)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
go tool pprof ./errors.test mem.prof
web
*/
func BenchmarkMarshal(b *testing.B) {
	var err error = NewCode(0, errCode, errMsg+"awesrdtfghjklsajghfdjkshdhgagdkaskdhakhkj")
	for i := 0; i < 0; i++ {
		err = Wrap(err, errTrace)
	}

	char := ' '
	_ = char
	str := `go test -benchmem -run=^$ -bench "^(BenchmarkMarshal)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c`
	b.Run("range", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, c := range str {
				char = c
			}
		}
		b.StopTimer()
	})
	b.Run("countEscape", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			countEscape(str)
		}
		b.StopTimer()
	})
	for i := 0; i < 5; i++ {
		b.Run("MarshalJSON", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MarshalJSON(err)
			}
			b.StopTimer()
		})
		b.Run("MarshalJSON2", func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				MarshalJSON2(err)
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
}
func Benchmark_JSON(b *testing.B) {
	for _, depth := range []int{1, 10} {
		name := fmt.Sprintf("%s-%d", "MarshalJSON", depth)
		b.Run(name, func(b *testing.B) {
			var err error = NewCode(0, errCode, errMsg)
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
			lxtErrs[i] = NewCode(0, errCode, errMsg)
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
			var err error = NewCode(0, errCode, errMsg)
			for j := 0; j < depth; j++ {
				err = Wrap(err, errTrace)
			}
			MarshalText(err)
		}},
		{lxt, "json", func(depth int) {
			var err error = NewCode(0, errCode, errMsg)
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
