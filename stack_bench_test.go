package errors

import (
	"fmt"
	"runtime"
	"testing"
)

func BenchmarkNewStack(b *testing.B) {
	type run struct {
		depth    int                   //嵌套深度
		funcName string                //函数名字
		f        func(skip, depth int) //调用方法
	}

	runtimeCaller1 := func(skip, depth int) {
		rpc := [1]uintptr{}
		runtime.Callers(skip+3, rpc[:])
		return
	}
	runtimeCaller16 := func(skip, depth int) {
		rpc := [16]uintptr{}
		runtime.Callers(skip+3, rpc[:])
		return
	}
	newDepthStack := func(skip, depth int) {
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
		{1, "newStack", newDepthStack},
		{16, "newStack", newDepthStack},
		{16, "newDepthErr", newDepthErr},
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

func BenchmarkStackCaller(b *testing.B) {
	// BenchmarkNewStack(b)
	type run struct {
		funcName string         //函数名字
		f        func(s *stack) //调用方法
		stacks   []stack
	}
	parse := func(s *stack) {
		s.parse()
	}
	parseSlow := func(s *stack) {
		s.parseSlow()
	}
	depths := []int{1, 16}
	nCallers := []int{10, 0}
	runs := []run{
		{"s.parse", parse, nil},
		{"s.parseSlow", parseSlow, nil},
	}
	for _, d := range depths {
		_, ok := mStacksCache[d]
		if !ok {
			b.Fatalf("mStacksCache[%d] is not exist", d)
		}
	}
	for _, depth := range depths {
		for _, nCaller := range nCallers {
			for _, r := range runs {
				name := fmt.Sprintf("%s-%d-%d", r.funcName, depth, nCaller)
				b.Run(name, func(b *testing.B) {
					b.StopTimer()
					//如果在这里 调用newStack 会花费大量时间,导致测试很慢
					stacksCache := mStacksCache[depth]
					if nCaller <= 0 || nCaller > len(stacksCache) {
						nCaller = len(stacksCache)
					}
					r.stacks = make([]stack, b.N)
					for i := range r.stacks {
						r.stacks[i] = stacksCache[i%nCaller]
						r.stacks[i].callers = nil
					}
					b.StartTimer()
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						r.f(&r.stacks[i])
					}
					b.StopTimer()
				})
			}
		}
	}
}
