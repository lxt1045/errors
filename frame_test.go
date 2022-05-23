package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Frame(t *testing.T) {
	t.Run("frame.String", func(t *testing.T) {
		pcs := [1]uintptr{}
		_, s := runtime.Callers(baseSkip, pcs[:]), NewFrame(0)
		f, _ := runtime.CallersFrames(pcs[:]).Next()
		assert.Equal(t, s.String(), toCaller(f).String())
	})

	t.Run("frame.parse", func(t *testing.T) {
		s := NewFrame(0)
		delete(mFrames, s[0])
		caller := s.String()
		cacheCaller := s.String()

		assert.Equal(t, caller, cacheCaller)

	})
}

func Test_NewFrame(t *testing.T) {
	s := NewFrame(0)
	t.Log(s.String())
	bs, err := s.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
	bs, err = json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func BenchmarkNewFrame(b *testing.B) {
	type run struct {
		depth    int                   //嵌套深度
		funcName string                //函数名字
		f        func(skip, depth int) //调用方法
	}

	runtimeCaller1 := func(skip, depth int) {
		pc := [1]uintptr{}
		runtime.Callers(skip+3, pc[:])
	}
	runtimeCaller16 := func(skip, depth int) {
		pc := [32]uintptr{}
		runtime.Callers(skip+3, pc[:])
	}
	NewFrame := func(skip, depth int) {
		NewFrame(skip + 1)
	}
	buildStack := func(skip, depth int) {
		buildStack(skip + 1)
	}
	newDepthErr := func(skip, depth int) {
		_ = NewErr(0, "ssss")
	}
	runs := []run{
		{1, "runtimeCaller", runtimeCaller1},
		{16, "runtimeCaller", runtimeCaller16},
		{1, "NewFrame", NewFrame},
		{16, "NewFrame", NewFrame},
		{1, "buildStack", buildStack},
		{16, "buildStack", buildStack},
		{0, "newDepthErr", newDepthErr},
	}
	for _, r := range runs {
		name := fmt.Sprintf("%s-%d", r.funcName, r.depth)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f(1, r.depth)
			}
			b.StopTimer()
		})
	}
}

func Benchmark_frame_parse(b *testing.B) {
	var frameCache = []frame{
		NewFrame(0), NewFrame(1),
		NewFrame(0), NewFrame(1),
		NewFrame(0), NewFrame(1),
		NewFrame(0), NewFrame(1),
		NewFrame(0), NewFrame(1),
	}
	runs := []struct {
		funcName string        //函数名字
		f        func(s frame) //调用方法
		frames   []frame
	}{
		{"s.parse", func(s frame) {
			parseFrame(s[0])
		}, nil},
		{"parse_no_cache", func(s frame) {
			f, _ := runtime.CallersFrames(s[:]).Next()
			_ = toCaller(f)
		}, nil},
		{"FuncForPC", func(s frame) {
			f := runtime.FuncForPC(s[0])
			file, line := f.FileLine(s[0])
			_ = toCaller(runtime.Frame{
				File:     file,
				Line:     line,
				Function: f.Name(),
			})
		}, nil},
	}

	for _, r := range runs {
		nCaller := len(frameCache)
		name := fmt.Sprintf("%s-%d", r.funcName, nCaller)
		b.Run(name, func(b *testing.B) {
			b.StopTimer()
			r.frames = make([]frame, b.N)
			for i := range r.frames {
				r.frames[i] = frameCache[i%nCaller]
			}
			b.StartTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f(r.frames[i])
			}
			b.StopTimer()
		})
	}
}

func BenchmarkFrameMarshal(b *testing.B) {
	s := buildFrame(0)

	b.Run("String", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = s.String()
		}
		b.StopTimer()
	})
	b.Run("MarshalJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = s.MarshalJSON()
		}
		b.StopTimer()
	})
	b.Run("json.Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(s)
		}
		b.StopTimer()
	})
}

func TestRuntimeCaller(t *testing.T) {
	t.Run("block", func(t *testing.T) {
		s := buildStack(0)
		traces, more, f := runtime.CallersFrames(s.pcCache[:s.npc]), true, runtime.Frame{}
		for more {
			f, more = traces.Next()
			t.Logf("%+v", f)
		}
		for _, pc := range s.pcCache[:s.npc] {
			f := runtime.FuncForPC(pc)
			file, line := f.FileLine(pc)
			t.Logf("Entry:%+v, Name:%+v, File:%+v, Line:%+v", f.Entry(), f.Name(), file, line)
			file, line = f.FileLine(f.Entry())
			t.Logf("Entry:%+v, Name:%+v, File:%+v, Line:%+v", f.Entry(), f.Name(), file, line)
		}
	})

}
