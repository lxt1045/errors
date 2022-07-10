package errors

import (
	stderrs "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStd(t *testing.T) {
	t.Run("GetCodeMsg.Code", func(t *testing.T) {
		err := NewErr(errCode, errMsg)
		cause := err.(*Code)
		assert.Equal(t, cause.Code(), errCode)
		assert.Equal(t, cause.Msg(), errMsg)
	})
	t.Run("Is.base", func(t *testing.T) {
		err := NewErr(errCode, errMsg)
		err1 := NewErr(errCode, errMsg)
		assert.True(t, Is(err, err1))
		err2 := stderrs.New(errMsg)
		assert.False(t, Is(err, err2))
	})
	t.Run("Is.wrap", func(t *testing.T) {
		err := NewErr(errCode, errMsg)
		err2 := stderrs.New(errMsg)
		err3 := fmt.Errorf(errTrace+":%w", err2)
		assert.True(t, Is(err3, err2))
		err4 := fmt.Errorf(errTrace+":%w", err)
		assert.True(t, Is(err4, err))
		err5 := Wrap(err, errTrace)
		assert.True(t, Is(err5, err))
	})
}

//
