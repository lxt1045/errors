package errors

import (
	"runtime"
	"strings"
	"sync"
	"unsafe"
)

const (
	baseSkip = 1

	// DefaultDepth 默认构建的调用栈深度
	DefaultDepth = 31
)

var (
	mStacks     = make(map[[DefaultDepth]uintptr][]string)
	mStacksLock sync.RWMutex
)

type stack struct {
	pcCache [DefaultDepth]uintptr
	npc     int
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

func (s *stack) Callers() (cs callers) {
	ok := false
	mStacksLock.RLock()
	cs, ok = mStacks[s.pcCache]
	mStacksLock.RUnlock()
	if ok {
		return
	}

	cs = parseStack(s.pcCache[:s.npc]) // 这步放在Lock()外虽然可能会造成重复计算,但是极大减少了锁争抢
	mStacksLock.Lock()
	mStacks[s.pcCache] = cs
	mStacksLock.Unlock()
	return
}

func parseStack(pcs []uintptr) (cs callers) {
	traces, more, f := runtime.CallersFrames(pcs), true, runtime.Frame{}
	for more {
		f, more = traces.Next()
		c := toCaller(f)
		if skipFile(c.File) && len(cs) > 0 {
			break
		}
		cs = append(cs, c.String())
		if strings.HasSuffix(f.Function, "main.main") && len(cs) > 0 {
			break
		}
	}
	return
}

func (s *stack) MarshalJSON() (bs []byte, err error) {
	cs := s.Callers()
	bs = make([]byte, 0, cs.jsonSize())
	bs = cs.json(bs)
	return
}

func (s *stack) String() string {
	cs := s.Callers()
	bs := make([]byte, 0, cs.textSize())
	bs = cs.text(bs)
	return *(*string)(unsafe.Pointer(&bs))
}
func (s *stack) fmt() (cs callers) {
	return s.Callers()
}

type callers []string

func (cs *callers) jsonSize() (l int) {
	if len(*cs) == 0 {
		return
	}
	l = len(*cs)*3 + 2 - 1
	for _, str := range *cs {
		l += len(str)
	}
	return
}

func (cs *callers) textSize() (l int) {
	if len(*cs) == 0 {
		return
	}
	l = len(*cs)*7 - 3
	for _, str := range *cs {
		l += len(str)
	}
	return
}

func (cs *callers) json(bs []byte) []byte {
	bs = append(bs, '[')
	for i, str := range *cs {
		if i != 0 {
			bs = append(bs, ',')
		}
		bs = append(bs, '"')
		bs = append(bs, str...)
		bs = append(bs, '"')
	}
	bs = append(bs, ']')
	return bs
}

func (cs *callers) text(bs []byte) []byte {
	for i, str := range *cs {
		if i != 0 {
			bs = append(bs, ", \n"...)
		}
		bs = append(bs, "    "...)
		bs = append(bs, str...)
	}
	return bs
}
