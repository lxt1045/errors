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
	"strconv"
	"strings"
	"sync"

	"github.com/rs/zerolog"
)

const (
	baseSkip     = 2
	DefaultDepth = 32 // 默认构建的调用栈深度
)

var (
	cacheStack   = StackCache[*callers]{}
	cacheCallers = StackCache[[]caller]{}
	cacheCaller  = RCUCache[uintptr, *caller]{}

	pool = sync.Pool{
		New: func() any { return &[DefaultDepth]uintptr{} },
	}
)

func NewCodeSlow(skip, code int, format string, a ...interface{}) (c *Code) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Code{code: code, msg: format}

	pcs := pool.Get().(*[DefaultDepth]uintptr)
	n := runtime.Callers(skip+baseSkip, pcs[:DefaultDepth-skip])
	// key := toString(pcs[:n])

	cs := cacheStack.Get(pcs, n)
	if cs == nil {
		cs = &callers{}
		for _, c := range parseSlow(pcs[:n]) {
			cs.stack = append(cs.stack, c.String())
		}
		l := 0
		for i, str := range cs.stack {
			// 检查是否需要转换 JSON 特殊字符串
			lStack, yes := countEscape(str)
			l += lStack
			if yes {
				cs.attr |= 1 << i
			}
		}
		cs.attr |= uint64(l) << 32

		// 加入
		cacheStack.Set(pcs, n, cs)
	}
	pool.Put(pcs)
	c.cache = cs
	return
}

// Clone 利用 code 和 msg 生成一个包含当前stack的新Error,
func Clone(err error, skips ...int) error {
	skip := 0
	if len(skips) > 0 {
		skip += skips[0]
	}
	if c, ok := err.(*Code); ok {
		return NewCode(skip+1, c.code, c.msg)
	}
	return err
}

func NewErr(code int, format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCode(1, code, format)
}

func NewCodeNoStack(code int, msg string) error {
	return &Code{
		code: code,
		msg:  msg,
	}
}

// New 替换 errors.New
func New(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCode(1, DefaultCode, format)
}

// Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCode(1, DefaultCode, format)
}

type Code struct {
	msg  string //业务错误信息
	code int    //业务错误码

	cache *callers
	skip  int
}

func (e *Code) WithErr(err error) *Code {
	if err == nil {
		return NewCode(1, e.code, e.msg)
	}
	if f := getErrorFunc(err); f != nil {
		return NewCode(1, e.code, e.msg+"; "+f(err))
	}
	return NewCode(1, e.code, e.msg+"; "+err.Error())
}

func (e *Code) Clone(msg ...string) *Code {
	if len(msg) > 0 {
		return NewCode(1, e.code, e.msg+"; "+strings.Join(msg, ";"))
	}
	return NewCode(1, e.code, e.msg)
}

func (e *Code) Clonef(format string, a ...interface{}) *Code {
	msg := fmt.Sprintf(e.msg+"; "+format, a...)
	return NewCode(1, e.code, msg)
}

func (e *Code) New(msg ...string) *Code {
	if len(msg) > 0 {
		return NewCode(1, e.code, strings.Join(msg, ";"))
	}
	return NewCode(1, e.code, e.msg)
}

func (e *Code) Newf(format string, a ...interface{}) *Code {
	msg := fmt.Sprintf(format, a...)
	return NewCode(1, e.code, msg)
}

func (e *Code) SkipClone(skip int, msg ...string) *Code {
	if len(msg) > 0 {
		return NewCode(skip+1, e.code, e.msg+"; "+strings.Join(msg, ";"))
	}
	return NewCode(skip+1, e.code, e.msg)
}

func (e *Code) SkipClonef(skip int, format string, a ...interface{}) *Code {
	msg := fmt.Sprintf(e.msg+"; "+format, a...)
	return NewCode(skip+1, e.code, msg)
}
func (e *Code) SkipNew(skip int, msg ...string) *Code {
	if len(msg) > 0 {
		return NewCode(skip+1, e.code, strings.Join(msg, ";"))
	}
	return NewCode(skip+1, e.code, e.msg)
}

