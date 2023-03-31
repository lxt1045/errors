package errors_test

import (
	"runtime"
	"strconv"
	"testing"

	"github.com/lxt1045/errors"
	"github.com/lxt1045/errors/logrus"
)

// sample 1
func LineByRuntime() string {
	_, file, n, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file + ":" + strconv.Itoa(n)
}

func TestSample(t *testing.T) {
	t.Run("LineByRuntime", func(t *testing.T) {
		t.Log(runtime.Caller(0))
	})
	t.Run("Line", func(t *testing.T) {
		t.Log(logrus.CallerFrame(errors.GetPC()))
	})
}

var gLine string

func BenchmarkUnmarshalInterface(b *testing.B) {

	b.Run("LineByRuntime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, gLine, _, _ = runtime.Caller(0)
		}
	})

	b.Run("Line", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gLine = logrus.CallerFrame(errors.GetPC()).File
		}
	})
}
