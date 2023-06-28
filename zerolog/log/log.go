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

package lxtlog

import (
	"context"
	"fmt"
	"io"
	"os"

	lxtlog "github.com/lxt1045/errors/zerolog"
	"github.com/rs/zerolog"
)

// Logger is the global logger.
var Logger = lxtlog.New(os.Stderr).With().Timestamp().Logger()

// Output duplicates the global logger and sets w as its output.
func Output(w io.Writer) lxtlog.Logger {
	return Logger.Output(w)
}

// With creates a child logger with the field added to its context.
func With() lxtlog.Context {
	return Logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func Level(level zerolog.Level) lxtlog.Logger {
	return Logger.Level(level)
}

// Sample returns a logger with the s sampler.
func Sample(s zerolog.Sampler) lxtlog.Logger {
	return Logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func Hook(h zerolog.Hook) lxtlog.Logger {
	return Logger.Hook(h)
}

// Err starts a new message with error level with err as a field if not nil or
// with info level if err is nil.
//
// You must call Msg on the returned event in order to send the event.
func Err(err error) *lxtlog.Event {
	return Logger.Err(err)
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func Trace() *lxtlog.Event {
	return Logger.Trace()
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func Debug() *lxtlog.Event {
	return Logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func Info() *lxtlog.Event {
	return Logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func Warn() *lxtlog.Event {
	return Logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func Error() *lxtlog.Event {
	return Logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func Fatal() *lxtlog.Event {
	return Logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func Panic() *lxtlog.Event {
	return Logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func WithLevel(level zerolog.Level) *lxtlog.Event {
	return Logger.WithLevel(level)
}

// Log starts a new message with no level. Setting lxtlog.GlobalLevel to
// lxtlog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func Log() *lxtlog.Event {
	return Logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func Print(v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msg(fmt.Sprint(v...))
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func Printf(format string, v ...interface{}) {
	Logger.Debug().CallerSkipFrame(1).Msgf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func Ctx(ctx context.Context) *lxtlog.Logger {
	return lxtlog.Ctx(ctx)
}
