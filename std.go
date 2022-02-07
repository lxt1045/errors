package errors

import (
	stderrs "errors"
	"fmt"
)

type BizErr interface {
	error
	GetCode() int
	GetMessage() string
}
type Error interface {
	BizErr

	Is(err error) bool
	Wrap(l string) Error
	Wrapf(format string, a ...interface{}) Error
	Unwrap() error
}

//Wrap ：添加中间栈帧（知乎记录当前的文件位置，不记录caller）；如果是std error则先转为Err
func Wrap(err error, trace string) error {
	return as(1, err).wrap(trace)
}

func Wrapf(err error, format string, a ...interface{}) error {
	return as(1, err).wrap(fmt.Sprintf(format, a...))
}

//Clone 利用 err.Code 和 Err.Message 生成一个包含当前stack的新Error,
func Clone(err error, trace string) Error {
	return as(1, err).clone(1, trace)
}

func Clonef(err error, format string, a ...interface{}) Error {
	return as(1, err).clone(1, fmt.Sprintf(format, a...))
}

//As std error --> Error
func As(err error) Error {
	return as(1, err)
}

//GetCode ...
func GetCode(err error) int {
	return as(1, err).Code
}

//GetMsg ...
func GetMsg(err error) string {
	return as(1, err).Message
}

//Is 检查code是不是一样的
func Is(err1, err2 error) bool {
	e1 := as(1, err1)
	if e1.Code == -1 {
		return stderrs.Is(err1, err2)
	}
	return e1.Is(err2)
}
