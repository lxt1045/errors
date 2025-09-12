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
	"context"
	"fmt"
	"io"
	"net"
	"reflect"
	"time"
	"unsafe"

	"github.com/lxt1045/errors"
	"github.com/rs/zerolog"
)

// Level defines log levels.
type Level zerolog.Level

const (
	// DebugLevel defines debug log level.
	DebugLevel = Level(zerolog.DebugLevel)
	// InfoLevel defines info log level.
	InfoLevel = Level(zerolog.InfoLevel)
	// WarnLevel defines warn log level.
	WarnLevel = Level(zerolog.WarnLevel)
	// ErrorLevel defines error log level.
	ErrorLevel = Level(zerolog.ErrorLevel)
	// FatalLevel defines fatal log level.
	FatalLevel = Level(zerolog.FatalLevel)
	// PanicLevel defines panic log level.
	PanicLevel = Level(zerolog.PanicLevel)
	// NoLevel defines an absent log level.
	NoLevel = Level(zerolog.NoLevel)
	// Disabled disables the logger.
	Disabled = Level(zerolog.Disabled)

	// TraceLevel defines trace log level.
	TraceLevel = Level(zerolog.TraceLevel)
	// Values less than TraceLevel are handled as numbers.
)

var _ = func() bool {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	// zerolog.DefaultContextLogger = func() *zerolog.Logger {
	// 	return nil
	// }()
	return true
}()

type Logger zerolog.Logger
type Context zerolog.Context
type Event zerolog.Event

func SetGlobalLevel(l Level) {
	zerolog.SetGlobalLevel(zerolog.Level(l))
}

func toZeroLogger(logger *Logger) *zerolog.Logger {
	return (*zerolog.Logger)(unsafe.Pointer(logger))
}

func toLogger(logger *zerolog.Logger) *Logger {
	return (*Logger)(unsafe.Pointer(logger))
}

func toZeroContext(ctx *Context) *zerolog.Context {
	return (*zerolog.Context)(unsafe.Pointer(ctx))
}

func toContext(ctx *zerolog.Context) *Context {
	return (*Context)(unsafe.Pointer(ctx))
}
func loggerToContext(logger *zerolog.Logger) *Context {
	return (*Context)(unsafe.Pointer(logger))
}

func toZeroEvent(event *Event) *zerolog.Event {
	return (*zerolog.Event)(unsafe.Pointer(event))
}

func toEvent(event *zerolog.Event) *Event {
	return (*Event)(unsafe.Pointer(event))
}

func New(w io.Writer) Logger {
	return Logger(zerolog.New(w))
}

// Nop returns a disabled logger for which all operation are no-op.
func Nop() Logger {
	return Logger(zerolog.Nop())
}

func (l Logger) WithContext(ctx context.Context) context.Context {
	return zerolog.Logger(l).WithContext(ctx)
}

func Ctx(ctx context.Context) *Logger {
	return toLogger(zerolog.Ctx(ctx))
}

// Output duplicates the current logger and sets w as its output.
func (l Logger) Output(w io.Writer) Logger {
	return Logger(l.Output(w))
}

// With creates a child logger with the field added to its context.
func (l Logger) With() Context {
	return Context(zerolog.Logger(l).With())
}

// UpdateContext updates the internal logger's context.
//
// Use this method with caution. If unsure, prefer the With method.
func (l *Logger) UpdateContext(update func(c Context) Context) {
	(*zerolog.Logger)(l).UpdateContext(
		func(c zerolog.Context) zerolog.Context {
			return zerolog.Context(update(Context(c)))
		},
	)
}

// Level creates a child logger with the minimum accepted level set to level.
func (l Logger) Level(lvl zerolog.Level) Logger {
	l = Logger(zerolog.Logger(l).Level(lvl))
	return l
}

// GetLevel returns the current Level of l.
func (l Logger) GetLevel() Level {
	return Level(zerolog.Logger(l).GetLevel())
}

// Sample returns a logger with the s sampler.
func (l Logger) Sample(s zerolog.Sampler) Logger {
	l = Logger(zerolog.Logger(l).Sample(s))
	return l
}

// Hook returns a logger with the h Hook.
func (l Logger) Hook(h zerolog.Hook) Logger {
	l = Logger(zerolog.Logger(l).Hook(h))
	return l
}

