package errors

import (
	stderrs "errors"
)

const (
	DefaultCode = -1
	DefaultMsg  = ""
)

// Is 检查code是不是一样的
func Is(err1, target error) bool {
	e1, ok := err1.(*Code)
	if !ok || e1.Code() == DefaultCode {
		return stderrs.Is(err1, target)
	}
	return e1.Is(target)
}

func WithErr(err error) (e *Code) {
	e, ok := err.(*Code)
	if ok {
		return
	}
	return NewCode(1, DefaultCode, err.Error())
}

func AsCode(err error) (e *Code) {
	e, ok := err.(*Code)
	if ok {
		return
	}
	return &Code{
		code: DefaultCode,
		msg:  err.Error(),
	}
}
