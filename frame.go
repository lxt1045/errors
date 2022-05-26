package errors

import (
	"runtime"
	"sync"
)

var (
	mFrames    = make(map[uintptr]string)
	mFuncsLock sync.RWMutex
)

type frame [1]uintptr

func NewFrame(skips ...int) (s frame) { //nolit
	skip := 1 + baseSkip
	if len(skips) > 0 {
		skip += skips[0]
	}
	runtime.Callers(skip, s[:])
	return
}
func buildFrame(skip int) (s frame) {
	runtime.Callers(skip+1+baseSkip, s[:])
	return
}

func parseFrame(pc uintptr) (c string) {
	mFuncsLock.RLock()
	ok := false
	c, ok = mFrames[pc]
	mFuncsLock.RUnlock()
	if ok {
		return
	}

	f, _ := runtime.CallersFrames([]uintptr{pc}).Next()
	c = toCaller(f).String()

	mFuncsLock.Lock()
	mFrames[pc] = c
	mFuncsLock.Unlock()
	return
}

func (s frame) MarshalJSON() (bs []byte, err error) {
	c := parseFrame(s[0])
	h := `{"frame":"`
	bs = make([]byte, 0, len(c)+2)
	bs = append(bs, h...)
	bs = append(bs, c...)
	bs = append(bs, '"')
	bs = append(bs, '}')

	return
}

func (s frame) String() string {
	return parseFrame(s[0])
}
