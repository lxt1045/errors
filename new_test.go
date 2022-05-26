package errors

import (
	"fmt"
	"runtime"
	"testing"

	pkgerrs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_new(t *testing.T) {
	t.Run("New", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), New(errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), DefaultCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("Errorf", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), Errorf(errFormat, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), DefaultCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("Errf", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), Errf(errCode, errFormat, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("NewErr", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, err := runtime.Callers(baseSkip, pcs[:]), NewErr(errCode, errMsg)
		e := err.(*Cause)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("NewErrSkip", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(baseSkip, pcs[:]), NewErrSkip(0, errCode, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), errMsg)
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("NewErrfSkip", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		npc, e := runtime.Callers(baseSkip, pcs[:]), NewErrfSkip(0, errCode, errFormat, errMsg)
		assert.Equal(t, e.Code(), errCode)
		assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		assert.True(t, equalStack(t, pcs[:npc], e.stack.pcCache[:e.stack.npc]))
	})

	t.Run("Wrap", func(t *testing.T) {
		pcs := [1]uintptr{}
		err := NewErrfSkip(0, errCode, errFormat, errMsg)
		_, e := runtime.Callers(baseSkip, pcs[:]), Wrap(err, errTrace)
		// assert.Equal(t, e.Code(), errCode)
		// assert.Equal(t, e.Message(), fmt.Sprintf(errFormat, errMsg))
		_ = e
	})
}

//
func Benchmark_new(b *testing.B) {
	runs := []struct {
		funcName string //函数名字
		f        func() //调用方法
	}{
		{"runtime.Caller", func() {
			runtime.Caller(2)
		}},
		{"NewFrame", func() {
			NewFrame(0)
		}},
		{"runtime.Callers", func() {
			var pcs [DefaultDepth]uintptr
			runtime.Callers(3, pcs[:])
		}},
		{"buildStack", func() {
			buildStack(1 + 1)
		}},
		{"New", func() {
			_ = New("ye error")
		}},
		{"NewErr", func() {
			_ = NewErr(-1, "ye error")
		}},
		{"pkgNew", func() {
			_ = pkgerrs.New("ye error")
		}},
	}
	for _, r := range runs {
		name := r.funcName
		b.Run(name, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				r.f()
			}
			b.StopTimer()
		})
	}
}