// Trace starts a new message with trace level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Trace() *Event {
	return toEvent((*zerolog.Logger)(l).Trace().Timestamp())
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *Event {
	return toEvent((*zerolog.Logger)(l).Debug().Timestamp())
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *Event {
	return toEvent((*zerolog.Logger)(l).Info().Timestamp())
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *Event {
	return toEvent((*zerolog.Logger)(l).Warn().Timestamp())
}

func (l *Logger) Error() *Event {
	return toEvent((*zerolog.Logger)(l).Error().Timestamp())
}

func (l *Logger) Err(err error) *Event {
	return toEvent((*zerolog.Logger)(l).Err(err))
}

func (l *Logger) Fatal() *Event {
	return toEvent((*zerolog.Logger)(l).Fatal())
}

func (l *Logger) Panic() *Event {
	return toEvent((*zerolog.Logger)(l).Panic())
}

func (l *Logger) WithLevel(level zerolog.Level) *Event {
	return toEvent((*zerolog.Logger)(l).WithLevel(level))
}

func (l *Logger) Log() *Event {
	return toEvent((*zerolog.Logger)(l).Log())
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	(*zerolog.Logger)(l).Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	(*zerolog.Logger)(l).Printf(format, v...)
}

// Write implements the io.Writer interface. This is useful to set as a writer
// for the standard library log.
func (l Logger) Write(p []byte) (n int, err error) {
	return (zerolog.Logger)(l).Write(p)
}

// Logger returns the logger with the context previously set.
func (c Context) Logger() Logger {
	return Logger(zerolog.Context(c).Logger())
}

// Fields is a helper function to use a map or slice to set fields using type assertion.
// Only map[string]interface{} and []interface{} are accepted. []interface{} must
// alternate string keys and arbitrary values, and extraneous ones are ignored.
func (c Context) Fields(fields interface{}) Context {
	c = Context(zerolog.Context(c).Fields(fields))
	return c
}

// Dict adds the field key with the dict to the logger context.
func (c Context) Dict(key string, dict *zerolog.Event) Context {
	c = Context(zerolog.Context(c).Dict(key, dict))
	return c
}

// Array adds the field key with an array to the event context.
// Use zerolog.Arr() to create the array or pass a type that
// implement the LogArrayMarshaler interface.
func (c Context) Array(key string, arr zerolog.LogArrayMarshaler) Context {
	c = Context(zerolog.Context(c).Array(key, arr))
	return c
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (c Context) Object(key string, obj zerolog.LogObjectMarshaler) Context {
	c = Context(zerolog.Context(c).Object(key, obj))
	return c
}

// EmbedObject marshals and Embeds an object that implement the LogObjectMarshaler interface.
func (c Context) EmbedObject(obj zerolog.LogObjectMarshaler) Context {
	c = Context(zerolog.Context(c).EmbedObject(obj))
	return c
}

// Str adds the field key with val as a string to the logger context.
func (c Context) Str(key, val string) Context {
	c = Context(zerolog.Context(c).Str(key, val))
	return c
}

// Strs adds the field key with val as a string to the logger context.
func (c Context) Strs(key string, vals []string) Context {
	c = Context(zerolog.Context(c).Strs(key, vals))
	return c
}

// Stringer adds the field key with val.String() (or null if val is nil) to the logger context.
func (c Context) Stringer(key string, val fmt.Stringer) Context {
	c = Context(zerolog.Context(c).Stringer(key, val))
	return c
}

// Bytes adds the field key with val as a []byte to the logger context.
func (c Context) Bytes(key string, val []byte) Context {
	c = Context(zerolog.Context(c).Bytes(key, val))
	return c
}

// Hex adds the field key with val as a hex string to the logger context.
func (c Context) Hex(key string, val []byte) Context {
	c = Context(zerolog.Context(c).Hex(key, val))
	return c
}

// RawJSON adds already encoded JSON to context.
//
// No sanity check is performed on b; it must not contain carriage returns and
// be valid JSON.
func (c Context) RawJSON(key string, b []byte) Context {
	c = Context(zerolog.Context(c).RawJSON(key, b))
	return c
}

// AnErr adds the field key with serialized err to the logger context.
func (c Context) AnErr(key string, err error) Context {
	c = Context(zerolog.Context(c).AnErr(key, err))
	return c
}

// Errs adds the field key with errs as an array of serialized errors to the
// logger context.
func (c Context) Errs(key string, errs []error) Context {
	c = Context(zerolog.Context(c).Errs(key, errs))
	return c
}

// Err adds the field "error" with serialized err to the logger context.
func (c Context) Err(err error) Context {
	c = Context(zerolog.Context(c).Err(err))
	return c
}

// Bool adds the field key with val as a bool to the logger context.
func (c Context) Bool(key string, b bool) Context {
	c = Context(zerolog.Context(c).Bool(key, b))
	return c
}

// Bools adds the field key with val as a []bool to the logger context.
func (c Context) Bools(key string, b []bool) Context {
	c = Context(zerolog.Context(c).Bools(key, b))
	return c
}

// Int adds the field key with i as a int to the logger context.
func (c Context) Int(key string, i int) Context {
	c = Context(zerolog.Context(c).Int(key, i))
	return c
}

// Ints adds the field key with i as a []int to the logger context.
func (c Context) Ints(key string, i []int) Context {
	c = Context(zerolog.Context(c).Ints(key, i))
	return c
}

// Int8 adds the field key with i as a int8 to the logger context.
func (c Context) Int8(key string, i int8) Context {
	c = Context(zerolog.Context(c).Int8(key, i))
	return c
}

// Ints8 adds the field key with i as a []int8 to the logger context.
func (c Context) Ints8(key string, i []int8) Context {
	c = Context(zerolog.Context(c).Ints8(key, i))
	return c
}

// Int16 adds the field key with i as a int16 to the logger context.
func (c Context) Int16(key string, i int16) Context {
	c = Context(zerolog.Context(c).Int16(key, i))
	return c
}

// Ints16 adds the field key with i as a []int16 to the logger context.
func (c Context) Ints16(key string, i []int16) Context {
	c = Context(zerolog.Context(c).Ints16(key, i))
	return c
}

// Int32 adds the field key with i as a int32 to the logger context.
func (c Context) Int32(key string, i int32) Context {
	c = Context(zerolog.Context(c).Int32(key, i))
	return c
}

// Ints32 adds the field key with i as a []int32 to the logger context.
func (c Context) Ints32(key string, i []int32) Context {
	c = Context(zerolog.Context(c).Ints32(key, i))
	return c
}

// Int64 adds the field key with i as a int64 to the logger context.
func (c Context) Int64(key string, i int64) Context {
	c = Context(zerolog.Context(c).Int64(key, i))
	return c
}

// Ints64 adds the field key with i as a []int64 to the logger context.
func (c Context) Ints64(key string, i []int64) Context {
	c = Context(zerolog.Context(c).Ints64(key, i))
	return c
}

// Uint adds the field key with i as a uint to the logger context.
func (c Context) Uint(key string, i uint) Context {
	c = Context(zerolog.Context(c).Uint(key, i))
	return c
}

// Uints adds the field key with i as a []uint to the logger context.
func (c Context) Uints(key string, i []uint) Context {
	c = Context(zerolog.Context(c).Uints(key, i))
	return c
}

// Uint8 adds the field key with i as a uint8 to the logger context.
func (c Context) Uint8(key string, i uint8) Context {
	c = Context(zerolog.Context(c).Uint8(key, i))
	return c
}

// Uints8 adds the field key with i as a []uint8 to the logger context.
func (c Context) Uints8(key string, i []uint8) Context {
	c = Context(zerolog.Context(c).Uints8(key, i))
	return c
}

// Uint16 adds the field key with i as a uint16 to the logger context.
func (c Context) Uint16(key string, i uint16) Context {
	c = Context(zerolog.Context(c).Uint16(key, i))
	return c
}

// Uints16 adds the field key with i as a []uint16 to the logger context.
func (c Context) Uints16(key string, i []uint16) Context {
	c = Context(zerolog.Context(c).Uints16(key, i))
	return c
}

// Uint32 adds the field key with i as a uint32 to the logger context.
func (c Context) Uint32(key string, i uint32) Context {
	c = Context(zerolog.Context(c).Uint32(key, i))
	return c
}

// Uints32 adds the field key with i as a []uint32 to the logger context.
func (c Context) Uints32(key string, i []uint32) Context {
	c = Context(zerolog.Context(c).Uints32(key, i))
	return c
}

// Uint64 adds the field key with i as a uint64 to the logger context.
func (c Context) Uint64(key string, i uint64) Context {
	c = Context(zerolog.Context(c).Uint64(key, i))
	return c
}

// Uints64 adds the field key with i as a []uint64 to the logger context.
func (c Context) Uints64(key string, i []uint64) Context {
	c = Context(zerolog.Context(c).Uints64(key, i))
	return c
}

// Float32 adds the field key with f as a float32 to the logger context.
func (c Context) Float32(key string, f float32) Context {
	c = Context(zerolog.Context(c).Float32(key, f))
	return c
}

// Floats32 adds the field key with f as a []float32 to the logger context.
func (c Context) Floats32(key string, f []float32) Context {
	c = Context(zerolog.Context(c).Floats32(key, f))
	return c
}

// Float64 adds the field key with f as a float64 to the logger context.
func (c Context) Float64(key string, f float64) Context {
	c = Context(zerolog.Context(c).Float64(key, f))
	return c
}

// Floats64 adds the field key with f as a []float64 to the logger context.
func (c Context) Floats64(key string, f []float64) Context {
	c = Context(zerolog.Context(c).Floats64(key, f))
	return c
}

func (c Context) Timestamp() Context {
	c = Context(zerolog.Context(c).Timestamp())
	return c
}

// Time adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (c Context) Time(key string, t time.Time) Context {
	c = Context(zerolog.Context(c).Time(key, t))
	return c
}

// Times adds the field key with t formated as string using zerolog.TimeFieldFormat.
func (c Context) Times(key string, t []time.Time) Context {
	c = Context(zerolog.Context(c).Times(key, t))
	return c
}

// Dur adds the fields key with d divided by unit and stored as a float.
func (c Context) Dur(key string, d time.Duration) Context {
	c = Context(zerolog.Context(c).Dur(key, d))
	return c
}

// Durs adds the fields key with d divided by unit and stored as a float.
func (c Context) Durs(key string, d []time.Duration) Context {
	c = Context(zerolog.Context(c).Durs(key, d))
	return c
}

func (c Context) Interface(key string, i interface{}) Context {
	c = Context(zerolog.Context(c).Interface(key, i))
	return c
}

type callerHook struct {
	callerSkipFrameCount int
}

func (ch callerHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if ch.callerSkipFrameCount < 1 {
		ch.callerSkipFrameCount = 2
	}
	cs := errors.CallersSkip(ch.callerSkipFrameCount - 1)
	e.Str(
		zerolog.CallerFieldName,
		cs[0].FileLine, //zerolog.CallerMarshalFunc(0, cs[0].File, cs[0].Line),
	)
}

func (c Context) Caller() Context {
	logger := c.Logger().Hook(&callerHook{2 + zerolog.CallerSkipFrameCount})
	c = *loggerToContext((*zerolog.Logger)(&logger))
	return c
}

func (c Context) CallerWithSkipFrameCount(skipFrameCount int) Context {
	logger := c.Logger().Hook(&callerHook{skipFrameCount + zerolog.CallerSkipFrameCount})
	c = *loggerToContext((*zerolog.Logger)(&logger))
	return c
}

// Stack enables stack trace printing for the error passed to Err().
func (c Context) Stack() Context {
	c = Context(zerolog.Context(c).Stack())
	return c
}

// IPAddr adds IPv4 or IPv6 Address to the context
func (c Context) IPAddr(key string, ip net.IP) Context {
	c = Context(zerolog.Context(c).IPAddr(key, ip))
	return c
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the context
func (c Context) IPPrefix(key string, pfx net.IPNet) Context {
	c = Context(zerolog.Context(c).IPPrefix(key, pfx))
	return c
}

// MACAddr adds MAC address to the context
func (c Context) MACAddr(key string, ha net.HardwareAddr) Context {
	c = Context(zerolog.Context(c).MACAddr(key, ha))
	return c
}

// func (e *Event) Enabled() bool {
// 	return e != nil && e.level != Disabled
// }

// Discard disables the event so Msg(f) won't print it.
func (e *Event) Discard() *Event {
	return toEvent(toZeroEvent(e).Discard())
}

// func (e *Event) Msg(msg string) {
// 	e.Event.Msg(msg)
// }

// func (e *Event) Send() {
// 	e.Event.Send()
// }

// func (e *Event) Msgf(format string, v ...interface{}) {
// 	e.Event.Msgf(format, v...)
// }

// func (e *Event) MsgFunc(createMsg func() string) {
// 	e.Event.MsgFunc(createMsg)
// }

func (e *Event) Fields(fields interface{}) *Event {
	return toEvent(toZeroEvent(e).Fields(fields))
}

// Dict adds the field key with a dict to the event context.
// Use zerolog.Dict() to create the dictionary.
func (e *Event) Dict(key string, dict *Event) *Event {
	return toEvent(toZeroEvent(e).Dict(key, toZeroEvent(dict)))
}

func Dict() *Event {
	return toEvent(zerolog.Dict())
}

func (e *Event) Array(key string, arr zerolog.LogArrayMarshaler) *Event {
	return toEvent(toZeroEvent(e).Array(key, arr))
}

// Object marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) Object(key string, obj zerolog.LogObjectMarshaler) *Event {
	return toEvent(toZeroEvent(e).Object(key, obj))
}

// Func allows an anonymous func to run only if the event is enabled.
func (e *Event) Func(f func(e *Event)) *Event {
	return toEvent(
		toZeroEvent(e).Func(
			func(e *zerolog.Event) {
				f(toEvent(e))
			},
		),
	)
}

// EmbedObject marshals an object that implement the LogObjectMarshaler interface.
func (e *Event) EmbedObject(obj zerolog.LogObjectMarshaler) *Event {
	return toEvent(toZeroEvent(e).EmbedObject(obj))
}

// Str adds the field key with val as a string to the *Event context.
func (e *Event) Str(key, val string) *Event {
	return toEvent(toZeroEvent(e).Str(key, val))
}

// Strs adds the field key with vals as a []string to the *Event context.
func (e *Event) Strs(key string, vals []string) *Event {
	return toEvent(toZeroEvent(e).Strs(key, vals))
}

// Stringer adds the field key with val.String() (or null if val is nil)
// to the *Event context.
func (e *Event) Stringer(key string, val fmt.Stringer) *Event {
	return toEvent(toZeroEvent(e).Stringer(key, val))
}

func (e *Event) Stringers(key string, vals []fmt.Stringer) *Event {
	return toEvent(toZeroEvent(e).Stringers(key, vals))
}

func (e *Event) Bytes(key string, val []byte) *Event {
	return toEvent(toZeroEvent(e).Bytes(key, val))
}

// Hex adds the field key with val as a hex string to the *Event context.
func (e *Event) Hex(key string, val []byte) *Event {
	return toEvent(toZeroEvent(e).Hex(key, val))
}

func (e *Event) RawJSON(key string, b []byte) *Event {
	return toEvent(toZeroEvent(e).RawJSON(key, b))
}

// RawCBOR adds already encoded CBOR to the log line under key.
//
// No sanity check is performed on b
// Note: The full featureset of CBOR is supported as data will not be mapped to json but stored as data-url
func (e *Event) RawCBOR(key string, b []byte) *Event {
	return toEvent(toZeroEvent(e).RawCBOR(key, b))
}

// AnErr adds the field key with serialized err to the *Event context.
// If err is nil, no field is added.
func (e *Event) AnErr(key string, err error) *Event {
	return toEvent(toZeroEvent(e).AnErr(key, err))
}

// Errs adds the field key with errs as an array of serialized errors to the
// *Event context.
func (e *Event) Errs(key string, errs []error) *Event {
	return toEvent(toZeroEvent(e).Errs(key, errs))
}

func (e *Event) Err(err error) *Event {
	return toEvent(toZeroEvent(e).Err(err))
}

// Stack enables stack trace printing for the error passed to Err().
//
// ErrorStackMarshaler must be set for this method to do something.
func (e *Event) Stack() *Event {
	return toEvent(toZeroEvent(e).Stack())
}

// Bool adds the field key with val as a bool to the *Event context.
func (e *Event) Bool(key string, b bool) *Event {
	return toEvent(toZeroEvent(e).Bool(key, b))
}

// Bools adds the field key with val as a []bool to the *Event context.
func (e *Event) Bools(key string, b []bool) *Event {
	return toEvent(toZeroEvent(e).Bools(key, b))
}

// Int adds the field key with i as a int to the *Event context.
func (e *Event) Int(key string, i int) *Event {
	return toEvent(toZeroEvent(e).Int(key, i))
}

// Ints adds the field key with i as a []int to the *Event context.
func (e *Event) Ints(key string, i []int) *Event {
	return toEvent(toZeroEvent(e).Ints(key, i))
}

// Int8 adds the field key with i as a int8 to the *Event context.
func (e *Event) Int8(key string, i int8) *Event {
	return toEvent(toZeroEvent(e).Int8(key, i))
}

// Ints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Ints8(key string, i []int8) *Event {
	return toEvent(toZeroEvent(e).Ints8(key, i))
}

// Int16 adds the field key with i as a int16 to the *Event context.
func (e *Event) Int16(key string, i int16) *Event {
	return toEvent(toZeroEvent(e).Int16(key, i))
}

// Ints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Ints16(key string, i []int16) *Event {
	return toEvent(toZeroEvent(e).Ints16(key, i))
}

