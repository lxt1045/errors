package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_stack(t *testing.T) {
	t.Run("NewStack", func(t *testing.T) {
		rpcs := [DefaultDepth]uintptr{}
		nrpc, s := runtime.Callers(baseSkip, rpcs[:]), NewStack(0, 0)
		assert.Equal(t, nrpc, s.nrpc)
		for i := 1; i < nrpc; i++ {
			assert.Equal(t, rpcs[i], s.rpcCache[i])
		}
	})

	t.Run("buildStack", func(t *testing.T) {
		rpcs := [DefaultDepth]uintptr{}
		nrpc, s := runtime.Callers(baseSkip, rpcs[:]), buildStack(0)
		assert.Equal(t, nrpc, s.nrpc)
		for i := 1; i < nrpc; i++ {
			assert.Equal(t, rpcs[i], s.rpcCache[i])
		}
	})

	t.Run("stack.Callers", func(t *testing.T) {
		rpcs := [DefaultDepth]uintptr{}
		nrpc, s := runtime.Callers(baseSkip, rpcs[:]), buildStack(0)
		assert.Equal(t, nrpc, s.nrpc)

		callers := s.Callers()
		for i, caller := range callers {
			f, _ := runtime.CallersFrames([]uintptr{rpcs[i]}).Next()
			c := toCaller(f)
			assert.Equal(t, caller, c.String())
		}
	})

	t.Run("stack.parse", func(t *testing.T) {
		s := buildStack(0)
		callers := s.Callers()
		cacheCallers := s.Callers()

		assert.Equal(t, len(callers), len(cacheCallers))
		for i, caller := range callers {
			assert.Equal(t, caller, cacheCallers[i])
		}
	})

	t.Run("stack.parseSlow", func(t *testing.T) {
		rpcs := [DefaultDepth]uintptr{}
		nrpc, s := runtime.Callers(baseSkip, rpcs[:]), buildStack(0)
		assert.Equal(t, nrpc, s.nrpc)

		s.parseSlow()
		for i, caller := range s.callers {
			f, _ := runtime.CallersFrames([]uintptr{rpcs[i]}).Next()
			c := toCaller(f)
			assert.Equal(t, caller.File, c.File)
			assert.Equal(t, caller.FuncName, c.FuncName)
		}
	})

	t.Run("stack.json", func(t *testing.T) {
		s := NewStack(0, 0)
		bs1, err := s.MarshalJSON()
		assert.Nil(t, err)
		bs2, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.Equal(t, bs1, bs2)

		st := struct {
			Stack []string `json:"stack"`
		}{}
		err = json.Unmarshal(bs1, &st)
		assert.Nil(t, err)
		assert.NotEqual(t, len(st.Stack), 0)
		assert.Equal(t, len(s.Callers()), len(st.Stack))
		for i := range st.Stack {
			assert.Equal(t, s.Callers()[i], st.Stack[i])
		}
	})

	t.Run("stack.String", func(t *testing.T) {
		s := &stack{
			callers: []caller{
				{"file1:88", "func1"},
			},
		}
		str := `    (file1:88) func1`
		assert.Equal(t, s.String(), str)
	})
}

func F1() {
	f := NewStack(0, 0)
	fmt.Println(f)

	bs, err := json.Marshal(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(bs))
}

func F2() {
	F1()
}

func F3() {
	F2()
}

func Test_NewTraces(t *testing.T) {
	F3()

	go F3()
	time.Sleep(1 * time.Second)
}

func Test_NewStack(t *testing.T) {
	s := NewStack(0, 0)
	fmt.Println(s.Callers())
	bs, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bs))
}

func Test_NewStack1(t *testing.T) {
	s := NewStack(0, 0)
	fmt.Println(s.Callers())
	bs, err := s.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bs))
}
