package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_frame(t *testing.T) {
	t.Run("NewFrame", func(t *testing.T) {
		rpcs := [1]uintptr{}
		_, s := runtime.Callers(baseSkip, rpcs[:]), NewFrame(0, "test")
		assert.NotEqual(t, rpcs[0], s.rpc[0])

	})

	t.Run("buildFrame", func(t *testing.T) {
		rpcs := [1]uintptr{}
		_, s := runtime.Callers(baseSkip, rpcs[:]), buildFrame(0, "test")
		assert.NotEqual(t, rpcs[0], s.rpc[0])
	})

	t.Run("frame.Callers", func(t *testing.T) {
		rpcs := [1]uintptr{}
		_, s := runtime.Callers(baseSkip, rpcs[:]), buildFrame(0, "test")

		caller := s.Caller()
		f, _ := runtime.CallersFrames(rpcs[:]).Next()
		c := toCaller(f)
		assert.Equal(t, caller, c.String())

	})

	t.Run("frame.parse", func(t *testing.T) {
		s := buildFrame(0, "")
		caller := s.Caller()
		cacheCaller := s.Caller()

		assert.Equal(t, caller, cacheCaller)

	})

	t.Run("frame.parseSlow", func(t *testing.T) {
		rpcs := [1]uintptr{}
		_, s := runtime.Callers(baseSkip, rpcs[:]), NewFrame(0, "test")

		s.parseSlow()

		f, _ := runtime.CallersFrames(rpcs[:]).Next()
		c := toCaller(f)
		assert.Equal(t, s.caller.File, c.File)
		assert.Equal(t, s.caller.FuncName, c.FuncName)
	})

	t.Run("frame.json", func(t *testing.T) {
		trace := "test"
		s := NewFrame(0, trace)
		bs1, err := s.MarshalJSON()
		assert.Nil(t, err)
		bs2, err := json.Marshal(s)
		assert.Nil(t, err)
		assert.Equal(t, bs1, bs2)

		st := struct {
			Trace  string `json:"trace"`
			Caller string `json:"caller"`
		}{}
		err = json.Unmarshal(bs1, &st)
		assert.Nil(t, err)
		assert.Equal(t, st.Trace, trace)
		assert.Equal(t, s.Caller(), st.Caller)

	})

	t.Run("frame.String", func(t *testing.T) {
		s := &frame{
			trace:  "test",
			caller: caller{"file1:88", "func1"},
		}
		str := "test,\n    (file1:88) func1"
		assert.Equal(t, s.String(), str)
	})
}

func Test_NewFrame(t *testing.T) {
	s := NewFrame(0, "test")
	fmt.Println(s.Caller())
	bs, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(bs))
}

func BenchmarkNewFrame(b *testing.B) {
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
		rpc := [32]uintptr{}
		runtime.Callers(skip+3, rpc[:])
		return
	}
	buildFrame := func(skip, depth int) {
		buildFrame(skip+1, "")
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
		{1, "buildFrame", buildFrame},
		{16, "buildFrame", buildFrame},
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

type M struct {
	sync.RWMutex
	m map[int]int
}

func NewM(n int) (m *M) {
	m = &M{
		m: make(map[int]int, n),
	}
	for i := 0; i < n; i++ {
		m.m[i] = i
	}
	return
}

func (m *M) Get(k int) int {
	m.RLock()
	defer m.RUnlock()
	return m.m[k]
}

type SyncM struct {
	sync.Map
}

func NewSyncM(n int) (m *SyncM) {
	m = &SyncM{}
	for i := 0; i < n; i++ {
		m.Store(i, i)
	}
	return
}

func (m *SyncM) Get(k int) int {
	v, _ := m.Load(k)
	n, _ := v.(int)
	return n
}

func BenchmarkMap(b *testing.B) {
	N := 6
	for x := 0; x < N; x++ {
		n := 10
		for i := 0; i < x; i++ {
			n = 10 * n
		}
		stdMap, syncMap := NewM(n), NewSyncM(n)

		b.Run(fmt.Sprintf("syncMap-%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				k := i % n
				syncMap.Get(k)
			}
			b.StopTimer()
		})
		b.Run(fmt.Sprintf("stdMap-%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				k := i % n
				stdMap.Get(k)
			}
			b.StopTimer()
		})

	}
}