// Int32 adds the field key with i as a int32 to the *Event context.
func (e *Event) Int32(key string, i int32) *Event {
	return toEvent(toZeroEvent(e).Int32(key, i))
}

// Ints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Ints32(key string, i []int32) *Event {
	return toEvent(toZeroEvent(e).Ints32(key, i))
}

// Int64 adds the field key with i as a int64 to the *Event context.
func (e *Event) Int64(key string, i int64) *Event {
	return toEvent(toZeroEvent(e).Int64(key, i))
}

// Ints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Ints64(key string, i []int64) *Event {
	return toEvent(toZeroEvent(e).Ints64(key, i))
}

// Uint adds the field key with i as a uint to the *Event context.
func (e *Event) Uint(key string, i uint) *Event {
	return toEvent(toZeroEvent(e).Uint(key, i))
}

// Uints adds the field key with i as a []int to the *Event context.
func (e *Event) Uints(key string, i []uint) *Event {
	return toEvent(toZeroEvent(e).Uints(key, i))
}

// Uint8 adds the field key with i as a uint8 to the *Event context.
func (e *Event) Uint8(key string, i uint8) *Event {
	return toEvent(toZeroEvent(e).Uint8(key, i))
}

// Uints8 adds the field key with i as a []int8 to the *Event context.
func (e *Event) Uints8(key string, i []uint8) *Event {
	return toEvent(toZeroEvent(e).Uints8(key, i))
}

