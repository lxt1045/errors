package errors

import (
	stderrs "errors"
	"fmt"
)

const (
	DefaultCode = -1
	DefaultMsg  = ""
)

type Error interface {
	error
	GetCode() int
	GetMessage() string

	Is(err error) bool
	Wrap(l string) Error
	Wrapf(format string, a ...interface{}) Error
	Unwrap() error
}

var _ Error = &Err{}

//New 用于替换 errors.New()
func New(msg string) Error {
	return &Err{
		Code:  DefaultCode,
		Msg:   msg,
		cause: buildStack(1),
	}
}

//Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) Error {
	return &Err{
		Code:  DefaultCode,
		Msg:   fmt.Sprintf(format, a...),
		cause: buildStack(1),
	}
}

func Errf(code int, format string, a ...interface{}) Error {
	return &Err{
		Code:  code,
		Msg:   fmt.Sprintf(format, a...),
		cause: buildStack(1),
	}
}

func NewErr(code int, msg string) Error {
	return &Err{
		Code:  code,
		Msg:   msg,
		cause: buildStack(1),
	}
}

func NewErrSkip(skip, code int, msg string) (err *Err) {
	err = &Err{
		Code:  code,
		Msg:   msg,
		cause: buildStack(skip + 1),
	}
	return
}

// Wrap ：添加中间栈帧（知乎记录当前的文件位置，不记录caller）；如果是std error则先转为Err
func Wrap(err error, trace string) Error {
	e, ok := err.(*Err)
	if !ok {
		e = &Err{
			Code:  DefaultCode,
			Msg:   err.Error(),
			cause: buildStack(1),
		}
	}
	e.traces = append(e.traces, buildFrame(1, trace))
	return e
}

func Wrapf(err error, format string, a ...interface{}) error {
	e, ok := err.(*Err)
	if !ok {
		e = &Err{
			Code:  DefaultCode,
			Msg:   err.Error(),
			cause: buildStack(1),
		}
	}
	e.traces = append(e.traces, buildFrame(1, fmt.Sprintf(format, a...)))
	return e
}

func CloneSkip(skip int, err error) Error {
	if e, ok := err.(*Err); ok {
		return &Err{
			Code:  e.GetCode(),
			Msg:   e.GetMessage(),
			cause: buildStack(skip + 1),
		}
	}
	return &Err{
		Code:  DefaultCode,
		Msg:   err.Error(),
		cause: buildStack(skip + 1),
	}
}

//Clone 利用 err.Code 和 Err.Msg 生成一个包含当前stack的新Error,
func Clone(err error) Error {
	if e, ok := err.(*Err); ok {
		return &Err{
			Code:  e.GetCode(),
			Msg:   e.GetMessage(),
			cause: buildStack(1),
		}
	}
	return &Err{
		Code:  DefaultCode,
		Msg:   err.Error(),
		cause: buildStack(1),
	}
}

//As std error --> Error
func As(err error) Error {
	if e, ok := err.(Error); ok {
		return e
	}
	s := "nil error"
	if err != nil {
		s = err.Error()
	}
	return &Err{
		Code:  DefaultCode,
		Msg:   s,
		cause: buildStack(1),
	}
}

//GetCode ...
func GetCode(err error) int {
	if e, ok := err.(Error); ok {
		return e.GetCode()
	}
	return DefaultCode
}

//GetMsg ...
func GetMsg(err error) string {
	if e, ok := err.(Error); ok {
		return e.GetMessage()
	}
	return DefaultMsg
}

//Is 检查code是不是一样的
func Is(err1, err2 error) bool {
	e1, ok := err1.(Error)
	if !ok || e1.GetCode() == DefaultCode {
		return stderrs.Is(err1, err2)
	}
	return e1.Is(err2)
}
