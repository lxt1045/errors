package errors

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const (
	baseSkip     = 1
	DefaultDepth = 31 // 默认构建的调用栈深度
)

var (
	mStacks     = make(map[[DefaultDepth]uintptr]*callers)
	mStacksLock sync.RWMutex
)

func NewCause(skip, code int, format string, a ...interface{}) (err *Cause) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	err = &Cause{code: code, msg: format}
	err.npc = runtime.Callers(skip+1+baseSkip, err.pcs[:])
	return
}

//CloneAs 利用 code 和 msg 生成一个包含当前stack的新Error,
func CloneAs(e error, skips ...int) *Cause {
	skip := 1 + baseSkip
	if len(skips) > 0 {
		skip += skips[0]
	}
	err := &Cause{}
	err.code, err.msg = GetCodeMsg(e)
	err.npc = runtime.Callers(1+baseSkip, err.pcs[:])
	return err
}

func NewErr(code int, format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	err := &Cause{code: code, msg: format}
	err.npc = runtime.Callers(1+baseSkip, err.pcs[:])
	return err
}

//New 替换 errors.New
func New(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	err := &Cause{code: DefaultCode, msg: format}
	/*
		TODO:
		可以通过asm的SP和BP拿列表,然后转成完整的pcs;
		因为可能内联,所以asm拿到的列表中有的pc可能实际会对应n个pc
		https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/traceback.go#L356
		https://github.com/golang/go/blob/dev.boringcrypto.go1.18/src/runtime/traceback.go#L194里的
		_func.funcFlag&funcFlag_TOPFRAME!=0  表示栈顶,不再继续执行,,,通过FuncForPC(pc)判断pc有效性?
	*/
	err.npc = runtime.Callers(1+baseSkip, err.pcs[:])
	return err
}

//Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	err := &Cause{code: DefaultCode, msg: format}
	err.npc = runtime.Callers(1+baseSkip, err.pcs[:])
	return err
}

type Cause struct {
	msg  string //业务错误信息
	code int    //业务错误码

	npc int
	pcs [DefaultDepth]uintptr
}

func (e *Cause) Code() int {
	return e.code
}

func (e *Cause) Message() string {
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
func (s *Cause) parse() (cs *callers) {
	ok := false
	mStacksLock.RLock()
	cs, ok = mStacks[s.pcs]
	mStacksLock.RUnlock()
	if ok {
		return
	}
	cs = &callers{}
	cs.stack = parseSlow(s.pcs[:s.npc]) // 这步放在Lock()外虽然可能会造成重复计算,但是极大减少了锁争抢
	l := 0
	for i, str := range cs.stack {
		lStack, yes := countEscape(str)
		l += lStack
		if yes {
			cs.attr |= 1 << i
		}
	}
	cs.attr |= uint64(l) << 32
	mStacksLock.Lock()
	mStacks[s.pcs] = cs
	mStacksLock.Unlock()
	return
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
	return fmtCause{code: strconv.Itoa(e.code), msg: e.msg, callers: e.parse()}
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
	// l, f.msgEscape = len(f.msg)*11/10, true
	l += len(f.code) + len(`{"code":,"msg":""}`)
	if len(f.stack) == 0 {
		return
	}
	l += len(`,"stack":[]`) + len(f.stack)*len(`,""`) - len(`,`) + (int(f.attr) >> 32)
	return
}

func (f *fmtCause) textSize() (l int) {
	l = 2 + len(f.code) + len(f.msg)
	if len(f.stack) == 0 {
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
	if len(f.stack) > 0 {
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
