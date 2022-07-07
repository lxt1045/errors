//go:build !amd64
// +build !amd64

package errors

import (
	"runtime"
)

func Wrap(err error, format string, is ...interface{}) error {
	if err == nil {
		return nil
	}
	e := &wrapper{
		err:    err,
		format: format,
		ifaces: is,
	}
	runtime.Callers(baseSkip, e.pc[:])
	return e
}

func NewLine(format string, ifaces ...interface{}) error {
	e := &wrapper{
		err:    nil,
		format: format,
		ifaces: ifaces,
	}
	runtime.Callers(baseSkip, e.pc[:])
	return e
}
