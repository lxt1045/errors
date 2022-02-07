package errors

import (
	"reflect"
)

type RPCErr interface {
	GetBizStatusCode() int32     // 获取业务的状态码
	GetBizStatusMessage() string // 获取业务的状态信息
}
type FrameErr interface {
	GetErrorCode() int // 获取框架的错误码
}

const (
	DefaultCode = -1
	DefaultMsg  = ""
)

func Must(ok bool) {
	if ok {
		return
	}
	throwErr(1, newErr(1, -1, "not ok"))
}

func MustNilErr(err error) {
	if IsNil(err) {
		return
	}
	throwErr(1, err)
}

func MustNil(obj interface{}, bizErr BizErr) {
	if IsNil(obj) {
		return
	}
	if bizErr != nil {
		if err, ok := bizErr.(*Err); ok {
			bizErr = err.clone(1, "") // 重建调用栈
		}
		throwErr(1, bizErr)
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
	if bizErr, ok := err.(BizErr); ok {
		panic(bizErr)
	}
	if rpcErr, ok := err.(RPCErr); ok {
		panic(newErr(skip+1, int(rpcErr.GetBizStatusCode()), rpcErr.GetBizStatusMessage()))
	}
	if frameErr, ok := err.(FrameErr); ok {
		panic(newErr(skip+1, frameErr.GetErrorCode(), err.Error()))
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

func TryCatch(fCatch func(err BizErr)) (deferFunc func()) {
	if fCatch == nil {
		panic("errors.TryCatch(fCatch), fCatch==nil")
	}
	return func() {
		err := recover()
		if err == nil {
			return
		}
		if bizErr, ok := err.(BizErr); ok {
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
