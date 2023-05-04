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
	"fmt"
	"unsafe"

	"github.com/lxt1045/errors"
	"go.uber.org/zap"
)

type SugaredLogger struct {
	zap.SugaredLogger
}

func toZapSugaredLogger(logger *SugaredLogger) *zap.SugaredLogger {
	return &logger.SugaredLogger
}

func toSugaredLogger(logger *zap.SugaredLogger) *SugaredLogger {
	return (*SugaredLogger)(unsafe.Pointer(logger))
}

func (log *SugaredLogger) getLogger() *Logger {
	return *(**Logger)(unsafe.Pointer(log))
}
func (log *Logger) Sugar() *SugaredLogger {
	sLogger := log.Logger.Sugar()
	return toSugaredLogger(sLogger)
}

func (s *SugaredLogger) Desugar() *Logger {
	base := s.SugaredLogger.Desugar()
	return toLogger(base)
}

func (s *SugaredLogger) Named(name string) *SugaredLogger {
	return toSugaredLogger(s.SugaredLogger.Named(name))
}

func (s *SugaredLogger) WithOptions(opts ...zap.Option) *SugaredLogger {
	return toSugaredLogger(s.SugaredLogger.WithOptions(opts...))
}

func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger {
	return toSugaredLogger(s.SugaredLogger.With(args...))
}

func getTemplateArgs(template string, fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return template
	}
	return fmt.Sprintf(template, fmtArgs...)
}

func getArgs(fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return ""
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	return fmt.Sprint(fmtArgs...)
}

func getArgsLn(fmtArgs []interface{}) string {
	if len(fmtArgs) == 0 {
		return ""
	}

	if len(fmtArgs) == 1 {
		if str, ok := fmtArgs[0].(string); ok {
			return str
		}
	}
	msg := fmt.Sprintln(fmtArgs...)
	return msg[:len(msg)-1]
}

func (s *SugaredLogger) Debug(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DebugLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Debug(msg, zap.String("caller", c.File))
}

// Info logs the provided arguments at [InfoLevel].
// Spaces are added between arguments when neither is a string.
func (s *SugaredLogger) Info(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.InfoLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Info(msg, zap.String("caller", c.File))
}

// Warn logs the provided arguments at [WarnLevel].
// Spaces are added between arguments when neither is a string.
func (s *SugaredLogger) Warn(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.WarnLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Warn(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Error(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.ErrorLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Error(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) DPanic(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DPanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.DPanic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Panic(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.PanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Panic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Fatal(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.FatalLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgs(args)
	s.getLogger().Logger.Fatal(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Debugf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DebugLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Debug(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Infof(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.InfoLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Info(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Warnf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.WarnLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Warn(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Errorf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.ErrorLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Error(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) DPanicf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DPanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.DPanic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Panicf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.PanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Panic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Fatalf(template string, args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.FatalLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getTemplateArgs(template, args)
	s.getLogger().Logger.Fatal(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DebugLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Debugw(msg, keysAndValues)
}

func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.InfoLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Infow(msg, keysAndValues)
}

func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.WarnLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Warnw(msg, keysAndValues)
}

func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.ErrorLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Errorw(msg, keysAndValues)
}

func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DPanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.DPanicw(msg, keysAndValues)
}

func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.PanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Panicw(msg, keysAndValues)
}

func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.FatalLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	keysAndValues = append(keysAndValues, zap.String("caller", c.File))
	s.SugaredLogger.Fatalw(msg, keysAndValues)
}

func (s *SugaredLogger) Debugln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DebugLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Debug(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Infoln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.InfoLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Info(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Warnln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.WarnLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Warn(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Errorln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.ErrorLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Error(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) DPanicln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.DPanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.DPanic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Panicln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.PanicLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Panic(msg, zap.String("caller", c.File))
}

func (s *SugaredLogger) Fatalln(args ...interface{}) {
	if !s.getLogger().getZapCore().Enabled(zap.FatalLevel) {
		return
	}
	c := CallerFrame(errors.GetPC())
	msg := getArgsLn(args)
	s.getLogger().Logger.Fatal(msg, zap.String("caller", c.File))
}
