package errors

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"unsafe"
)

var (
	mFramesCache unsafe.Pointer = func() unsafe.Pointer {
		m := make(map[uintptr]*frame)
		return unsafe.Pointer(&m)
	}()
)

type wrapper struct {
	pc     [1]uintptr
	err    error
	format string
	ifaces []interface{}
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

func (e *wrapper) parse() (f *frame) {
	mFC := *(*map[uintptr]*frame)(atomic.LoadPointer(&mFramesCache))
	f, ok := mFC[e.pc[0]]
	if !ok {
		f = &frame{}
		// file, n := runtime.FuncForPC(e.pc).FileLine(e.pc)
		cf, _ := runtime.CallersFrames(e.pc[:]).Next()
		f.stack = toCaller(cf).String()
		l, yes := countEscape(f.stack)
		f.attr = uint64(l) << 32
		if yes {
			f.attr |= 1
		}

		mFC2 := make(map[uintptr]*frame, len(mFC)+10)
		mFC2[e.pc[0]] = f
		for {
			p := atomic.LoadPointer(&mFramesCache)
			mFC = *(*map[uintptr]*frame)(p)
			for k, v := range mFC {
				mFC2[k] = v
			}
			swapped := atomic.CompareAndSwapPointer(&mFramesCache, p, unsafe.Pointer(&mFC2))
			if swapped {
				break
			}
		}
	}
	return f
}

func (e *wrapper) fmt() fmtWrapper {
	if len(e.ifaces) > 0 {
		trace := fmt.Sprintf(e.format, e.ifaces...)
		return fmtWrapper{trace: trace, frame: e.parse()}
	}
	return fmtWrapper{trace: e.format, frame: e.parse()}
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
