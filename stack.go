package errors

import (
	"bytes"
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
	// 调用栈名字的缓存
	mStacks     = make(map[[DefaultDepth]uintptr][]string)
	mStacksLock sync.RWMutex
)

// 用于描述调用栈
type stack struct {
	npc     int                   // runtime 的调用栈
	pcCache [DefaultDepth]uintptr // runtime 的调用栈缓存,避免一次内存分配
}

// NewStack 锚定调用栈
func NewStack(skip int, depths ...int) (s *stack) { //nolint
	depth := DefaultDepth
	if len(depths) > 0 && depths[0] > 0 && depths[0] < DefaultDepth {
		depth = depths[0]
	}
	s = &stack{}
	s.npc = runtime.Callers(skip+1+baseSkip, s.pcCache[:depth])
	return
}

func buildStack(skip int) (s stack) {
	s.npc = runtime.Callers(skip+1+baseSkip, s.pcCache[:])
	return
}

func (s *stack) Callers() (callers []string) {
	ok := false
	mStacksLock.RLock()
	callers, ok = mStacks[s.pcCache]
	mStacksLock.RUnlock()
	if ok {
		return
	}

	callers = parse(s.pcCache[:s.npc]) // 这步放在Lock()外虽然可能会造成重复计算,但是极大减少了锁争抢
	mStacksLock.Lock()
	mStacks[s.pcCache] = callers
	mStacksLock.Unlock()
	return
}

func parse(pcs []uintptr) (callers []string) {
	traces, more, f := runtime.CallersFrames(pcs), true, runtime.Frame{}
	for more {
		f, more = traces.Next()
		c := toCaller(f)
		if skipFile(c.File) && len(callers) > 0 {
			break
		}
		callers = append(callers, c.String())
		if strings.HasSuffix(f.Function, "main.main") && len(callers) > 0 {
			break
		}
	}
	return
}

func (s *stack) json(buf *bytes.Buffer) {
	callers := s.Callers()
	buf.WriteString(`[`)
	for i, caller := range callers {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteByte('"')
		buf.WriteString(caller)
		buf.WriteByte('"')
	}
	buf.WriteString("]")
}

func (s *stack) text(buf *bytes.Buffer) {
	callers := s.Callers()
	for i, caller := range callers {
		if i != 0 {
			buf.WriteString(", \n")
		}
		buf.WriteString("    ")
		buf.WriteString(caller)
	}
}

func (s *stack) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	buf.WriteString(`{"stack":`)
	s.json(buf)
	buf.WriteString("}")
	return buf.Bytes(), nil
}

func (s *stack) String() string {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	s.text(buf)
	return bufToString(buf)
}