// Uint16 adds the field key with i as a uint16 to the *Event context.
func (e *Event) Uint16(key string, i uint16) *Event {
	return toEvent(toZeroEvent(e).Uint16(key, i))
}

// Uints16 adds the field key with i as a []int16 to the *Event context.
func (e *Event) Uints16(key string, i []uint16) *Event {
	return toEvent(toZeroEvent(e).Uints16(key, i))
}

// Uint32 adds the field key with i as a uint32 to the *Event context.
func (e *Event) Uint32(key string, i uint32) *Event {
	return toEvent(toZeroEvent(e).Uint32(key, i))
}

// Uints32 adds the field key with i as a []int32 to the *Event context.
func (e *Event) Uints32(key string, i []uint32) *Event {
	return toEvent(toZeroEvent(e).Uints32(key, i))
}

// Uint64 adds the field key with i as a uint64 to the *Event context.
func (e *Event) Uint64(key string, i uint64) *Event {
	return toEvent(toZeroEvent(e).Uint64(key, i))
}

// Uints64 adds the field key with i as a []int64 to the *Event context.
func (e *Event) Uints64(key string, i []uint64) *Event {
	return toEvent(toZeroEvent(e).Uints64(key, i))
}

// Float32 adds the field key with f as a float32 to the *Event context.
func (e *Event) Float32(key string, f float32) *Event {
	return toEvent(toZeroEvent(e).Float32(key, f))
}

