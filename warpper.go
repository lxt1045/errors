package errors

import (
	"fmt"
	"runtime"
	"sync"
)

var (
	mFrames    = make(map[uintptr]frame)
	mFuncsLock sync.RWMutex
)

type wrapper struct {
	trace string
	err   error
	pc    [1]uintptr
}

func Wrap(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	e := &wrapper{
		trace: format,
		err:   err,
	}
	runtime.Callers(1+baseSkip, e.pc[:])
	return e
}

func (e *wrapper) Unwrap() error {
	return e.err
}

func (e *wrapper) Error() string {
	cache := e.fmt()
	buf := NewWriteBuffer(cache.textSize())
	cache.text(buf)
	return buf.String()
}

func (e *wrapper) MarshalJSON() ([]byte, error) {
	return MarshalJSON(e), nil
}

func (e *wrapper) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			s.Write(MarshalText(e))
			return
		}
		fallthrough
	case 's':
		s.Write([]byte(e.Error()))
		return
	case 'q':
		fmt.Fprintf(s, "%q", e.Error())
	}
}

func (e *wrapper) parse() *frame {
	mFuncsLock.RLock()
	c, ok := mFrames[e.pc[0]]
	mFuncsLock.RUnlock()
	if !ok {
		f, _ := runtime.CallersFrames(e.pc[:]).Next()
		c.stack = toCaller(f).String()
		l, yes := countEscape(c.stack)
		c.attr = uint64(l) << 32
		if yes {
			c.attr |= 1
		}

		mFuncsLock.Lock()
		mFrames[e.pc[0]] = c
		mFuncsLock.Unlock()
	}

	return &c
}

func (e *wrapper) fmt() fmtWrapper {
	return fmtWrapper{trace: e.trace, frame: e.parse()}
}

type frame struct {
	stack string
	attr  uint64 // count:escape ==> uint32:uint32
}
type fmtWrapper struct {
	trace       string
	traceEscape bool
	*frame
}

func (f *fmtWrapper) jsonSize() (l int) {
	l, f.traceEscape = countEscape(f.trace)
	l += len(`{"trace":"","caller":""}`) + (int(f.attr) >> 32)
	return
}

func (f *fmtWrapper) textSize() int {
	return len(",\n    ;") + len(f.trace) + len(f.stack)
}

func (f *fmtWrapper) json(buf *writeBuffer) {
	buf.WriteString(`{"trace":"`)
	if !f.traceEscape {
		buf.WriteString(f.trace)
	} else {
		buf.WriteEscape(f.trace)
	}
	buf.WriteString(`","caller":"`)
	if (f.attr & 1) == 0 {
		buf.WriteString(f.stack)
	} else {
		buf.WriteEscape(f.stack)
	}
	buf.WriteString(`"}`)
	return
}

func (f *fmtWrapper) text(buf *writeBuffer) {
	buf.WriteString(f.trace)
	buf.WriteString(",\n    ")
	buf.WriteString(f.stack)
	buf.WriteByte(';')
	return
}
