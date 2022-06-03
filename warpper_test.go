package errors

import (
	"encoding/json"
	"errors"
	stderrs "errors"
	"fmt"
	"runtime"
	"testing"

	pkgerrs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWarpNew(t *testing.T) {
	err := stderrs.New(errMsg)
	e := Wrap(err, errTrace).(*wrapper)
	t.Log("\n", e)
	t.Log(e.parse())
	bs, err := e.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
	bs, err = json.Marshal(e)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bs))
}

func Test_wrapper(t *testing.T) {
	err := stderrs.New(errMsg)
	t.Run("Wrap.nil", func(t *testing.T) {
		assert.Nil(t, Wrap(nil, errMsg))
	})
	t.Run("Wrap.stderr", func(t *testing.T) {
		pcs := [1]uintptr{}
		_, e := runtime.Callers(baseSkip, pcs[:]), Wrap(err, errTrace).(*wrapper)
		f, _ := runtime.CallersFrames(pcs[:]).Next()
		assert.Equal(t, e.parse().stack, toCaller(f).String())
	})
	t.Run("wrapper.parse.cache", func(t *testing.T) {
		e := Wrap(err, errTrace).(*wrapper)
		delete(mFrames, e.pc[0])
		caller := e.parse()
		cacheCaller := e.parse()
		assert.Equal(t, caller, cacheCaller)
	})
	t.Run("wrapper.parse", func(t *testing.T) {
		e := Wrap(err, errTrace).(*wrapper)
		e.pc = testFrame
		caller := e.parse().stack
		str := testFrameFunc
		assert.Equal(t, caller, str)
	})
	t.Run("wrapper.Unwrap", func(t *testing.T) {
		e := Wrap(err, errTrace).(*wrapper)
		err1 := e.Unwrap()
		assert.Equal(t, MarshalText(err), MarshalText(err1))
	})
	t.Run("wrapper.MarshalJSON", func(t *testing.T) {
		e := Wrap(err, errTrace).(*wrapper)
		e.pc = testFrame
		bs, err := json.Marshal(e)
		assert.Nil(t, err)
		str := `{"cause":"msg!","wrapper":[{"trace":"trace!","caller":"(file1:88) func1"}]}`
		assert.Equal(t, string(bs), str)
	})
	t.Run("wrapper.Error", func(t *testing.T) {
		e := Wrap(err, errTrace)
		e = Wrap(e, errTrace)
		w := e.(*wrapper)
		w.pc = testFrame
		bs := e.Error()
		assert.Equal(t, string(bs), string(w.Error()))
		str := "trace!,\n    (file1:88) func1;"
		assert.Equal(t, string(bs), str)
	})
	t.Run("wrapper.%+v", func(t *testing.T) {
		e := Wrap(err, errTrace)
		e.(*wrapper).pc = testFrame
		e = Wrap(e, errTrace)
		w := e.(*wrapper)
		w.pc = testFrame
		bs := fmt.Sprintf("%s", e)
		assert.Equal(t, string(bs), string(w.Error()))
		bs1 := fmt.Sprintf("%v", e)
		assert.Equal(t, string(bs), string(bs1))
		bs2 := fmt.Sprintf("%+v", e)
		str := "msg!;\ntrace!,\n    (file1:88) func1;\ntrace!,\n    (file1:88) func1;"
		assert.Equal(t, string(bs2), str)
	})

}
func BenchmarkWrap(b *testing.B) {
	runs := []struct {
		funcName string                //函数名字
		f        func(depth int) error //调用方法
	}{
		{"stdWrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = fmt.Errorf("%w", errors.New(errTrace))
			}
			return err
		}},
		{"Wrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = Wrap(err, errTrace)
			}
			return err
		}},
		{"pkg.Wrap", func(depth int) error {
			err := errors.New(errMsg)
			for i := 0; i < depth; i++ {
				err = pkgerrs.Wrap(err, errTrace)
			}
			return err
		}},
	}
	depths := []int{1, 10} //嵌套深度
	for _, r := range runs {
		for _, depth := range depths {
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

func BenchmarkWrapperMarshal(b *testing.B) {
	// TODO
	var err error = NewCause(0, errCode, errMsg)
	for i := 0; i < 0; i++ {
		err = Wrap(err, errTrace)
	}
	b.Run("json", func(b *testing.B) {
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
