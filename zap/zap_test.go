package zap

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/lxt1045/errors"
	lxtlog "github.com/lxt1045/errors/logrus"
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLog(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			// zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			// zapcore.AddSync(io.Discard),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(true))
		logger.Info("some log messages",
			zap.String("string", `some string format log information`),
			zap.Int("int", 3),
		)
	})

	t.Run("2", func(t *testing.T) {
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			// zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			// zapcore.AddSync(io.Discard),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
		logger := New(core, zap.WithCaller(false))
		logger.Info("some log messages",
			zap.String("string", `some string format log information`),
			zap.Int("int", 3),
		)
	})

	t.Run("3", func(t *testing.T) {
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			// zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			// zapcore.AddSync(io.Discard),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
		logger := zap.New(core, zap.WithCaller(false))
		sugar := logger.Sugar()
		sugar.Info("some log messages",
			"string", `some string format log information`,
			"int", 3,
		)
	})

	t.Run("4", func(t *testing.T) {
		cfg := zap.NewProductionConfig()
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			// zapcore.NewConsoleEncoder(cfg.EncoderConfig),
			// zapcore.AddSync(io.Discard),
			zapcore.AddSync(os.Stdout),
			zapcore.InfoLevel,
		)
		logger := New(core, zap.WithCaller(false))
		sugar := logger.Sugar()
		sugar.Info("some log messages",
			"string", `some string format log information`,
			"int", 3,
			"backoff", time.Second,
		)
	})
}

func BenchmarkLog(b *testing.B) {
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
		logger := New(core, zap.WithCaller(false))
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
		logger := New(core, zap.WithCaller(false))
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
}

/*
go test -benchmem -run=^$ -bench ^BenchmarkZapCaller$ github.com/lxt1045/errors/zap -count=1 -v -cpuprofile cpu.prof -c
go tool pprof ./json.test cpu.prof
*/
func BenchmarkZapCaller(b *testing.B) {
	b.StopTimer()
	b.ReportAllocs()
	cfg := zap.NewProductionConfig()
	core := zapcore.NewCore(
		// zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(io.Discard),
		zapcore.InfoLevel,
	)
	logger := New(core, zap.WithCaller(false))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("some log messages",
			zap.String("string", `some string format log information`),
			zap.Int("int", 3),
		)
	}
}
