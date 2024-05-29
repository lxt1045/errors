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

func TestSet(t *testing.T) {
	t.Run("Set", func(t *testing.T) {
		// t.Error("err")
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
}
