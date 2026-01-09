//go:build !amd64 && amd64p32 && arm64
// +build !amd64,amd64p32,arm64

package errors

var NewCode = NewCodeSlow
