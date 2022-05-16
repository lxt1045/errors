package errors_test

import (
	"encoding/json"
	"fmt"
	"testing"

	errs "github.com/lxt1045/errors"
)

func ferr1() error {
	err := errs.NewErr(1600002, "message ferr1")
	return err
}

func ferr2() error {
	err := ferr1()
	err = errs.Wrap(err, "log ferr2")
	return err
}

func Test_New(t *testing.T) {
	err := errs.NewErr(1600002, "message ferr1")
	fmt.Println("err:\n", err)
}

func Test_Wrap(t *testing.T) {
	err := ferr2()
	err = errs.Wrap(err, "warp ...")
	fmt.Println("err:\n", err)
}

func Test_JSON(t *testing.T) {
	err := ferr2()
	err = errs.Wrap(err, "warp1 ...")
	err = errs.Wrap(err, "warp2 ...")

	bs, e := json.Marshal(err)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(bs))
}

type X struct {
}

func (X) Err() error {
	err := ferr2()
	err = errs.Wrap(err, "warp1 ...")
	err = errs.Wrap(err, "warp2 ...")
	return err
}

func Test_JSON2(t *testing.T) {
	err := X{}.Err()
	bs, e := json.Marshal(err)
	if e != nil {
		t.Fatal(e)
	}
	fmt.Println(string(bs))
}
