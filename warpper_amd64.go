//go:build amd64
// +build amd64

package errors

import (
	_ "unsafe"
)

//go:noinline
//Wrap 需要禁止内联，因为内联后，将无法通过 BP 获取正确的PC。
func Wrap(err error, format string, ifaces ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrapper{
		pc:     GetPC(),
		err:    err,
		format: format,
		ifaces: ifaces,
	}
}

//go:noinline
func NewLine(format string, ifaces ...interface{}) error {
	return &wrapper{
		pc:     GetPC(),
		err:    nil,
		format: format,
		ifaces: ifaces,
	}
}