// Floats32 adds the field key with f as a []float32 to the *Event context.
func (e *Event) Floats32(key string, f []float32) *Event {
	return toEvent(toZeroEvent(e).Floats32(key, f))
}

// Float64 adds the field key with f as a float64 to the *Event context.
func (e *Event) Float64(key string, f float64) *Event {
	return toEvent(toZeroEvent(e).Float64(key, f))
}

// Floats64 adds the field key with f as a []float64 to the *Event context.
func (e *Event) Floats64(key string, f []float64) *Event {
	return toEvent(toZeroEvent(e).Floats64(key, f))
}

// Timestamp adds the current local time as UNIX timestamp to the *Event context with the "time" key.
// To customize the key name, change zerolog.TimestampFieldName.
//
// NOTE: It won't dedupe the "time" key if the *Event (or *Context) has one
// already.
func (e *Event) Timestamp() *Event {
	return toEvent(toZeroEvent(e).Timestamp())
}

// Time adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (e *Event) Time(key string, t time.Time) *Event {
	return toEvent(toZeroEvent(e).Time(key, t))
}

// Times adds the field key with t formatted as string using zerolog.TimeFieldFormat.
func (e *Event) Times(key string, t []time.Time) *Event {
	return toEvent(toZeroEvent(e).Times(key, t))
}

