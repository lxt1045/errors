package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"unsafe"

	stderrors "errors"

	pkgerrs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

const (
	errCode   = 88888
	errMsg    = "msg!"
	errTrace  = "trace!"
	errFormat = "format:%v"
)

var (
	testPCs = [DefaultDepth]uintptr{
		0: 189989,
	}
	testFrame     = [1]uintptr{189989}
	testFrameFunc = "(file1:88) func1"
)

func TestMain(m *testing.M) {
	// mStacks[testPCs] = &callers{stack: []string{testFrameFunc}, attr: uint64(len(testFrame) << 32)}

	// mFrames[testFrame[0]] = frame{stack: testFrameFunc, attr: uint64(len(testFrameFunc)) << 32}
	mFramesCache = func() unsafe.Pointer {
		m := map[uintptr]*frame{
			testFrame[0]: &frame{stack: testFrameFunc, attr: uint64(len(testFrameFunc)) << 32},
		}
		return unsafe.Pointer(&m)
	}()
	m.Run()
}

func TestNew(t *testing.T) {
	t.Run("NewCause", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(1, pcs[:]), NewCause(0, errCode, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Msg(), errMsg)
		assert.True(t, len(e.cache.stack) > 0)
		stack := parseSlow(pcs[:npc])
		assert.Equal(t, stack, e.cache.stack)
	})

	t.Run("NewCausef", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(1, pcs[:]), NewCause(0, errCode, errFormat, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Msg(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, len(e.cache.stack) > 0)
		stack := parseSlow(pcs[:npc])
		assert.Equal(t, stack, e.cache.stack)
	})

	t.Run("Wrap", func(t *testing.T) {
		pcs := [1]uintptr{}
		err := NewCause(0, errCode, errFormat, errMsg)
		_, e := runtime.Callers(1, pcs[:]), Wrap(err, errTrace)
		// assert.Equal(t, e.Code(), errCode)
		// assert.Equal(t, e.Msg(), fmt.Sprintf(errFormat, errMsg))
		_ = e
	})
}

func Test_Cause(t *testing.T) {
	t.Run("NewCause", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(1, pcs[:]), NewCause(0, errCode, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Msg(), errMsg)
		assert.True(t, len(e.cache.stack) > 0)
		stack := parseSlow(pcs[:npc])
		assert.Equal(t, stack, e.cache.stack)
	})

	t.Run("Is", func(t *testing.T) {
		err := NewCause(0, errCode, errMsg)
		err1 := NewCause(0, errCode, errMsg)
		assert.True(t, err.Is(err1))
		err2 := NewCause(0, errCode+1, errMsg)
		assert.False(t, err.Is(err2))
	})

	t.Run("json", func(t *testing.T) {
		c := NewCause(0, errCode, errMsg)

		bs1, err := c.MarshalJSON()
		assert.Nil(t, err)
		bs2, err := json.Marshal(c)
		assert.Nil(t, err)
		assert.Equal(t, bs1, bs2)

		st := struct {
			Code  int      `json:"code"`
			Msg   string   `json:"msg"`
			Stack []string `json:"stack"`
		}{}
		err = json.Unmarshal(bs1, &st)
		assert.Nil(t, err)
		assert.Equal(t, c.msg, st.Msg)
		assert.Equal(t, c.code, st.Code)
		assert.Equal(t, c.fmt().stack, st.Stack)
	})

	t.Run("Error", func(t *testing.T) {
		c := &Cause{
			code: errCode,
			msg:  errMsg,
			cache: &callers{
				stack: []string{"(file1:88) func1"},
			},
		}
		str := "88888, msg!;\n    (file1:88) func1;"
		assert.Equal(t, c.Error(), str)
		errStr := str
		assert.Equal(t, c.Error(), errStr)
		assert.Equal(t, fmt.Sprint(c), errStr)
	})
}

func Test_Text(t *testing.T) {
	ferr1 := func() error {
		err := NewErr(1600002, "message ferr1")
		return err
	}
	ferr2 := func() error {
		err := ferr1()
		err = Wrap(err, "log ferr2")
		return err
	}
	err := ferr2()
	fmt.Printf("c:%[1]c,\nv:%[1]v,\n+v:%+[1]v\nq:%[1]q\n", err)
}

//*/

