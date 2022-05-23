package main

import (
	"log"
	"unsafe"

	"github.com/lxt1045/errors"
)

func FErr() (err error) {
	defer TryCatchErr(&err)
	errors.NilErr(err)
	return
}

func TryCatchErr(perrr *error) {
	p := uintptr(unsafe.Pointer(perrr))
	e := recover()
	if e == nil {
		return
	}
	perr := (unsafe.Pointer)(p)
	if perr == nil {
		panic(e)
	}
	ok := true
	perr1 := (*error)(perr)
	if *perr1, ok = e.(*errors.Cause); ok {
		return
	}

	//其他错误则再次抛出
	panic(e)
}

func main() {
	defer func() {
		log.Println(recover())
	}()

	go func() {
		// Panic occurs in a goroutine
		panic("A bad boy stole a server")
	}()
	return
}
