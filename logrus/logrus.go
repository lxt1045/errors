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

	"github.com/lxt1045/errors"
	"github.com/sirupsen/logrus"
)

type Fields = logrus.Fields
type Level = logrus.Level

type Entry struct {
	logrus.Entry
}

func toLogrusEntry(entry *Entry) *logrus.Entry {
	return &entry.Entry
}

func toEntry(entry *logrus.Entry) *Entry {
	return (*Entry)(unsafe.Pointer(entry))
}

func NewEntry(logger *logrus.Logger) *Entry {
	return toEntry(logrus.NewEntry(logger))
}

func (entry *Entry) AddCaller(pc errors.PC) *logrus.Entry {
	c := pc.CallerFrame()
	return toLogrusEntry(entry).WithFields(logrus.Fields{
		logrus.FieldKeyFunc: c.Func,
		logrus.FieldKeyFile: c.FileLine,
	})
}

func (entry *Entry) WithError(err error) *Entry {
	return toEntry(toLogrusEntry(entry).WithError(err))
}

func (entry *Entry) WithContext(ctx context.Context) *Entry {
	return toEntry(toLogrusEntry(entry).WithContext(ctx))
}

func (entry *Entry) WithField(key string, value interface{}) *Entry {
	return toEntry(toLogrusEntry(entry).WithField(key, value))
}

func (entry *Entry) WithFields(fields logrus.Fields) *Entry {
	return toEntry(toLogrusEntry(entry).WithFields(fields))
}

func (entry *Entry) WithTime(t time.Time) *Entry {
	return toEntry(toLogrusEntry(entry).WithTime(t))
}

//go:noinline
func (entry *Entry) Trace(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Trace(args...)
}

//go:noinline
func (entry *Entry) Debug(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Debug(args...)
}

//go:noinline
func (entry *Entry) Info(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Info(args...)
}

//go:noinline
func (entry *Entry) Print(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Print(args...)
}

//go:noinline
func (entry *Entry) Warn(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warn(args...)
}

//go:noinline
func (entry *Entry) Warning(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warning(args...)
}

//go:noinline
func (entry *Entry) Error(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Error(args...)
}

//go:noinline
func (entry *Entry) Fatal(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Fatal(args...)
}

//go:noinline
func (entry *Entry) Panic(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Panic(args...)
}

func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Tracef(format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Debugf(format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Infof(format, args...)
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Printf(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warnf(format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warningf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Errorf(format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Fatalf(format, args...)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Panicf(format, args...)
}

func (entry *Entry) Traceln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Traceln(args...)
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Debugln(args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Infoln(args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Println(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warnln(args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Warningln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Errorln(args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Fatalln(args...)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.AddCaller(errors.GetPC()).Panicln(args...)
}
