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

package logrus

import (
	"context"
	"time"
	"unsafe"
	_ "unsafe" //nolint:bgolint

	"github.com/lxt1045/errors"
	"github.com/sirupsen/logrus"
)

var (
	StandardLogger  = logrus.StandardLogger
	SetOutput       = logrus.SetOutput
	SetFormatter    = logrus.SetFormatter
	SetReportCaller = logrus.SetReportCaller
	SetLevel        = logrus.SetLevel
	GetLevel        = logrus.GetLevel
	IsLevelEnabled  = logrus.IsLevelEnabled
	AddHook         = logrus.AddHook
)

// WithError creates an entry from the standard logger and adds an error to it, using the value defined in ErrorKey as key.
func WithError(err error) *Entry {
	return toEntry(logrus.WithError(err))
}

func WithContext(ctx context.Context) *Entry {
	return toEntry(logrus.WithContext(ctx))
}

func WithField(key string, value interface{}) *Entry {
	return toEntry(logrus.WithField(key, value))
}

func WithFields(fields logrus.Fields) *Entry {
	return toEntry(logrus.WithFields(fields))
}

func WithTime(t time.Time) *Entry {
	return toEntry(logrus.WithTime(t))
}

type Logger struct {
	logrus.Logger
}

func toLogrusLogger(logger *Logger) *logrus.Logger {
	return &logger.Logger
}

func toLogger(logger *logrus.Logger) *Logger {
	return (*Logger)(unsafe.Pointer(logger))
}

func New() *Logger {
	return toLogger(logrus.New())
}
func (logger *Logger) AddCaller(pc uintptr) *Entry {
	c := CallerFrame(pc)
	return logger.WithFields(logrus.Fields{
		logrus.FieldKeyFunc: c.Func,
		logrus.FieldKeyFile: c.File,
	})
}

func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return toEntry(toLogrusLogger(logger).WithField(key, value))
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields logrus.Fields) *Entry {
	return toEntry(toLogrusLogger(logger).WithFields(fields))
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	return toEntry(toLogrusLogger(logger).WithError(err))
}

// Add a context to the log entry.
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	return toEntry(toLogrusLogger(logger).WithContext(ctx))
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	return toEntry(toLogrusLogger(logger).WithTime(t))
}

func (logger *Logger) Logf(level logrus.Level, format string, args ...interface{}) {
	logger.logf(level, errors.GetPC(), format, args...)
}

