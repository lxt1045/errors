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
	"unsafe"
)

const (
	baseSkip     = 2
	DefaultDepth = 32 // 默认构建的调用栈深度
)

var (
	// cacheStack  = AtomicCache[string, *callers]{}
	cacheStack = Cache[string, *callers]{}

	pool = sync.Pool{
		New: func() any { return &[DefaultDepth]uintptr{} },
	}
)

func toString(p []uintptr) string {
	bs := (*[DefaultDepth * 8]byte)(unsafe.Pointer(&p[0]))[:len(p)*8]
	return *(*string)(unsafe.Pointer(&bs))
}
func fromString(str string) (p []uintptr) {
	bs := *(*[]byte)(unsafe.Pointer(&str))
	p = (*[DefaultDepth]uintptr)(unsafe.Pointer(&bs[0]))[:len(str)/8]
	return
}

func NewCodeSlow(skip, code int, format string, a ...interface{}) (c *Code) {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	c = &Code{code: code, msg: format}

	pcs := pool.Get().(*[DefaultDepth]uintptr)
	n := runtime.Callers(skip+baseSkip, pcs[:DefaultDepth-skip])
	key := toString(pcs[:n])

	cs := cacheStack.Get(key)
	if cs == nil {
		cs = &callers{}
		cs.stack = parseSlow(pcs[:n])
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
		cacheStack.Set(key, cs)
	}
	pool.Put(pcs)
	c.cache = cs
	return
}

//Clone 利用 code 和 msg 生成一个包含当前stack的新Error,
func Clone(err error, skips ...int) error {
	skip := 1
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

//New 替换 errors.New
func New(format string, a ...interface{}) error {
	if len(a) > 0 {
		format = fmt.Sprintf(format, a...)
	}
	return NewCode(1, DefaultCode, format)
}

//Errorf 替换 fmt.Errorf
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

//Error error interface, 序列化为string, 包含调用栈
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
func (e *Code) fmt() (cs fmtCode) {
	return fmtCode{code: strconv.Itoa(e.code), msg: e.msg, callers: e.cache}
}

type callers struct {
	stack []string
	attr  uint64 // count:escape ==> uint32:uint32
}
type fmtCode struct {
	code      string
	msg       string
	msgEscape bool
	*callers
}

func (f *fmtCode) jsonSize() (l int) {
	l, f.msgEscape = countEscape(f.msg)
	l += len(f.code) + len(`{"code":,"msg":""}`)
	if len(f.stack) == 0 {
		return
	}
	l += len(`,"stack":[]`) + len(f.stack)*len(`,""`) - len(`,`) + (int(f.attr) >> 32)
	return
}

func (f *fmtCode) textSize() (l int) {
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

func (f *fmtCode) text(buf *writeBuffer) {
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
