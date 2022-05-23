package errors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

// var pathSeparator = string([]byte{os.PathSeparator})
const pathSeparator = string(os.PathSeparator)

var (
	rootDirs = []string{"src/", "/pkg/mod/"} // file 会从rootDir开始截断

	// skipPkgs里的pkg会被忽略
	skipPrefixFiles = []string{
		"github.com/cloudwego/kitex",
		"testing/benchmark.go",
		"testing/testing.go",
	}

	pool = sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1024*2))
		},
	}
)

func MarshalJSON(err error) (bs []byte) {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	buf.WriteString(`{"wrapper":[`)
	cause := marshalJSON(buf, err)
	buf.WriteByte(']')
	buf.WriteString(`,"cause":`)
	marshalJSON2(buf, cause)
	buf.WriteString(`}`)
	bs = make([]byte, buf.Len())
	copy(bs, buf.Bytes())
	pool.Put(buf)
	return
}

func MarshalText(err error) (bs []byte) {
	buf := pool.Get().(*bytes.Buffer)
	buf.Reset()
	marshalText(buf, err)
	bs = make([]byte, buf.Len())
	copy(bs, buf.Bytes())
	pool.Put(buf)
	return
}

func marshalJSON(buf *bytes.Buffer, err error) (cause error) {
	errInner := errors.Unwrap(err)
	if errInner != nil {
		l := buf.Len()
		cause = marshalJSON(buf, errInner)
		if buf.Len() > l {
			buf.WriteByte(',')
		}
	} else {
		return err
	}
	marshalJSON2(buf, err)
	return
}

func marshalJSON2(buf *bytes.Buffer, err error) (cause error) {
	switch e := err.(type) {
	case *wrapper:
		e.json(buf)
	case *Cause:
		e.json(buf)
	case json.Marshaler:
		bs, err := e.MarshalJSON()
		if err != nil {
			panic(err)
		}
		buf.Write(bs)
	case fmt.Formatter:
		buf.WriteString(`{"trace":"`)
		buf.WriteString(fmt.Sprintf("%+v", err))
		buf.WriteString(`"}`)
	default:
		buf.WriteString(`{"trace":"`)
		buf.WriteString(e.Error())
		buf.WriteString(`"}`)
	}
	return
}

func marshalText(buf *bytes.Buffer, err error) {
	errInner := errors.Unwrap(err)
	if errInner != nil {
		marshalText(buf, errInner)
		buf.WriteByte('\n')
	}
	switch e := err.(type) {
	case *wrapper:
		e.text(buf)
	case *Cause:
		e.text(buf)
	case fmt.Formatter:
		buf.WriteString(`{"trace":"`)
		buf.WriteString(fmt.Sprintf("%+v", err))
		buf.WriteString(`"}`)
	default:
		buf.WriteString(`{"trace":"`)
		buf.WriteString(e.Error())
		buf.WriteString(`"}`)
	}
}

func toCaller(f runtime.Frame) caller { // nolint:gocritic
	funcName, file, line := f.Function, f.File, f.Line

	// /xxx@v0.0.3-0.20211019092134-6247f1f99488/...
	i := strings.Index(file, "@")
	if i > 0 {
		j := strings.Index(file[i+1:], pathSeparator)
		if j > 0 {
			file = file[:i] + file[i+1+j:]
		}
	}

	i = strings.LastIndex(funcName, pathSeparator)
	if i > 0 {
		rootDir := funcName[:i]
		funcName = funcName[i+1:]
		i = strings.Index(file, rootDir)
		if i > 0 {
			file = file[i:]
		}
	}
	if i <= 0 {
		for _, rootDir := range rootDirs {
			i = strings.Index(file, rootDir)
			if i > 0 {
				i += len(rootDir)
				file = file[i:]
				break
			}
		}
	}

	return caller{
		File: file,
		Line: line,
		Func: funcName, // 获取函数名
	}
}

type caller struct {
	File string
	Line int
	Func string
}

func (c caller) String() (s string) {
	line := strconv.Itoa(c.Line)
	return "(" + c.File + ":" + line + ") " + c.Func
}

func bufToString(buf *bytes.Buffer) string {
	bs := buf.Bytes()
	return *(*string)(unsafe.Pointer(&bs))
}

func skipFile(f string) bool {
	for _, skipPkg := range skipPrefixFiles {
		if strings.HasPrefix(f, skipPkg) {
			return true
		}
	}
	return false
}
