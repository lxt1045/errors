package jmp

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func TestExample(t *testing.T) {
	ctx := context.TODO()
	log := zerolog.New(os.Stdout)

	err := func(ctx context.Context) (err error) {
		pc, err := Set()
		if err != nil {
			return
		}
		err = func() (err error) {
			return nil
		}()
		Try(pc, err)

		err = func() (err error) {
			err = errors.New("err")
			log.Info().Caller().Err(err).Send()
			return
		}()
		Try(pc, err)
		log.Info().Caller().Err(err).Msg("after Try(pc,err)")

		return
	}(ctx)

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

func TestSet(t *testing.T) {
	if !testing.Verbose() {
		t.Error("err")
	}

	t.Run("Set", func(t *testing.T) {
		func() {
			defer func() {
				t.Log("outer defer")
			}()
			err := func() (err error) {
				defer func() {
					t.Log("inner defer")
				}()
				t.Log("1")
				pc, err := Set()
				t.Logf("defer:0x%x, defer_offset:%d, jump:%+v", pc._defer, defer_offset, pc)
				if err != nil {
					t.Log("3")
					t.Log("Setjmp() get error:", err)
					return
				}
				defer func() {
					t.Log("inner defer 2")
				}()
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				Try(pc, nil)
				t.Log("4")
				Try(pc, err3)
				t.Log("5")
				return
			}()
			t.Log("outer err:", err)
		}()
	})

	t.Run("jump-deep", func(t *testing.T) {
		func() {
			defer func() {
				t.Log("outer defer")
			}()
			err := func() (err error) {
				defer func() {
					t.Log("inner defer")
				}()
				t.Log("1")
				pc, err := Set()
				t.Logf("defer:0x%x, defer_offset:%d, jump:%+v", pc._defer, defer_offset, pc)
				if err != nil {
					t.Log("3")
					t.Log("Setjmp() get error:", err)
					return
				}
				defer func() {
					t.Log("inner defer 2")
				}()
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				Try(pc, nil)
				t.Log("4")
				// TryLong(pc, err3)
				func() {
					func() {
						pc1, _ := Set()
						_ = pc1
						err4 := fmt.Errorf("error 4")
						t.Log("5")
						TryLong(pc, err4)
					}()
				}()
				TryLong(pc, err3)
				t.Log("6")
				Try(pc, err3)
				t.Log("7")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}