func Test_JSON(t *testing.T) {
	ferr1 := func() error {
		err := NewErr(1600002, "message ferr1")
		return err
	}
	ferr2 := func() error {
		err := ferr1()
		err = Wrap(err, "log ferr2")
		return err
	}
	err := ferr2()
	err = Wrap(err, "wrap1 ...")
	err = Wrap(err, "wrap2 ...")

	fmt.Println(string(MarshalJSON(err)))

	bs, e := json.MarshalIndent(err, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(bs))
}

/*
go test -benchmem -run=^$ -bench "^(BenchmarkNewCause1)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./errors.test cpu.prof
go test -benchmem -run=^$ -bench "^(BenchmarkNewCause1)$" github.com/lxt1045/errors -test.memprofilerate=1 -count=1 -v -memprofile mem.prof -c
go tool pprof ./errors.test mem.prof
*/
func BenchmarkNewCause1(b *testing.B) {
	b.Run("NewCause", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewCause(0, 0, errMsg)
		}
	})
	b.Run("NewCause2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewCause2(0, 0, errMsg)
		}
	})
	b.Run("NewCause", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewCause(0, 0, errMsg)
		}
	})
	b.Run("NewCause2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NewCause2(0, 0, errMsg)
		}
	})
}

func Test_Xi(t *testing.T) {
	m := make(map[*[DefaultDepth]uintptr]*callers)
	a := [DefaultDepth]uintptr{
		1: 222,
		2: 333,
	}
	b := a
	m[&a] = &callers{
		attr: 88888,
	}
	fmt.Println(m[&a])
	fmt.Println(m[&b])
}

func BenchmarkEscape(b *testing.B) {
	b.Run("escape", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := [DefaultDepth]uintptr{}
			_ = buildStack(pcs[:])
		}
	})
	b.Run("escape-0", func(b *testing.B) {
		pcs := [DefaultDepth]uintptr{}
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = buildStack(pcs[:])
		}
	})
	b.Run("escape", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := pool.Get().(*[DefaultDepth]uintptr)
			_ = buildStack(pcs[:])
			pool.Put(pcs)
		}
	})
	b.Run("pool", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := pool.Get().(*[DefaultDepth]uintptr)
			pool.Put(pcs)
		}
	})

	b.Run("not-escape", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := [DefaultDepth]uintptr{}
			p := uintptr(unsafe.Pointer(&pcs))
			pp := (*[DefaultDepth]uintptr)(unsafe.Pointer(p))[:]
			_ = buildStack(pp)
			runtime.KeepAlive(&pcs)
		}
	})
}

func BenchmarkNewCause(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"runtime.Callers", func() {
			pc := [DefaultDepth]uintptr{}
			runtime.Callers(1, pc[:])
		}},
		{"lxt.NewCause", func() {
			NewCause(0, 0, errMsg)
		}},
		{"lxt.NewCause2", func() {
			NewCause2(0, 0, errMsg)
		}},
	}
	depths := []int{1, 10, 100} //嵌套深度
	for _, r := range runs {
		for _, depth := range depths {
			name := fmt.Sprintf("%s-%d", r.funcName, depth)
			b.Run(name, func(b *testing.B) {
				b.StopTimer()
				deepCall(depth, func() {
					b.StartTimer()
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						r.f()
					}
					b.StopTimer()
				})
			})
		}
	}
}

func BenchmarkCaseMarshal(b *testing.B) {
	err := NewCause(0, 0, errMsg)

	b.Run("Error", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = err.Error()
		}
		b.StopTimer()
	})
	b.Run("MarshalJSON", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = err.MarshalJSON()
		}
		b.StopTimer()
	})
	b.Run("json.Marshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(err)
		}
		b.StopTimer()
	})
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

//
func BenchmarkNew(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"std.New", func() {
			_ = stderrors.New("ye error")
		}},
		{"runtime.Caller", func() {
			runtime.Caller(2)
		}},
		{"runtime.Callers", func() {
			var pcs [DefaultDepth]uintptr
			runtime.Callers(3, pcs[:])
		}},
		{"pkg.New", func() {
			_ = pkgerrs.New("ye error")
		}},
		{"pkg.WithStack", func() {
			_ = pkgerrs.WithStack(stderrors.New("ye error"))
		}},
		{"lxt.New", func() {
			_ = New("ye error")
		}},
		{"lxt.NewErr", func() {
			_ = NewErr(-1, "ye error")
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