func (logger *Logger) logf(level logrus.Level, pc uintptr, format string, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.AddCaller(pc)
		entry.Logf(level, format, args...)
	}
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.logf(logrus.TraceLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.logf(logrus.DebugLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.logf(logrus.InfoLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	entry := logger.AddCaller(errors.GetPC())
	entry.Printf(format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.logf(logrus.WarnLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.logf(logrus.WarnLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.logf(logrus.ErrorLevel, errors.GetPC(), format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.logf(logrus.FatalLevel, errors.GetPC(), format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.logf(logrus.PanicLevel, errors.GetPC(), format, args...)
}

// Log will log a message at the level given as parameter.
// Warning: using Log at Panic or Fatal level will not respectively Panic nor Exit.
// For this behaviour Logger.Panic or Logger.Fatal should be used instead.
func (logger *Logger) Log(level logrus.Level, args ...interface{}) {
	logger.log(level, errors.GetPC(), args...)
}
func (logger *Logger) log(level logrus.Level, pc uintptr, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.AddCaller(pc)
		entry.Log(level, args...)
	}
}

func (logger *Logger) LogFn(level logrus.Level, fn logrus.LogFunction) {
	logger.log(level, errors.GetPC(), fn()...)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.log(logrus.TraceLevel, errors.GetPC(), args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.log(logrus.DebugLevel, errors.GetPC(), args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.log(logrus.InfoLevel, errors.GetPC(), args...)
}

func (logger *Logger) Print(args ...interface{}) {
	entry := logger.AddCaller(errors.GetPC())
	entry.Print(args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.log(logrus.WarnLevel, errors.GetPC(), args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.log(logrus.WarnLevel, errors.GetPC(), args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.log(logrus.ErrorLevel, errors.GetPC(), args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.log(logrus.FatalLevel, errors.GetPC(), args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.log(logrus.PanicLevel, errors.GetPC(), args...)
}

func (logger *Logger) TraceFn(fn logrus.LogFunction) {
	logger.log(logrus.TraceLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) DebugFn(fn logrus.LogFunction) {
	logger.log(logrus.DebugLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) InfoFn(fn logrus.LogFunction) {
	logger.log(logrus.InfoLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) PrintFn(fn logrus.LogFunction) {
	entry := logger.AddCaller(errors.GetPC())
	entry.Print(fn()...)
}

func (logger *Logger) WarnFn(fn logrus.LogFunction) {
	logger.log(logrus.WarnLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) WarningFn(fn logrus.LogFunction) {
	logger.log(logrus.WarnLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) ErrorFn(fn logrus.LogFunction) {
	logger.log(logrus.ErrorLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) FatalFn(fn logrus.LogFunction) {
	logger.log(logrus.FatalLevel, errors.GetPC(), fn()...)
	logger.Exit(1)
}

func (logger *Logger) PanicFn(fn logrus.LogFunction) {
	logger.log(logrus.PanicLevel, errors.GetPC(), fn()...)
}

func (logger *Logger) Logln(level logrus.Level, args ...interface{}) {
	logger.logln(level, errors.GetPC(), args...)
}
func (logger *Logger) logln(level logrus.Level, pc uintptr, args ...interface{}) {
	if logger.IsLevelEnabled(level) {
		entry := logger.AddCaller(pc)
		entry.Logln(level, args...)
	}
}

func (logger *Logger) Traceln(args ...interface{}) {
	logger.logln(logrus.TraceLevel, errors.GetPC(), args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.logln(logrus.DebugLevel, errors.GetPC(), args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.logln(logrus.InfoLevel, errors.GetPC(), args...)
}

func (logger *Logger) Println(args ...interface{}) {
	entry := logger.AddCaller(errors.GetPC())
	entry.Println(args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.logln(logrus.WarnLevel, errors.GetPC(), args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.logln(logrus.WarnLevel, errors.GetPC(), args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.logln(logrus.ErrorLevel, errors.GetPC(), args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.logln(logrus.FatalLevel, errors.GetPC(), args...)
	logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.logln(logrus.PanicLevel, errors.GetPC(), args...)
}

func (logger *Logger) Exit(code int) {
	toLogrusLogger(logger).Exit(code)
}

//

//

func Trace(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.TraceLevel, errors.GetPC(), args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.DebugLevel, errors.GetPC(), args...)
}

// Print logs a message at level Info on the standard logger.
func Print(args ...interface{}) {
	entry := toLogger(logrus.StandardLogger()).AddCaller(errors.GetPC())
	entry.Print(args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.InfoLevel, errors.GetPC(), args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.WarnLevel, errors.GetPC(), args...)
}

// Warning logs a message at level Warn on the standard logger.
func Warning(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.WarnLevel, errors.GetPC(), args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.ErrorLevel, errors.GetPC(), args...)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	toLogger(logrus.StandardLogger()).log(logrus.PanicLevel, errors.GetPC(), args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	logger := toLogger(logrus.StandardLogger())
	logger.log(logrus.FatalLevel, errors.GetPC(), args...)
	logger.Exit(1)
}

// TraceFn logs a message from a func at level Trace on the standard logger.
func TraceFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.TraceLevel, errors.GetPC(), fn()...)
}

// DebugFn logs a message from a func at level Debug on the standard logger.
func DebugFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.DebugLevel, errors.GetPC(), fn()...)
}

// PrintFn logs a message from a func at level Info on the standard logger.
func PrintFn(fn logrus.LogFunction) {
	logger := toLogger(logrus.StandardLogger())
	entry := logger.AddCaller(errors.GetPC())
	entry.Print(fn()...)
}

// InfoFn logs a message from a func at level Info on the standard logger.
func InfoFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.InfoLevel, errors.GetPC(), fn()...)
}

// WarnFn logs a message from a func at level Warn on the standard logger.
func WarnFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.WarnLevel, errors.GetPC(), fn()...)
}

// WarningFn logs a message from a func at level Warn on the standard logger.
func WarningFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.WarnLevel, errors.GetPC(), fn()...)
}

// ErrorFn logs a message from a func at level Error on the standard logger.
func ErrorFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.ErrorLevel, errors.GetPC(), fn()...)
}

// PanicFn logs a message from a func at level Panic on the standard logger.
func PanicFn(fn logrus.LogFunction) {
	toLogger(logrus.StandardLogger()).log(logrus.PanicLevel, errors.GetPC(), fn()...)
}

// FatalFn logs a message from a func at level Fatal on the standard logger then the process will exit with status set to 1.
func FatalFn(fn logrus.LogFunction) {
	logger := toLogger(logrus.StandardLogger())
	logger.log(logrus.FatalLevel, errors.GetPC(), fn()...)
	logger.Exit(1)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.TraceLevel, errors.GetPC(), format, args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.DebugLevel, errors.GetPC(), format, args...)
}

// Printf logs a message at level Info on the standard logger.
func Printf(format string, args ...interface{}) {
	entry := toLogger(logrus.StandardLogger()).AddCaller(errors.GetPC())
	entry.Printf(format, args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.InfoLevel, errors.GetPC(), format, args...)
}

// Warnf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.WarnLevel, errors.GetPC(), format, args...)
}

// Warningf logs a message at level Warn on the standard logger.
func Warningf(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.WarnLevel, errors.GetPC(), format, args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.ErrorLevel, errors.GetPC(), format, args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	toLogger(logrus.StandardLogger()).logf(logrus.PanicLevel, errors.GetPC(), format, args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	logger := toLogger(logrus.StandardLogger())
	logger.logf(logrus.FatalLevel, errors.GetPC(), format, args...)
	logger.Exit(1)
}

// Traceln logs a message at level Trace on the standard logger.
func Traceln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.TraceLevel, errors.GetPC(), args...)
}

// Debugln logs a message at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.DebugLevel, errors.GetPC(), args...)
}

// Println logs a message at level Info on the standard logger.
func Println(args ...interface{}) {
	entry := toLogger(logrus.StandardLogger()).AddCaller(errors.GetPC())
	entry.Println(args...)
}

// Infoln logs a message at level Info on the standard logger.
func Infoln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.InfoLevel, errors.GetPC(), args...)
}

// Warnln logs a message at level Warn on the standard logger.
func Warnln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.WarnLevel, errors.GetPC(), args...)
}

// Warningln logs a message at level Warn on the standard logger.
func Warningln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.WarnLevel, errors.GetPC(), args...)
}

// Errorln logs a message at level Error on the standard logger.
func Errorln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.ErrorLevel, errors.GetPC(), args...)
}

// Panicln logs a message at level Panic on the standard logger.
func Panicln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.PanicLevel, errors.GetPC(), args...)
}

// Fatalln logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalln(args ...interface{}) {
	toLogger(logrus.StandardLogger()).logln(logrus.FatalLevel, errors.GetPC(), args...)
	toLogger(logrus.StandardLogger()).Exit(1)
}
