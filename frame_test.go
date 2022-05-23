package errors

import (
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
	fmt.Println(s.String())
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
		return
	}
	runtimeCaller16 := func(skip, depth int) {
		pc := [32]uintptr{}
		runtime.Callers(skip+3, pc[:])
		return
	}
	NewFrame := func(skip, depth int) {
		NewFrame(skip + 1)
		return
	}
	buildStack := func(skip, depth int) {
		buildStack(skip + 1)
		return
	}
	newDepthErr := func(skip, depth int) {
		NewErr(0, "ssss") //nolint
		return
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
