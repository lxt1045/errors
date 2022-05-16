package errors

import (
	"testing"
)

var (
	bizErr1 = NewErr(111, "msg1")
)

func getBizErr() error {
	return NewErr(88, "msg")
}

func testCatchErr() (errRet error) {
	defer TryCatchErr(&errRet)()

	err := getBizErr()
	MustNil(err, bizErr1)
	return
}

func testCatchErr1() (errRet error) {
	defer TryCatchErr(&errRet)()

	err := getBizErr()
	MustNilErr(err)
	return
}

func TestTryCatch(t *testing.T) {
	err := testCatchErr()
	t.Logf("testCatchErr:%+v", err)

	err = testCatchErr1()
	t.Logf("testCatchErr1:%+v", err)
}
