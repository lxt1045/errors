package errors

import (
	"fmt"
	"reflect"
)

type pcErr interface {
	GetBizStatusCode() int32     // 获取业务的状态码
	GetBizStatusMessage() string // 获取业务的状态信息
}
type TraceErr interface {
	GetErrorCode() int // 获取框架的错误码
}

func OK(ok bool, err *Cause) {
	if ok {
		return
	}
	if err != nil {
		panic(err) //重新生成调用栈
	}
	panic(&Cause{
		Code:  DefaultCode,
		Msg:   "not ok",
		stack: buildStack(1),
	})
}

func NilErr(err error) {
	if IsNil(err) {
		return
	}
	e, ok := err.(*Cause)
	if !ok {
		e = as(err, buildStack(1))
	}
	panic(e)
}

func Nil(obj interface{}, err *Cause) {
	if IsNil(obj) {
		return
	}
	if err != nil {
		panic(err) //重新生成调用栈
	}
	panic(&Cause{
		Code:  DefaultCode,
		Msg:   "not nil",
		stack: buildStack(1),
	})
}

func Nilf(obj interface{}, code int, format string, a ...interface{}) {
	if IsNil(obj) {
		return
	}
	panic(&Cause{
		Code:  code,
		Msg:   fmt.Sprintf(format, a...),
		stack: buildStack(1),
	})
}

func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}
	val := reflect.ValueOf(object)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
		reflect.Slice, reflect.Interface:
		return val.IsNil()
	default:
		return false
	}
}

func TryByFunc(fCatch func(interface{}) bool) {
	e := recover()
	if e == nil {
		return
	}
	if fCatch != nil && fCatch(e) {
		return
	}
	panic(e)
}

func TryErr(perr *error) {
	e := recover()
	if e == nil {
		return
	}
	if perr == nil {
		panic(e)
	}
	ok := true
	if *perr, ok = e.(*Cause); ok {
		return
	}
	panic(e)
}
