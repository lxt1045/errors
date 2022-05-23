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
	errTrace1 = "trace!1"
	errFormat = "format:%v"
)

func Test_Err(t *testing.T) {
	t.Run("buildCause", func(t *testing.T) {
		pcs := [DefaultDepth]uintptr{}
		_, s := runtime.Callers(baseSkip, pcs[:]), buildCause(errCode, errMsg, buildStack(0))
		assert.Equal(t, s.GetCode(), errCode)
		assert.Equal(t, s.GetMsg(), errMsg)
		assert.True(t, equalCaller(pcs[0], s.stack.pcCache[0]))
		assert.Equal(t, pcs[1:], s.stack.pcCache[1:])
	})

	t.Run("Is", func(t *testing.T) {
		err := NewErrSkip(0, errCode, errMsg)
		err1 := NewErrSkip(0, errCode, errTrace)
		assert.True(t, err.Is(err1))
		err2 := NewErrSkip(0, errCode+1, errMsg)
		assert.False(t, err.Is(err2))
	})

	t.Run("json", func(t *testing.T) {
		s := NewErrSkip(0, errCode, errMsg)

		bs1, err := s.MarshalJSON()
		assert.Nil(t, err)
		bs2, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.Equal(t, bs1, bs2)

		st := struct {
			Code  int      `json:"code"`
			Msg   string   `json:"msg"`
			Stack []string `json:"stack"`
		}{}
		err = json.Unmarshal(bs1, &st)
		assert.Nil(t, err)
		assert.Equal(t, s.Msg, st.Msg)
		assert.Equal(t, s.stack.Callers(), st.Stack)
	})

	t.Run("Error", func(t *testing.T) {
		s := &Cause{
			Code: errCode,
			Msg:  errMsg,
			stack: stack{
				npc:     1,
				pcCache: testPCs,
			},
		}
		str := "88888, msg!;\n    (file1:88) func1;"
		assert.Equal(t, s.Error(), str)
		errStr := str
		assert.Equal(t, s.Error(), errStr)
		assert.Equal(t, fmt.Sprint(s), errStr)
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
	fmt.Printf("s:%[1]s,\nv:%[1]v,\n+v:%+[1]v\nq:%[1]q\n", err)
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
