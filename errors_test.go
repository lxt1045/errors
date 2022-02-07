package errors_test

import (
	"fmt"
	"testing"

	errs "github.com/lxt1045/errors"
)

func ferr1() error {
	err := errs.NewErr(1600002, "message ferr1", "log ferr1")
	return err
}

func ferr2() error {
	err := ferr1()
	err = errs.Wrap(err, "log ferr2")
	return err
}

func Test_New(t *testing.T) {
	err := errs.NewErr(1600002, "message ferr1", "log ferr1")
	fmt.Println("err:\n", err)
	errs.Layout = errs.LayoutTypeJSON
	fmt.Println(err)
}

func Test_Wrap(t *testing.T) {
	err := ferr2()
	err = errs.Wrap(err, "warp ...")
	fmt.Println("err:\n", err)
	errs.Layout = errs.LayoutTypeJSON
	fmt.Println(err)
}
