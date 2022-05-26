package errors

import (
	"fmt"
	"testing"

	pkgerrs "github.com/pkg/errors"
)

func TestMarshalJSON(t *testing.T) {
	for _, depth := range []int{0, 10} {
		name := fmt.Sprintf("%s-%d", "MarshalJSON", depth)
		t.Run(name, func(t *testing.T) {
			err := pkgerrs.New("test")
			_ = err
		})
	}
}