func (e *Code) SkipNewf(skip int, format string, a ...interface{}) *Code {
	msg := fmt.Sprintf(format, a...)
	return NewCode(skip+1, e.code, msg)
}

func (e *Code) Code() int {
	return e.code
}

func (e *Code) Msg() string {
	return e.msg
}

func (e *Code) Is(err error) bool {
	to, ok := err.(*Code)
	return ok && e.code != -1 && e.code == to.code
}

// Error error interface, 序列化为string, 包含调用栈
func (e *Code) Error() string {
	cache := e.fmt()
	buf := NewWriteBuffer(cache.textSize())
	cache.text(buf)
	return buf.String()
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *Code) MarshalJSON() (bs []byte, err error) {
	cache := e.fmt()
	buf := NewWriteBuffer(cache.jsonSize())
	cache.json(buf)
	return buf.Bytes(), nil
}
func parseSlow(pcs []uintptr) (cs []caller) {
	traces, more, f := runtime.CallersFrames(pcs), true, runtime.Frame{}
	for more {
		f, more = traces.Next()
		c := toCaller(f)
		if skipFile(c.FileLine) && len(cs) > 0 {
			break
		}
		cs = append(cs, c)
		if strings.HasSuffix(f.Function, "main.main") && len(cs) > 0 {
			break
		}
	}
	return
}
func (e *Code) fmt() (cs fmtCode) {
	return fmtCode{code: strconv.Itoa(e.code), msg: e.msg, callers: e.cache, skip: e.skip}
}

type callers struct {
	stack []string
	attr  uint64 // count:escape ==> uint32:uint32
}
type fmtCode struct {
	code      string
	msg       string
	skip      int
	msgEscape bool
	*callers
}

func (f *fmtCode) jsonSize() (l int) {
	l, f.msgEscape = countEscape(f.msg)
	l += len(f.code) + len(`{"code":,"msg":""}`)
	if len(f.stack) <= f.skip {
		return
	}
	l += len(`,"stack":[]`) + (len(f.stack)-f.skip)*len(`,""`) - len(`,`) + (int(f.attr) >> 32)
	return
}

func (f *fmtCode) textSize() (l int) {
	l = len(", ") + len(f.code) + len(f.msg)
	if f.callers == nil || len(f.stack) <= f.skip {
		return
	}
	l += (len(f.stack) - f.skip) * len(", \n    ")
	for _, str := range f.stack[f.skip:] {
		l += len(str) + 3
	}
	return
}

func (f *fmtCode) json(buf *writeBuffer) {
	buf.WriteString(`{"code":`)
	buf.WriteString(f.code)
	buf.WriteString(`,"msg":"`)
	if !f.msgEscape {
		buf.WriteString(f.msg)
	} else {
		buf.WriteEscape(f.msg)
	}
	buf.WriteByte('"')
	if len(f.stack) > f.skip {
		buf.WriteString(`,"stack":[`)
		for i, str := range f.stack[f.skip:] {
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
}

func (f *fmtCode) text(buf *writeBuffer) {
	buf.WriteString(f.code)
	buf.WriteString(", ")
	buf.WriteString(f.msg)
	if f.callers != nil && len(f.stack) > f.skip {
		buf.WriteString(";\n")
		for i, str := range f.stack[f.skip:] {
			if i != 0 {
				buf.WriteString(", \n")
			}
			buf.WriteString("    ")
			buf.WriteString(str)
		}
		buf.WriteByte(';')
	}
}

// MarshalZerologObject for zerolog
func (e *Code) MarshalZerologObject(evt *zerolog.Event) {
	evt.Int("code", e.code)
	evt.Str("msg", e.msg)
	evt.Array("stack", e)
	return
}

func (e *Code) MarshalZerologArray(a *zerolog.Array) {
	if e.cache != nil && len(e.cache.stack) > e.skip {
		for _, str := range e.cache.stack[e.skip:] {
			a.Str(str)
		}
	}
}
