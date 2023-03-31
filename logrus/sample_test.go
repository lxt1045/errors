package logrus

import (
	"bytes"
	"context"
	"runtime"
	"strconv"
	"testing"

	"github.com/lxt1045/errors"
	"github.com/sirupsen/logrus"
)

// sample 1
func LineByRuntime() string {
	_, file, n, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return file + ":" + strconv.Itoa(n)
}
func LineByLog() string {
	pc := errors.GetPC()
	cf := CallerFrame(pc)
	return cf.File // + ":" + strconv.Itoa(cf.Line)
}

func TestSample(t *testing.T) {
	t.Run("LineByRuntime", func(t *testing.T) {
		t.Log(runtime.Caller(0))
	})
	t.Run("Line", func(t *testing.T) {
		t.Log("line:", LineByLog())
	})
}

var gLine string

func BenchmarkUnmarshalInterface(b *testing.B) {
	b.Run("LineByRuntime", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, gLine, _, _ = runtime.Caller(0)
		}
	})

	b.Run("Log", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			gLine = LineByLog()
		}
	})

}

func TestLog(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		w := &bytes.Buffer{}
		logrus.SetReportCaller(true)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(w)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		// h := &Hook{AppName: "awesome-web"}
		// logrus.AddHook(h)
		logrus.Info("info msg")

		t.Log(w.String())
	})

	t.Run("2", func(t *testing.T) {
		w := &bytes.Buffer{}
		logrus.SetReportCaller(false)
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetOutput(w)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		// h := &Hook{AppName: "awesome-web"}
		// logrus.AddHook(h)

		w.Reset()
		WithContext(context.TODO()).Info("info msg")
		t.Log(w.String())
		w.Reset()
		WithContext(context.TODO()).Debug("info msg")
		WithContext(context.TODO()).Info("info msg")
		WithContext(context.TODO()).Warnf("info msg")
		t.Log(w.String())
	})
}

func BenchmarkLog(b *testing.B) {
	bs := make([]byte, 1<<20)
	w := bytes.NewBuffer(bs)
	logrus.SetReportCaller(true)
	logrus.SetOutput(w)
	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})
	// h := &Hook{AppName: "awesome-web"}
	// logrus.AddHook(h)
	logrus.Info("info msg")
	// b.Log(w.String())

	ctx := context.TODO()

	b.Run("logrus+caller", func(b *testing.B) {
		logrus.SetReportCaller(true)
		for i := 0; i < b.N; i++ {
			logrus.WithContext(ctx).Info("info msg")
			if w.Len() > len(bs)-64 {
				w.Reset()
			}
		}
	})

	b.Run("logrus", func(b *testing.B) {
		logrus.SetReportCaller(false)
		for i := 0; i < b.N; i++ {
			WithContext(ctx).Info("info msg")
			if w.Len() > len(bs)-64 {
				w.Reset()
			}
		}
	})

	b.Run("logrus+lxt caller", func(b *testing.B) {
		logrus.SetReportCaller(false)
		for i := 0; i < b.N; i++ {
			logrus.WithContext(ctx).Info("info msg")
			if w.Len() > len(bs)-64 {
				w.Reset()
			}
		}
	})
}
