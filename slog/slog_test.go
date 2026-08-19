package slog

import (
	"io"
	"log/slog"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("std", func(t *testing.T) {
		logger := slog.New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug, // slog记录所有的日志
				AddSource: true,            // 显示文件行号
			}))
		logger.Error(
			"some log messages",
			"string", `some string format log information`,
			"int", 3,
		)
	})
	t.Run("std-lxt", func(t *testing.T) {
		logger := New(slog.NewTextHandler(os.Stdout,
			&slog.HandlerOptions{
				Level:     slog.LevelDebug, // slog记录所有的日志
				AddSource: true,            // 显示文件行号
			}))
		logger.With().Error(
			"some log messages",
			"string", `some string format log information`,
			"int", 3,
		)
	})
}
func BenchmarkLog(b *testing.B) {
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

	b.Run("std+lxt+caller", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(slog.NewTextHandler(io.Discard,
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
	b.Run("std+lxt+caller1", func(b *testing.B) {
		b.StopTimer()
		b.ReportAllocs()
		logger := New(slog.NewTextHandler(io.Discard,
			&slog.HandlerOptions{
				Level: slog.LevelDebug, // slog记录所有的日志
				// AddSource: true,            // 显示文件行号
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