// Dur adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Dur(key string, d time.Duration) *Event {
	return toEvent(toZeroEvent(e).Dur(key, d))
}

// Durs adds the field key with duration d stored as zerolog.DurationFieldUnit.
// If zerolog.DurationFieldInteger is true, durations are rendered as integer
// instead of float.
func (e *Event) Durs(key string, d []time.Duration) *Event {
	return toEvent(toZeroEvent(e).Durs(key, d))
}

// TimeDiff adds the field key with positive duration between time t and start.
// If time t is not greater than start, duration will be 0.
// Duration format follows the same principle as Dur().
func (e *Event) TimeDiff(key string, t time.Time, start time.Time) *Event {
	return toEvent(toZeroEvent(e).TimeDiff(key, t, start))
}

// Any is a wrapper around Event.Interface.
func (e *Event) Any(key string, i interface{}) *Event {
	return toEvent(toZeroEvent(e).Any(key, i))
}

// Interface adds the field key with i marshaled using reflection.
func (e *Event) Interface(key string, i interface{}) *Event {
	return toEvent(toZeroEvent(e).Interface(key, i))
}

// Type adds the field key with val's type using reflection.
func (e *Event) Type(key string, val interface{}) *Event {
	return toEvent(toZeroEvent(e).Type(key, val))
}

