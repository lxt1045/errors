package zerolog

import (
	"io"
	"log/slog"
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

func TestStr(t *testing.T) {
	// bs := [...]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17}
	// const M = "\001\002\003\004\005\006\007\010\011\012"
	const M = "\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10"
	// M := string(bs)
	t.Logf("%X", M)
	for i := 0; i < len(M); i++ {
		t.Logf("%d:%0x", i, byte(M[i]))
		t.Logf("%d:%0d", i, byte(M[i]))
	}
}
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

	t.Run("lxt-zerolog-context-caller", func(t *testing.T) {
		logger := New(os.Stdout)
		log := logger.
			With().
			Caller().Logger()
		log.Info().
			Str("string", `some string format log information`).
			Int("int", 3).
			Send()
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
	b.Run("zerolog+lxt std context-caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(io.Discard)
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			log := logger.With().Logger()
			log.Infoln(
				"string", `some string format log information`,
				"int", 3,
				"some log messages",
			)
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
			c := errors.GetPC().CallerFrame()
			io.Discard.Write([]byte(zap.String("caller", c.FileLine).String))
		}
	})

	b.Run("std+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := slog.New(slog.NewTextHandler(io.Discard,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug, // slog记录所有的日志
				AddSource: true,            // 显示文件行号
			}))
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			logger.Info(
				"some log messages",
				"string", `some string format log information`,
				"int", 3,
			)
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
