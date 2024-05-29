package zerolog

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/lxt1045/errors"
	lxtlog "github.com/lxt1045/errors/logrus"
	lxtzaplog "github.com/lxt1045/errors/zap"
	"github.com/rs/zerolog"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLog(t *testing.T) {
	t.Run("fatal-zerolog", func(t *testing.T) {
		defer func() {
			t.Log("in defer")
		}()
		logger := zerolog.New(os.Stdout)
		logger.Info().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("panic-zerolog", func(t *testing.T) {
		defer func() {
			if e := recover(); e != nil {
				t.Log("in defer:", e)
			}
		}()
		logger := zerolog.New(os.Stdout)
		logger.Panic().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("lxt-zerolog", func(t *testing.T) {
		logger := New(os.Stdout)
		logger.Info().
			Str("string", `some string format log information`).
			Timestamp().
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("zerolog", func(t *testing.T) {
		logger := zerolog.New(os.Stdout)
		logger.Info().
			Str("string", `some string format log information`).
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
	t.Run("zerolog-caller", func(t *testing.T) {
		logger := zerolog.New(os.Stdout)
		logger.Info().
			Caller().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})

	//

	t.Run("lxt-zerolog-context-caller", func(t *testing.T) {
		logger := New(os.Stdout)
		log := logger.
			With().
			Caller().Logger()
		log.Info().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})
	t.Run("zerolog-context-caller", func(t *testing.T) {
		logger := zerolog.New(os.Stdout)
		log := logger.
			With().
			Caller().Logger()
		log.Info().
			Str("string", `some string format log information`).
			Int("int", 3).
			Msg("some log messages")
	})
}

func BenchmarkLog(b *testing.B) {
	b.Run("zerolog", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := zerolog.New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})
	b.Run("zerolog+lxt", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})
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

	b.Run("zerolog+context-caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := zerolog.New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			log := logger.
				With().
				Caller().Logger()
			log.Info().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})
	b.Run("zerolog+lxt context-caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			log := logger.
				With().
				Caller().Logger()
			log.Info().
				Str("string", `some string format log information`).
				Int("int", 3).
				Msg("some log messages")
		}
	})

	b.Run("logrus", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		// logrus.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(logrus.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})
	b.Run("logrus+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := logrus.New()
		logger.SetOutput(io.Discard)
		logger.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(logrus.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})
	b.Run("logrus+lxt caller", func(b *testing.B) {
		// logrus.SetReportCaller(false)
		b.StopTimer()
		b.ReportAllocs()
		logger := lxtlog.New()
		logger.SetOutput(io.Discard)
		// logrus.SetReportCaller(true)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.WithFields(lxtlog.Fields{
				"string": "some string format log information",
				"int":    3,
			}).Info("some log messages")
		}
	})

	b.Run("zap", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})
	b.Run("zap+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(true))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})
	b.Run("zap+lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := lxtzaplog.New(core, zap.WithCaller(false))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info("some log messages",
				zap.String("string", `some string format log information`),
				zap.Int("int", 3),
			)
		}
	})

	b.Run("zap-sugar", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core)
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})
	b.Run("zap-sugar+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(true))
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})

	b.Run("zap-sugar+lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			// zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			zapcore.AddSync(io.Discard),
			zapcore.InfoLevel,
		)
		logger := lxtzaplog.New(core, zap.WithCaller(false))
		sugar := logger.Sugar()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			sugar.Info("some log messages",
				"string", `some string format log information`,
				"int", 3,
				"backoff", time.Second,
			)
		}
	})

	b.Run("lxt caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			c := CallerFrame(errors.GetPC())
			io.Discard.Write([]byte(zap.String("caller", c.File).String))
		}
	})
}

/*
go test -benchmem -run=^$ -bench ^BenchmarkZeroCaller$ github.com/lxt1045/errors/zerolog -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./json.test cpu.prof
*/
func BenchmarkZeroCaller(b *testing.B) {
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
}
