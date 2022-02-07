package errors

import (
	"fmt"
	"strconv"
	"strings"
)

var Layout = LayoutTypeText

type LayoutType int8

const (
	LayoutTypeJSON LayoutType = 1 // 以json格式输出
	LayoutTypeText LayoutType = 2 // 已text格式输出
)

//New 用于替换 errors.New()
func New(msg string) Error {
	return &Err{
		Code:    -1,
		Message: msg,
		cause:   newBreadcrumb(1, DefaultLevel, ""),
	}
}

//Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) Error {
	msg := fmt.Sprintf(format, a...)
	return &Err{
		Code:    -1,
		Message: msg,
		cause:   newBreadcrumb(1, DefaultLevel, ""),
	}
}

func Errf(code int, format string, a ...interface{}) Error {
	msg := fmt.Sprintf(format, a...)
	return &Err{
		Code:    code,
		Message: msg,
		cause:   newBreadcrumb(1, DefaultLevel, ""),
	}
}
func newErr(skip, code int, msg string) Error {
	return &Err{
		Code:    code,
		Message: msg,
		cause:   newBreadcrumb(1+skip, DefaultLevel, ""),
	}
}
func newErrf(skip, code int, format string, a ...interface{}) Error {
	msg := fmt.Sprintf(format, a...)
	return &Err{
		Code:    code,
		Message: msg,
		cause:   newBreadcrumb(1+skip, DefaultLevel, ""),
	}
}

func NewErr(code int, msg string, trace string) Error {
	return &Err{
		Code:    code,
		Message: msg,
		cause:   newBreadcrumb(1, DefaultLevel, trace),
	}
}

func (e Err) clone(skip int, trace string) Error {
	return &Err{
		Code:    e.Code,
		Message: e.Message,
		cause:   newBreadcrumb(1+skip, DefaultLevel, trace),
	}
}

func as(skip int, err error) *Err {
	if e, ok := err.(*Err); ok {
		return e
	}
	s := "nil error"
	if err != nil {
		s = err.Error()
	}
	return &Err{
		Code:    -1,
		Message: s,
		cause:   newBreadcrumb(1+skip, DefaultLevel, ""),
	}
}

//Err 一个 wrap error,可以打印 NewErr 时的调用栈
type Err struct {
	Code    int    //业务错误码
	Message string //业务错误信息

	cause  breadcrumb   //错误的现场
	traces []breadcrumb //Warp()组成的路径
}

func (e Err) GetCode() int {
	return e.Code
}
func (e Err) GetMessage() string {
	return e.Message
}

func (e *Err) Wrap(trace string) Error {
	return e.wrap(trace)
}
func (e *Err) Wrapf(format string, a ...interface{}) Error {
	return e.wrap(fmt.Sprintf(format, a...))
}
func (e *Err) wrap(trace string) Error {
	e.traces = append(e.traces, newBreadcrumb(2, 1, trace))
	return e
}
func (e *Err) Unwrap() error {
	if len(e.traces) <= 0 {
		return e
	}
	e.traces = e.traces[:len(e.traces)-1]
	return e
}
func (e *Err) Is(err error) bool {
	to, ok := err.(Error)
	if !ok {
		return false
	}
	if e.GetCode() != -1 && e.GetCode() == to.GetCode() {
		return true
	}
	return false
}

//Error error interface, 序列化为string, 包含调用栈
func (e *Err) Error() string {
	return e.serialize()
}

func (e *Err) String() string {
	return e.serialize()
}

func (e *Err) serialize() string {
	buf := new(strings.Builder)
	if Layout == LayoutTypeJSON {
		return e.json(buf).String()
	}
	return e.text(buf).String()
}

func (e *Err) json(buf *strings.Builder) *strings.Builder {
	buf.Write([]byte(`{"code":`))
	buf.WriteString(strconv.Itoa(e.Code))
	buf.Write([]byte(`,"message":`))
	buf.WriteByte('"')
	buf.WriteString(e.Message)
	buf.WriteByte('"')
	buf.Write([]byte(`,"cause":`))
	e.cause.json(buf)
	if len(e.traces) > 0 {
		buf.Write([]byte(`,"traces":[`))
		for i, b := range e.traces {
			if i != 0 {
				buf.WriteByte(',')
			}
			b.json(buf)
		}
		buf.WriteByte(']')
	}
	buf.WriteByte('}')

	return buf
}

func (e *Err) text(buf *strings.Builder) *strings.Builder {
	buf.WriteString(strconv.Itoa(e.Code))
	buf.Write([]byte(`, `))
	if e.Message == "" {
		e.Message = "-"
	}
	buf.WriteString(e.Message)
	buf.Write([]byte("\n"))

	e.cause.text(buf)
	buf.Write([]byte(";\n"))

	for _, b := range e.traces {
		b.text(buf)
		buf.Write([]byte(";\n"))
	}

	return buf
}
