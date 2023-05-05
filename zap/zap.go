// MIT License
//
// Copyright (c) 2021 Xiantu Li
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package zap

import (
	"unsafe"

	"github.com/lxt1045/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zap.Logger
}

func toZapLogger(logger *Logger) *zap.Logger {
	return &logger.Logger
}

func toLogger(logger *zap.Logger) *Logger {
	return (*Logger)(unsafe.Pointer(logger))
}

func (log *Logger) getZapCore() zapcore.Core {
	return *(*zapcore.Core)(unsafe.Pointer(log))
}

func New(core zapcore.Core, options ...zap.Option) *Logger {
	// options = append(options, zap.WithCaller(false))
	logger := zap.New(core, options...)
	return toLogger(logger)
}

func (log *Logger) Log(lvl zapcore.Level, msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(lvl) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Log(lvl, msg, fields...)
}

func (log *Logger) Debug(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.DebugLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Debug(msg, fields...)
}

func (log *Logger) Info(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.InfoLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Info(msg, fields...)
}

func (log *Logger) Warn(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.WarnLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Warn(msg, fields...)
}

func (log *Logger) Error(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.ErrorLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Error(msg, fields...)
}

func (log *Logger) DPanic(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.DPanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.DPanic(msg, fields...)
}

func (log *Logger) Panic(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.PanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Panic(msg, fields...)
}

func (log *Logger) Fatal(msg string, fields ...zap.Field) {
	if !log.getZapCore().Enabled(zap.FatalLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	fields = append(fields, zap.String("caller", c.File))
	log.Logger.Fatal(msg, fields...)
}
