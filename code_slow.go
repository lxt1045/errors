//go:build (!amd64 && !amd64p32 && !arm64) || cgo
// +build !amd64,!amd64p32,!arm64 cgo

package errors

var NewCode = NewCodeSlow
