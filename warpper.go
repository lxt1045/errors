// MIT License
//
// Copyright (c) 2021 Xiantu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
	pc  [1]uintptr
	err error
	msg string
}

func WrapSlow(err error, format string, ifaces ...interface{}) error {
	if err == nil {
		return nil
	}
	if len(ifaces) > 0 {
		format = fmt.Sprintf(format, ifaces...)
	}
	e := &wrapper{
		err: err,
		msg: format,
	}
	runtime.Callers(baseSkip, e.pc[:])
	return e
}

func NewLineSlow(format string, ifaces ...interface{}) error {
	if len(ifaces) > 0 {
		format = fmt.Sprintf(format, ifaces...)
	}
	e := &wrapper{
		err: nil,
		msg: format,
	}
	runtime.Callers(baseSkip, e.pc[:])
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

var cacheWrapper = RCUCache[[1]uintptr, *frame]{
	New: func(k [1]uintptr) (v *frame) {
		f := &frame{}
		cf, _ := runtime.CallersFrames(k[:]).Next()
		f.stack = toCaller(cf).String()
		l, yes := countEscape(f.stack)
		f.attr = uint64(l) << 32
		if yes {
			f.attr |= 1
		}
		return f
	},
}

func (e *wrapper) parse2() (f *frame) {
	return cacheWrapper.Get(e.pc)
}

func (e *wrapper) fmt() fmtWrapper {
	return fmtWrapper{trace: e.msg, frame: e.parse()}
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
}

func (f *fmtWrapper) text(buf *writeBuffer) {
	buf.WriteString(f.trace)
	buf.WriteString(",\n    ")
	buf.WriteString(f.stack)
	buf.WriteByte(';')
}

type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}
