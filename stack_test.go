package errors_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/lxt1045/errors"
)

func F1() {
	f := errors.NewStack(0, 0)
	fmt.Println(f)
}

func F2() {
	F1()
}

func F3() {
	F2()
}

func Test_NewFrames(t *testing.T) {
	f := errors.NewStack(0, 0)
	fmt.Println(f)

	go F3()
	time.Sleep(1 * time.Second)
}

func Test_NewStack(t *testing.T) {
	s := errors.NewStack(0, 0)
	fmt.Println(s.String())
	errors.Layout = errors.LayoutTypeJSON
	fmt.Println(s.String())
}
