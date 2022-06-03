package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"

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
	mStacks[testPCs] = &callers{stack: []string{testFrameFunc}, attr: uint64(len(testFrame) << 32)}

	mFrames[testFrame[0]] = frame{stack: testFrameFunc, attr: uint64(len(testFrameFunc)) << 32}
	m.Run()
}

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), New(errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), DefaultCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("Errorf", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), Errorf(errFormat, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), DefaultCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("Errf", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), NewErr(errCode, errFormat, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("NewErr", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), NewErr(errCode, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("NewCause", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(baseSkip, pcs[:]), NewCause(0, errCode, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("NewCausef", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(baseSkip, pcs[:]), NewCause(0, errCode, errFormat, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.pcs[:e.npc]))
	})

	t.Run("Wrap", func(t *testing.T) {
		pcs := [1]uintptr{}
		err := NewCause(0, errCode, errFormat, errMsg)
		_, e := runtime.Callers(baseSkip, pcs[:]), Wrap(err, errTrace)
		// assert.Equal(t, e.Code(), errCode)
		// assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		_ = e
	})
}
func Test_Cause(t *testing.T) {
	t.Run("NewCause", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		_, c := runtime.Callers(baseSkip, pcs[:]), NewCause(0, errCode, errMsg)
		assert.Equal(t, c.Code(), errCode)
		assert.Equal(t, c.Message(), errMsg)
		assert.True(t, equalCaller(pcs[0], c.pcs[0]))
		assert.Equal(t, pcs[1:], c.pcs[1:])
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
			npc:  1,
			pcs:  testPCs,
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

func BenchmarkNewCause(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"runtime.Callers", func() {
			pc := [DefaultDepth]uintptr{}
			runtime.Callers(baseSkip, pc[:])
		}},
		{"lxt.NewCause", func() {
			NewCause(0, 0, errMsg)
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
