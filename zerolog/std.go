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

package zerolog

import (
	"fmt"
	"os"
	"unsafe"

	"github.com/lxt1045/errors"
	"github.com/rs/zerolog"
)

type stdLogger Logger

func ToStd(logger *Logger) *stdLogger {
	return (*stdLogger)(unsafe.Pointer(logger))
}
func (logger *Logger) PointerToStd() *stdLogger {
	return (*stdLogger)(unsafe.Pointer(logger))
}
func (logger Logger) ToStd() stdLogger {
	return *(*stdLogger)(unsafe.Pointer(&logger))
}

func (l *stdLogger) Debug(args ...interface{}) {
	(*Logger)(l).Debug().print(errors.GetPC(), args...)
}
func (l *stdLogger) Debugf(format string, args ...interface{}) {
	(*Logger)(l).Debug().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Debugln(args ...interface{}) {
	(*Logger)(l).Debug().print(errors.GetPC(), args...)
}
func (l *stdLogger) Error(args ...interface{}) {
	(*Logger)(l).Error().print(errors.GetPC(), args...)
}
func (l *stdLogger) Errorf(format string, args ...interface{}) {
	(*Logger)(l).Error().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Errorln(args ...interface{}) {
	(*Logger)(l).Error().print(errors.GetPC(), args...)
}
func (l *stdLogger) Info(args ...interface{}) {
	(*Logger)(l).Info().print(errors.GetPC(), args...)
}
func (l *stdLogger) Infof(format string, args ...interface{}) {
	(*Logger)(l).Info().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Infoln(args ...interface{}) {
	(*Logger)(l).Info().print(errors.GetPC(), args...)
}
func (l *stdLogger) Warn(args ...interface{}) {
	(*Logger)(l).Warn().print(errors.GetPC(), args...)
}
func (l *stdLogger) Warnf(format string, args ...interface{}) {
	(*Logger)(l).Warn().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Warnln(args ...interface{}) {
	(*Logger)(l).Warn().print(errors.GetPC(), args...)
}

func (l *stdLogger) Fatal(args ...interface{}) {
	(*Logger)(l).Fatal().print(errors.GetPC(), args...)
}
func (l *stdLogger) Fatalf(format string, args ...interface{}) {
	(*Logger)(l).Fatal().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Fatalln(args ...interface{}) {
	(*Logger)(l).Fatal().print(errors.GetPC(), args...)
}

func (l *stdLogger) Panic(args ...interface{}) {
	(*Logger)(l).Panic().print(errors.GetPC(), args...)
}
func (l *stdLogger) Panicf(format string, args ...interface{}) {
	(*Logger)(l).Panic().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Panicln(args ...interface{}) {
	(*Logger)(l).Panic().print(errors.GetPC(), args...)
}

func (l *stdLogger) Print(args ...interface{}) {
	(*Logger)(l).Info().print(errors.GetPC(), args...)
}
func (l *stdLogger) Printf(format string, args ...interface{}) {
	(*Logger)(l).Info().printf(errors.GetPC(), format, args...)
}
func (l *stdLogger) Println(args ...interface{}) {
	(*Logger)(l).Info().print(errors.GetPC(), args...)
}

func (e *Event) print(pc errors.PC, args ...interface{}) {
	if (*zerolog.Event)(e).Enabled() {
		// e.CallerSkipFrame(1).Msg(fmt.Sprint(args...))
		// e = e.Timestamp().Str(
		e = e.Str(
			zerolog.CallerFieldName,
			pc.CallerFrame().FileLine,
		)
		e.Msg(fmt.Sprint(args...))
	}
}

func (e *Event) printf(pc errors.PC, format string, args ...interface{}) {
	if (*zerolog.Event)(e).Enabled() {
		// e.CallerSkipFrame(1).Msg(fmt.Sprint(args...))
		// e = e.Timestamp().Str(
		e = e.Str(
			zerolog.CallerFieldName,
			pc.CallerFrame().FileLine,
		)
		e.Msgf(format, args...)
	}
}

// Logger is the global logger.
var StdLogger = New(os.Stderr)

func Debug(args ...interface{}) {
	StdLogger.Debug().print(errors.GetPC(), args...)
}
func Debugf(format string, args ...interface{}) {
	StdLogger.Debug().printf(errors.GetPC(), format, args...)
}
func Debugln(args ...interface{}) {
	StdLogger.Debug().print(errors.GetPC(), args...)
}
func Error(args ...interface{}) {
	StdLogger.Error().print(errors.GetPC(), args...)
}
func Errorf(format string, args ...interface{}) {
	StdLogger.Error().printf(errors.GetPC(), format, args...)
}
func Errorln(args ...interface{}) {
	StdLogger.Error().print(errors.GetPC(), args...)
}
func Info(args ...interface{}) {
	StdLogger.Info().print(errors.GetPC(), args...)
}
func Infof(format string, args ...interface{}) {
	StdLogger.Info().printf(errors.GetPC(), format, args...)
}
func Infoln(args ...interface{}) {
	StdLogger.Info().print(errors.GetPC(), args...)
}

func Fatal(args ...interface{}) {
	StdLogger.Fatal().print(errors.GetPC(), args...)
}
func Fatalf(format string, args ...interface{}) {
	StdLogger.Fatal().printf(errors.GetPC(), format, args...)
}
func Fatalln(args ...interface{}) {
	StdLogger.Fatal().print(errors.GetPC(), args...)
}

func Panic(args ...interface{}) {
	StdLogger.Panic().print(errors.GetPC(), args...)
}
func Panicf(format string, args ...interface{}) {
	StdLogger.Panic().printf(errors.GetPC(), format, args...)
}
func Panicln(args ...interface{}) {
	StdLogger.Panic().print(errors.GetPC(), args...)
}

func Print(args ...interface{}) {
	StdLogger.Info().print(errors.GetPC(), args...)
}
func Printf(format string, args ...interface{}) {
	StdLogger.Info().printf(errors.GetPC(), format, args...)
}
func Println(args ...interface{}) {
	StdLogger.Info().print(errors.GetPC(), args...)
}
func Warn(args ...interface{}) {
	StdLogger.Warn().print(errors.GetPC(), args...)
}
func Warnf(format string, args ...interface{}) {
	StdLogger.Warn().printf(errors.GetPC(), format, args...)
}
func Warnln(args ...interface{}) {
	StdLogger.Warn().print(errors.GetPC(), args...)
}
