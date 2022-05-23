package errors

import (
	"strconv"
	"unsafe"
)

//Cause 是一个带错误栈信息的error
type Cause struct {
	code int    //业务错误码
	msg  string //业务错误信息

	stack stack //错误的现场
}

func buildCause(code int, msg string, stack stack) (e *Cause) {
	return &Cause{
		code:  code,
		msg:   msg,
		stack: stack,
	}
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
	bs := make([]byte, 0, cache.textSize())
	bs = cache.text(bs)
	return *(*string)(unsafe.Pointer(&bs))
}

// MarshalJSON json.Marshaler的方法, json.Marshal 里调用
func (e *Cause) MarshalJSON() (bs []byte, err error) {
	cache := e.fmt()
	bs = make([]byte, 0, cache.jsonSize())
	bs = cache.json(bs)
	return
}

func (e *Cause) fmt() (cs fmtCause) {
	return fmtCause{strconv.Itoa(e.code), e.msg, e.stack.fmt()}
}

type fmtCause struct {
	code       string
	msg        string
	stackCache fmtStack
}

func (f *fmtCause) jsonSize() int {
	return 8 + len(f.code) + 8 + 1 + len(f.msg) + len(f.stackCache)*9 + 1 + f.stackCache.jsonSize()
}
func (f *fmtCause) textSize() int {
	return 2 + len(f.code) + len(f.msg) + len(f.stackCache)*3 + f.stackCache.textSize()
}
func (f *fmtCause) json(bs []byte) []byte {
	bs = append(bs, `{"code":`...)
	bs = append(bs, f.code...)
	bs = append(bs, `,"msg":"`...)
	bs = append(bs, f.msg...)
	bs = append(bs, '"')
	if len(f.stackCache) > 0 {
		bs = append(bs, `,"stack":`...)
		bs = f.stackCache.json(bs)
	}
	bs = append(bs, '}')
	return bs
}

func (f *fmtCause) text(bs []byte) []byte {
	bs = append(bs, f.code...)
	bs = append(bs, ", "...)
	bs = append(bs, f.msg...)
	if len(f.stackCache) > 0 {
		bs = append(bs, ";\n"...)
		bs = f.stackCache.text(bs)
		bs = append(bs, ';')
	}
	return bs
}
