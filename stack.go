package errors

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	baseSkip = 1

	// DefaultDepth 默认构建的调用栈深度
	DefaultDepth = 32
)

var (
	rootDirs      = []string{"src", "/pkg/mod/"} // file 会从rootDir开始截断
	pathSeparator = string([]byte{os.PathSeparator})

	// skipPkgs里的pkg会被忽略
	skipPkgs = []string{
		"github.com/cloudwego/kitex",
	}

	// 调用栈名字的缓存
	mCallers     = make(map[[DefaultDepth]uintptr][]caller)
	mCallersLock sync.RWMutex
)

// 用于描述调用栈
type stack struct {
	nrpc     int                   // runtime 的调用栈
	rpcCache [DefaultDepth]uintptr // runtime 的调用栈缓存,避免一次内存分配
	callers  []caller              // 解析后的缓存
}

// NewStack 锚定调用栈
func NewStack(skip, depth int) *stack { //nolint
	s := buildStack(skip + 1)
	if depth > 0 && depth < s.nrpc {
		s.nrpc = depth
	}
	return &s
}

func buildStack(skip int) (s stack) {
	s.nrpc = runtime.Callers(skip+1+baseSkip, s.rpcCache[:]) //nolint
	return
}

func (s *stack) Callers() (callers []string) {
	s.parse()
	for _, c := range s.callers {
		callers = append(callers, c.String())
	}
	return
}

func (s *stack) parse() {
	if len(s.callers) > 0 {
		return
	}
	mCallersLock.RLock()
	ok := false
	s.callers, ok = mCallers[s.rpcCache]
	mCallersLock.RUnlock()
	if ok {
		return
	}

	s.parseSlow() // 这步放在Lock()外虽然可能会造成重复计算,但是极大减少了锁争抢
	mCallersLock.Lock()
	if _, ok := mCallers[s.rpcCache]; !ok {
		mCallers[s.rpcCache] = s.callers
	}
	mCallersLock.Unlock()
}

func (s *stack) parseSlow() {
	traces, more, f := runtime.CallersFrames(s.rpcCache[:s.nrpc]), true, runtime.Frame{}
	for more {
		f, more = traces.Next()
		if skipFunc(f.Function) && len(s.callers) > 0 {
			break
		}
		s.callers = append(s.callers, toCaller(f))
		if strings.HasSuffix(f.Function, "main.main") && len(s.callers) > 0 {
			break
		}
	}
	if len(s.callers) == 0 {
		s.callers = []caller{{File: "nil", FuncName: "nil"}}
	}
}

func (s *stack) json(buf *bytes.Buffer) {
	callers := s.Callers()
	buf.Write([]byte(`[`))
	for i, caller := range callers {
		if i != 0 {
			buf.Write([]byte(","))
		}
		buf.WriteByte('"')
		buf.WriteString(caller)
		buf.WriteByte('"')
	}
	buf.Write([]byte("]"))
}

func (s *stack) text(buf *bytes.Buffer) {
	callers := s.Callers()
	for i, caller := range callers {
		if i != 0 {
			buf.Write([]byte(", \n"))
		}
		buf.Write([]byte("    "))
		buf.WriteString(caller)
	}
}

func (s *stack) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.Write([]byte(`{"stack":`))
	s.json(buf)
	buf.Write([]byte("}"))
	return buf.Bytes(), nil
}

func (s *stack) String() string {
	buf := bytes.NewBuffer(nil)
	s.text(buf)
	return bufToString(buf)
}
