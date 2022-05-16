package errors

import (
	"fmt"
	"reflect"
)

type RPCErr interface {
	GetBizStatusCode() int32     // 获取业务的状态码
	GetBizStatusMessage() string // 获取业务的状态信息
}
type TraceErr interface {
	GetErrorCode() int // 获取框架的错误码
}

func Must(ok bool) {
	if ok {
		return
	}
	throwErr(1, NewErrSkip(1, -1, "not ok"))
}

func MustNilErr(err error) {
	if IsNil(err) {
		return
	}
	throwErr(1, err)
}

func MustNil(obj interface{}, bizErr Error) {
	if IsNil(obj) {
		return
	}
	if bizErr != nil {
		throwErr(1, CloneSkip(1, bizErr)) // 重建调用栈)
	}
	err, ok := obj.(error)
	if !ok {
		err = newErrf(1, DefaultCode, "%+v", obj)
	}
	throwErr(1, err)
}

func throwErr(skip int, err error) {
	if traceErr, ok := err.(*Err); ok {
		panic(traceErr) //重新生成调用栈
	}
	if bizErr, ok := err.(Error); ok {
		panic(bizErr)
	}
	if rpcErr, ok := err.(RPCErr); ok {
		panic(NewErrSkip(skip+1, int(rpcErr.GetBizStatusCode()), rpcErr.GetBizStatusMessage()))
	}
	if traceErr, ok := err.(TraceErr); ok {
		panic(NewErrSkip(skip+1, traceErr.GetErrorCode(), err.Error()))
	}
	panic(newErrf(skip+1, DefaultCode, "%+v", err))
}

func MustNilf(obj interface{}, code int, format string, a ...interface{}) {
	if IsNil(obj) {
		return
	}
	panic(newErrf(1, code, format, a...))
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

func TryCatch(fCatch func(err Error)) (deferFunc func()) {
	if fCatch == nil {
		panic("errors.TryCatch(fCatch), fCatch==nil")
	}
	return func() {
		err := recover()
		if err == nil {
			return
		}
		if bizErr, ok := err.(Error); ok {
			fCatch(bizErr)
			return
		}

		//其他错误则再次抛出
		panic(err)
	}
}

func TryCatchErr(perr *error) (deferFunc func()) {
	if perr == nil {
		panic("errors.TryCatchErr(perr), perr==nil")
	}
	return func() {
		rerr := recover()
		if rerr == nil {
			return
		}
		if bizErr, ok := rerr.(error); ok {
			*perr = bizErr
			return
		}

		//其他错误则再次抛出
		panic(perr)
	}
}

func newErrf(skip, code int, format string, a ...interface{}) Error {
	msg := fmt.Sprintf(format, a...)
	return &Err{
		Code:  code,
		Msg:   msg,
		cause: buildStack(1 + skip),
	}
}
