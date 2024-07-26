package jmp_test

import (
	"context"
	"os"
	"testing"

	"github.com/lxt1045/errors"
	"github.com/lxt1045/errors/jmp"
	"github.com/rs/zerolog"
)

func Do(ctx context.Context) (err error) {
	pc, err := jmp.Set()
	if err != nil {
		return
	}
	err = func() (err error) {
		return nil
	}()
	jmp.Try(pc, err)

	err = func() (err error) {
		err = errors.New("err")
		return
	}()
	jmp.Try(pc, err)
	return
}

func TestExample(t *testing.T) {
	ctx := context.TODO()
	log := zerolog.New(os.Stdout)

	err := Do(ctx)

	log.Info().Caller().Err(err).Msg("just for test")

}

var _ = func() bool {
	for _, v := range os.Args {
		if v == "-v" || v == "-test.v" {
			return true
		}
	}
	os.Args = append(os.Args, "-test.v")
	testing.Init() // 根据 go 全局变量的初始化顺序，全局变量优先init()函数执行
	return true
}()
