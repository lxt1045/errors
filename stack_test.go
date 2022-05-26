package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testPCs = [DefaultDepth]uintptr{
		0: 189989,
	}
	testFrame = frame{189989}
)

func TestMain(m *testing.M) {
	mStacks[testPCs] = []string{"(file1:88) func1"}

	mFrames[testFrame[0]] = "(file1:88) func1"
	m.Run()
}

func F1() *stack { //nolint
	return NewStack(0, 0)
}

func Test_NewStack(t *testing.T) {
	F := func() *stack {
		return F1()
	}

	s := F()

	_ = s.String()
	t.Log("\n", s)
	bs, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("MarshalIndent:", string(bs))

}

func Test_stack(t *testing.T) {
	t.Run("NewStack", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, s := runtime.Callers(baseSkip, pcs[:]), NewStack(0, 0)
		assert.True(t, equalStack(t, pcs[:npc], s.pcCache[:s.npc]))
	})

	t.Run("buildStack", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, s := runtime.Callers(baseSkip, pcs[:]), buildStack(0)
		assert.True(t, equalStack(t, pcs[:npc], s.pcCache[:s.npc]))
	})

	t.Run("stack.Callers", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, s := runtime.Callers(baseSkip, pcs[:]), buildStack(0)
		assert.Equal(t, npc, s.npc)

		callers := s.Callers()
		for i, caller := range callers {
			f, _ := runtime.CallersFrames([]uintptr{pcs[i]}).Next()
			c := toCaller(f)
			assert.Equal(t, caller, c.String())
		}
	})

	t.Run("stack.parse", func(t *testing.T) {
		s := buildStack(0)
		callers := s.Callers()
		cacheCallers := s.Callers()
		assert.Equal(t, callers, cacheCallers)
	})

	t.Run("stack.parseSlow", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, s := runtime.Callers(baseSkip, pcs[:]), buildStack(0)
		assert.Equal(t, npc, s.npc)

		for i, caller := range s.Callers() {
			f, _ := runtime.CallersFrames([]uintptr{pcs[i]}).Next()
			c := toCaller(f)
			assert.Equal(t, caller, c.String())
		}
	})

	t.Run("stack.json", func(t *testing.T) {
		s := NewStack(0, 0)
		bs1, err := s.MarshalJSON()
		assert.Nil(t, err)
		bs2, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.Equal(t, bs1, bs2)

		st := callers{}
		err = json.Unmarshal(bs1, &st)
		assert.Nil(t, err)
		assert.NotEqual(t, len(st), 0)
		assert.Equal(t, s.Callers(), st)
	})

	t.Run("stack.String", func(t *testing.T) {
		s := &stack{
			pcCache: testPCs,
		}
		str := `    (file1:88) func1`
		assert.Equal(t, s.String(), str)
	})
}

func BenchmarkNewStack2(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"lxt.buildStack", func() {
			buildStack(0)
		}},
		{"lxt.NewStack.0", func() {
			NewStack(0, 0)
		}},
		{"lxt.NewStack.16", func() {
			NewStack(0, 16)
		}},
		{"lxt.NewStack.8", func() {
			NewStack(0, 8)
		}},
		{"lxt.NewStack.4", func() {
			NewStack(0, 4)
		}},
	}
	for _, r := range runs {
		name := fmt.Sprintf("%s-1", r.funcName)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f()
			}
			b.StopTimer()
		})
		name = fmt.Sprintf("%s-5", r.funcName)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				func() {
					func() {
						func() {
							func() {
								func() {
									r.f()
								}()
							}()
						}()
					}()
				}()
			}
			b.StopTimer()
		})
	}
}

func BenchmarkNewStack(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"runtime.Callers", func() {
			pc := [DefaultDepth]uintptr{}
			runtime.Callers(baseSkip, pc[:])
		}},
		{"runtime.Callers.make", func() {
			pc := make([]uintptr, 32)
			runtime.Callers(baseSkip, pc[:])
		}},
		{"lxt.buildStack", func() {
			buildStack(1)
		}},
		{"lxt.NewErr", func() {
			_ = NewErr(0, "")
		}},
	}
	for _, r := range runs {
		name := fmt.Sprintf("%s-1", r.funcName)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f()
			}
			b.StopTimer()
		})
		name = fmt.Sprintf("%s-5", r.funcName)
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				func() {
					func() {
						func() {
							func() {
								func() {
									r.f()
								}()
							}()
						}()
					}()
				}()
			}
			b.StopTimer()
		})
	}
}

func Benchmark_Callers(b *testing.B) {
	var stackCache = []stack{
		buildStack(0), buildStack(1),
		buildStack(0), buildStack(1),
		buildStack(0), buildStack(1),
		buildStack(0), buildStack(1),
		buildStack(0), buildStack(1),
	}
	runs := []struct {
		funcName string         //函数名字
		f        func(s *stack) //调用方法
		stacks   []stack
	}{
		{"s.Callers", func(s *stack) {
			s.Callers()
		}, nil},
		{"parseStack", func(s *stack) {
			parseStack(s.pcCache[:s.npc])
		}, nil},
	}

	for _, r := range runs {
		nCaller := len(stackCache)
		name := fmt.Sprintf("%s-%d", r.funcName, nCaller)
		b.Run(name, func(b *testing.B) {
			r.stacks = make([]stack, b.N)
			for i := range r.stacks {
				r.stacks[i] = stackCache[i%nCaller]
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f(&r.stacks[i])
			}
			b.StopTimer()
		})
		b.Run(name+".warm_up", func(b *testing.B) {
			r.stacks = make([]stack, b.N)
			for i := range r.stacks {
				r.stacks[i] = stackCache[i%nCaller]
				r.f(&r.stacks[i])
			}
			b.ResetTimer()
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f(&r.stacks[i])
			}
			b.StopTimer()
		})
	}
}

func equalCaller(expected, actual uintptr) bool {
	fA, _ := runtime.CallersFrames([]uintptr{actual}).Next()
	cA := toCaller(fA)
	fE, _ := runtime.CallersFrames([]uintptr{expected}).Next()
	cE := toCaller(fE)
	return cA.String() == cE.String()
}

func equalStack(t *testing.T, expected, actual []uintptr) bool {
	_ = t
	if len(expected) != len(actual) {
		return false
	}
	if !equalCaller(expected[0], actual[0]) {
		return false
	}
	for i := 1; i < len(actual); i++ {
		if actual[i] != expected[i] {
			return false
		}
	}
	return true
}

func BenchmarkStackMarshal(b *testing.B) {
	s := buildStack(0)

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
