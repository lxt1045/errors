package errors

import (
	"errors"
	"fmt"
	"testing"

	pkgerrs "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func Test_Assert(t *testing.T) {
	err := NewCause(0, errCode, errMsg)
	t.Run("OK", func(t *testing.T) {
		OK(true, nil)
	})

	t.Run("!OK", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, err.code, err1.code)
			assert.Equal(t, err.msg, err1.msg)
			return true
		})
		OK(false, err)
	})

	t.Run("!OK", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, DefaultCode, err1.code)
			return true
		})
		OK(false, nil)
	})

	t.Run("NilErr", func(t *testing.T) {
		NilErr(nil)
		var err error
		NilErr(err)
	})

	t.Run("!NilErr", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, err.code, err1.code)
			return true
		})
		NilErr(err)
	})

	t.Run("!NilErr", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, DefaultCode, err1.code)
			assert.Equal(t, errMsg, err1.msg)
			return true
		})
		NilErr(errors.New(errMsg))
	})

	t.Run("!NilErr.pkg", func(t *testing.T) {
		pkgErr := pkgerrs.New(errMsg)
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, DefaultCode, err1.code)
			assert.Equal(t, fmt.Sprintf("%+v", pkgErr), err1.msg)
			return true
		})
		NilErr(pkgErr)
	})

	t.Run("Nil", func(t *testing.T) {
		Nil(nil, err)
		var err1 error
		Nil(err1, err)
	})

	t.Run("!Nil", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, err.code, err1.code)
			assert.Equal(t, err.msg, err1.msg)
			return true
		})
		Nil(err, err)
	})

	t.Run("!Nil", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, DefaultCode, err1.code)
			assert.NotEqual(t, errMsg, err1.msg)
			return true
		})
		Nil(err, nil)
	})

	t.Run("!Nilf", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1 := e.(*Cause)
			assert.Equal(t, errCode, err1.code)
			assert.Equal(t, errMsg, err1.msg)
			return true
		})
		Nilf(err, errCode, errMsg)
	})

	t.Run("IsNil", func(t *testing.T) {
		assert.True(t, IsNil(nil))
		assert.False(t, IsNil(interface{}("")))
		assert.True(t, IsNil(chan int(nil)))

		nilErr := error((*Cause)(nil))
		assert.False(t, nil == nilErr)
		assert.True(t, IsNil(nilErr))
	})

	t.Run("TryByFunc", func(t *testing.T) {
		defer TryByFunc(func(e interface{}) (ok bool) {
			err1, ok := e.(*Cause)
			assert.Equal(t, err.code, err1.code)
			assert.Equal(t, err.msg, err1.msg)
			return
		})
		NilErr(err)
	})

	t.Run("TryByFunc.nil", func(t *testing.T) {
		defer TryByFunc(func(interface{}) (ok bool) {
			assert.Fail(t, "may not in here")
			return false
		})
		NilErr(nil)
	})

	t.Run("TryByFunc.panic", func(t *testing.T) {
		assert.Panics(t, func() {
			defer TryByFunc(func(e interface{}) (ok bool) {
				return false
			})
			NilErr(err)
		})
	})

	t.Run("TryErr", func(t *testing.T) {
		e1 := func() (e error) {
			defer TryErr(&e)
			NilErr(err)
			return
		}()
		err1 := e1.(*Cause)
		assert.Equal(t, err.code, err1.code)
		assert.Equal(t, err.msg, err1.msg)
	})

	t.Run("TryErr.nil", func(t *testing.T) {
		e1 := func() (e error) {
			defer TryErr(&e)
			NilErr(nil)
			return
		}()
		assert.Nil(t, e1)
	})

	t.Run("TryErr.panic", func(t *testing.T) {
		assert.Panics(t, func() {
			var e error
			defer TryErr(&e)
			panic("")
		})
	})

	t.Run("TryErr.panic", func(t *testing.T) {
		assert.Panics(t, func() {
			var e error
			defer TryErr(&e)
			x := 0
			_ = 11 / x
		})
	})
}

// go test -benchmem -run=^$ -bench "^(BenchmarkTry)$" github.com/lxt1045/errors -count=1 -v -cpuprofile cpu.prof -c
// go tool pprof ./errors.test cpu.prof
// web
func BenchmarkTry(b *testing.B) {
	errNotNil := NewErr(errCode, errMsg)

	type run struct {
		funcName string       //函数名字
		f        func() error //调用方法
	}

	runs := []run{
		{"ifErr.nil", func() (err error) {
			if err == nil {
				return err
			}
			return nil
		}},
		{"ifErr.not-nil", func() (err error) {
			err = errNotNil
			if err == nil {
				return err
			}
			return nil
		}},
		{"defer", func() (err error) {
			defer func() {}()
			return
		}},
		{"perr", func() (err error) {
			defer func(perr *error) {
				*perr = nil
			}(&err)
			return
		}},
		{"TryByFunc.nil", func() (err error) {
			defer TryByFunc(func(e interface{}) (ok bool) {
				err, ok = e.(*Cause)
				return
			})
			NilErr(err)
			return
		}},
		{"TryByFunc.not-nil", func() (err error) {
			defer TryByFunc(func(e interface{}) (ok bool) {
				err, ok = e.(*Cause)
				return
			})
			err = errNotNil
			NilErr(err)
			return
		}},
		{"TryErr.nil", func() (err error) {
			defer TryErr(&err)
			NilErr(err)
			return
		}},
		{"TryErr.noy-nil", func() (err error) {
			defer TryErr(&err)
			err = errNotNil
			NilErr(err)
			return
		}},
	}

	for _, r := range runs {
		b.Run(r.funcName, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err := r.f()
				_ = err
				// fmt.Sprint(err)
			}
			b.StopTimer()
		})
	}
}
