package errors

import (
	stderrs "errors"
	"fmt"
)

const (
	DefaultCode = -1
	DefaultMsg  = ""
)

type BizErr interface {
	error
	GetCode() int
	GetMsg() string
}
type Error interface {
	BizErr

	Is(err error) bool
}

var _ Error = &Cause{}

//TODO: 兼容 errors 和 pkg/errors
//兼容 erors.IS?

//New 用于替换 errors.New()
func New(msg string) Error {
	return buildCause(DefaultCode, msg, buildStack(1))
}

//Errorf 替换 fmt.Errorf
func Errorf(format string, a ...interface{}) Error {
	return buildCause(DefaultCode, fmt.Sprintf(format, a...), buildStack(1))
}

func Errf(code int, format string, a ...interface{}) Error {
	return buildCause(code, fmt.Sprintf(format, a...), buildStack(1))
}

func NewErr(code int, msg string) Error {
	return buildCause(code, msg, buildStack(1))
}

func NewErrSkip(skip, code int, msg string) (err *Cause) {
	return buildCause(code, msg, buildStack(1+skip))
}

func NewErrfSkip(skip, code int, format string, a ...interface{}) (err *Cause) {
	return buildCause(code, fmt.Sprintf(format, a...), buildStack(1+skip))
}

func CloneSkip(skip int, err error) *Cause {
	e, ok := err.(*Cause)
	if !ok {
		return e
	}
	return as(err, buildStack(1+skip))
}

//Clone 利用 err.Code 和 Cause.Msg 生成一个包含当前stack的新Error,
func Clone(err error) Error {
	e, ok := err.(*Cause)
	if ok {
		return e
	}
	return as(err, buildStack(1))
}

//As std error --> Error
func As(err error) (e *Cause) {
	if e, ok := err.(*Cause); ok {
		return e
	}
	return as(err, buildStack(1))
}

func as(err error, stack stack) (e *Cause) {
	code, msg := DefaultCode, err.Error()
	if bizErr, ok := err.(BizErr); ok {
		code, msg = bizErr.GetCode(), bizErr.GetMsg()
	} else if pcErr, ok := err.(pcErr); ok {
		code, msg = int(pcErr.GetBizStatusCode()), pcErr.GetBizStatusMessage()
	} else {
		if traceErr, ok := err.(TraceErr); ok {
			code = traceErr.GetErrorCode()
		}
		if _, ok := err.(fmt.Formatter); ok {
			msg = fmt.Sprintf("%+v", err)
		}
	}
	return buildCause(code, msg, stack)
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
		return e.GetMsg()
	}
	return DefaultMsg
}

//Is 检查code是不是一样的
func Is(err1, err2 error) bool {
	e1, ok := err1.(*Cause)
	if !ok || e1.GetCode() == DefaultCode {
		return stderrs.Is(err1, err2)
	}
	return e1.Is(err2)
}
