package zerolog

import (
	"io"
	"os"
	"testing"

	"github.com/lxt1045/errors"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
)

func TestStdLog(t *testing.T) {
	t.Run("lxt-zerolog", func(t *testing.T) {
		logger := New(os.Stdout)
		logger.Info().
			Str("string", `some string format log information`).
			Timestamp().
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("lxt-zerolog-caller", func(t *testing.T) {
		logger := New(os.Stdout)
		logger.Info().
			Caller().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("lxt-zerolog-std-1", func(t *testing.T) {
		logger := New(os.Stdout).ToStd()
		logger.Info("string", `some string format log information`,
			"int", 3,
			"some log messages")
	})
	t.Run("lxt-zerolog-std-2", func(t *testing.T) {
		Info("string", `some string format log information`,
			"int", 3,
			"some log messages")
	})
}

func BenchmarkStdLog(b *testing.B) {
	b.Run("zerolog+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := zerolog.New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info().
				Caller().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})
	b.Run("zerolog+lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info().
				Caller().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})
	b.Run("zerolog+lxt+std caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(io.Discard).ToStd()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("string", `some string format log information`,
				"int", 3,
				"some log messages")
		}
	})

	b.Run("lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c := errors.GetPC().CallerFrame()
			io.Discard.Write([]byte(zap.String("caller", c.FileLine).String))
		}
	})
}