func (e *Event) CallerSkipFrame(skip int) *Event {
	return toEvent(toZeroEvent(e).CallerSkipFrame(skip))
}

// 使用heck的方式获取 zerolog.Event.skipFrame
var zeroEventSkipFrameOffset = func() uintptr {
	typ := reflect.TypeOf(zerolog.Event{})
	field, ok := typ.FieldByName("skipFrame")
	if !ok {
		panic("zerolog.Event.skipFrame not exist")
	}
	return field.Offset
}()

func getEventSkipFrame(e *Event) int {
	if e == nil {
		return zerolog.CallerSkipFrameCount
	}
	return *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(e)) + zeroEventSkipFrameOffset))
}

func (e *Event) Caller(skip ...int) *Event {
	if e == nil {
		return e
	}
	levelSkip := zerolog.CallerSkipFrameCount + getEventSkipFrame(e)
	if len(skip) > 0 && skip[0] > 0 {
		levelSkip += skip[0]
	}
	if levelSkip <= 2 {
		c := errors.GetPC().CallerFrame()
		e = e.Str(
			zerolog.CallerFieldName,
			c.FileLine, // zerolog.CallerMarshalFunc(0, c.File, c.Line),
		)
		return e
	}
	if levelSkip-2 <= errors.DefaultDepth {
		cs := errors.CallersSkip(levelSkip - 2)
		e = e.Str(
			zerolog.CallerFieldName,
			cs[0].FileLine, //zerolog.CallerMarshalFunc(0, cs[0].File, cs[0].Line),
		)
		return e
	}
	return toEvent(toZeroEvent(e).Caller(skip...))
}

// IPAddr adds IPv4 or IPv6 Address to the event
func (e *Event) IPAddr(key string, ip net.IP) *Event {
	return toEvent(toZeroEvent(e).IPAddr(key, ip))
}

// IPPrefix adds IPv4 or IPv6 Prefix (address and mask) to the event
func (e *Event) IPPrefix(key string, pfx net.IPNet) *Event {
	return toEvent(toZeroEvent(e).IPPrefix(key, pfx))
}

// MACAddr adds MAC address to the event
func (e *Event) MACAddr(key string, ha net.HardwareAddr) *Event {
	return toEvent(toZeroEvent(e).MACAddr(key, ha))
}

// MACAddr adds MAC address to the event
func (e *Event) Send() {
	toZeroEvent(e).Send()
}

func (e *Event) Msg(msg string) {
	toZeroEvent(e).Msg(msg)
}

func (e *Event) Msgf(format string, v ...interface{}) {
	toZeroEvent(e).Msgf(format, v...)
}
