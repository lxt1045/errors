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
	Code() int
	Msg() string
}
type Error interface {
	BizErr

	Is(err error) bool
}

var _ Error = &Cause{}

//TODO: 兼容 errors 和 pkg/errors
//兼容 erors.IS?

func GetCodeMsg(err error) (code int, msg string) {
	code, msg = DefaultCode, err.Error()
	if bizErr, ok := err.(BizErr); ok {
		code, msg = bizErr.Code(), bizErr.Msg()
	} else {
		if e := stderrs.Unwrap(err); e != nil {
			return GetCodeMsg(e)
		}
		if _, ok := err.(fmt.Formatter); ok {
			msg = fmt.Sprintf("%+v", err)
		}
	}
	return
}

//Is 检查code是不是一样的
func Is(err1, target error) bool {
	e1, ok := err1.(*Cause)
	if !ok || e1.Code() == DefaultCode {
		return stderrs.Is(err1, target)
	}
	return e1.Is(target)
}
