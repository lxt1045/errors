package errors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

const (
	baseSkip     = 2
	DefaultDepth = 32 // 默认构建的调用栈深度
)

type stackCacheType = map[[DefaultDepth]uintptr]*callers

var (
	mStackCache unsafe.Pointer = func() unsafe.Pointer {
		m := make(stackCacheType)
		return unsafe.Pointer(&m)
	}()

	pool = sync.Pool{
		New: func() any { return &[DefaultDepth]uintptr{} },
	}
)

//CloneAs 利用 code 和 msg 生成一个包含当前stack的新Error,
func CloneAs(e error, skips ...int) *Cause {
	skip := 1
	if len(skips) > 0 {
		skip += skips[0]
	}
	code, msg := GetCodeMsg(e)
	return NewCause(skip+1, code, msg)
}

func NewErr(code int, format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCause(1, code, format)
}

//New 替换 errors.New
func New(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCause(1, DefaultCode, format)
}

//Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCause(1, DefaultCode, format)
}

type Cause struct {
	msg  string //业务错误信息
	code int    //业务错误码

	cache *callers
}

func (e *Cause) Code() int {
	return e.code
}

func (e *Cause) Msg() string {
	return e.msg
}

func (e *Cause) Is(err error) bool {
	to, ok := err.(*Cause)
	return ok && e.code != -1 && e.code == to.code
}

//Error error interface, 序列化为string, 包含调用栈
func (e *Cause) Error() string {
	cache := e.fmt()
	buf := NewWriteBuffer(cache.textSize())
	cache.text(buf)
	return buf.String()
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *Cause) MarshalJSON() (bs []byte, err error) {
	cache := e.fmt()
	buf := NewWriteBuffer(cache.jsonSize())
	cache.json(buf)
	return buf.Bytes(), nil
}
func parseSlow(pcs []uintptr) (cs []string) {
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
func (e *Cause) fmt() (cs fmtCause) {
	return fmtCause{code: strconv.Itoa(e.code), msg: e.msg, callers: e.cache}
}

type callers struct {
	stack []string
	attr  uint64 // count:escape ==> uint32:uint32
}
type fmtCause struct {
	code      string
	msg       string
	msgEscape bool
	*callers
}

func (f *fmtCause) jsonSize() (l int) {
	l, f.msgEscape = countEscape(f.msg)
	l += len(f.code) + len(`{"code":,"msg":""}`)
	if len(f.stack) == 0 {
		return
	}
	l += len(`,"stack":[]`) + len(f.stack)*len(`,""`) - len(`,`) + (int(f.attr) >> 32)
	return
}

func (f *fmtCause) textSize() (l int) {
	l = 2 + len(f.code) + len(f.msg)
	if f.callers == nil || len(f.stack) == 0 {
		return
	}
	l += len(f.stack)*7 - 3
	for _, str := range f.stack {
		l += len(str) + 3
	}
	return
}

func (f *fmtCause) json(buf *writeBuffer) {
	buf.WriteString(`{"code":`)
	buf.WriteString(f.code)
	buf.WriteString(`,"msg":"`)
	if !f.msgEscape {
		buf.WriteString(f.msg)
	} else {
		buf.WriteEscape(f.msg)
	}
	buf.WriteByte('"')
	if len(f.stack) > 0 {
		buf.WriteString(`,"stack":[`)
		for i, str := range f.stack {
			if i != 0 {
				buf.WriteByte(',')
			}
			buf.WriteByte('"')
			if f.attr&(1<<i) == 0 {
				buf.WriteString(str)
			} else {
				buf.WriteEscape(str)
			}
			buf.WriteByte('"')
		}
		buf.WriteByte(']')
	}
	buf.WriteByte('}')
	return
}

func (f *fmtCause) text(buf *writeBuffer) {
	buf.WriteString(f.code)
	buf.WriteString(", ")
	buf.WriteString(f.msg)
	if f.callers != nil && len(f.stack) > 0 {
		buf.WriteString(";\n")
		for i, str := range f.stack {
			if i != 0 {
				buf.WriteString(", \n")
			}
			buf.WriteString("    ")
			buf.WriteString(str)
		}
		buf.WriteByte(';')
	}
	return
}
