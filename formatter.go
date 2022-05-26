package errors

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
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
)

func MarshalJSON(err error) (bs []byte) {
	bs = marshalJSON(1, bs, err)
	bs[len(bs)-1] = ']'
	bs = append(bs, '}')

	return
}

func marshalJSON(size int, bs []byte, err error) []byte {
	errInner := errors.Unwrap(err)
	switch e := err.(type) {
	case *wrapper:
		cache := e.fmt()
		if errInner != nil {
			needSize := cache.jsonSize() + 1
			bs = marshalJSON(size+needSize, bs, errInner)
			bs = cache.json(bs)
			bs = append(bs, ',')
			return bs
		}
		bs = tryGrow(bs, size+cache.jsonSize()+9+12)
		bs = append(bs, `{"cause":`...)
		bs = cache.json(bs)
		bs = append(bs, `,"wrapper":[`...)
	case *Cause:
		cache := e.fmt()
		if errInner != nil {
			needSize := cache.jsonSize() + 1
			bs = marshalJSON(size+needSize, bs, errInner)
			bs = cache.json(bs)
			bs = append(bs, ',')
			return bs
		}
		bs = tryGrow(bs, size+cache.jsonSize()+9+12)
		bs = append(bs, `{"cause":`...)
		bs = cache.json(bs)
		bs = append(bs, `,"wrapper":[`...)
	case fmt.Formatter:
		cache := fmt.Sprintf("%+v", err)
		if errInner != nil {
			needSize := len(cache) + 10 + 3
			bs = marshalJSON(size+needSize, bs, errInner)
			bs = append(bs, `{"trace":"`...)
			bs = append(bs, cache...)
			bs = append(bs, `"},`...)
			return bs
		}
		bs = tryGrow(bs, size+len(cache)+10+13)
		bs = append(bs, `{"cause":"`...)
		bs = append(bs, cache...)
		bs = append(bs, `","wrapper":[`...)
	default:
		cache := e.Error()
		if errInner != nil {
			needSize := len(cache) + 10 + 3
			bs = marshalJSON(size+needSize, bs, errInner)
			bs = append(bs, `{"trace":"`...)
			bs = append(bs, cache...)
			bs = append(bs, `"},`...)
			return bs
		}
		bs = tryGrow(bs, size+len(cache)+10+13)
		bs = append(bs, `{"cause":"`...)
		bs = append(bs, cache...)
		bs = append(bs, `","wrapper":[`...)
	}
	return bs
}

func MarshalText(err error) (bs []byte) {
	bs = marshalText(0, bs, err)
	return
}
func tryGrow(bs []byte, l int) []byte {
	if cap(bs) < l {
		bs2 := make([]byte, len(bs), l+len(bs))
		copy(bs2, bs)
		// fmt.Printf("tryGrow:%d\n", l)
		return bs2
	}
	return bs
}
func marshalText(size int, bs []byte, err error) []byte {
	errInner := errors.Unwrap(err)
	switch e := err.(type) {
	case *wrapper:
		cache := e.fmt()
		needSize := cache.textSize() + 1
		if errInner != nil {
			bs = marshalText(size+needSize, bs, errInner)
			bs = append(bs, '\n')
			return cache.text(bs)
		}
		bs = tryGrow(bs, size+needSize)
	case *Cause:
		cache := e.fmt()
		needSize := cache.textSize() + 1
		if errInner != nil {
			bs = marshalText(size+needSize, bs, errInner)
			bs = append(bs, '\n')
			return cache.text(bs)
		}
		bs = tryGrow(bs, size+needSize)
		bs = cache.text(bs)
	case fmt.Formatter:
		cache := fmt.Sprintf("%+v", err)
		needSize := len(cache) + 5 + 1
		if errInner != nil {
			bs = marshalText(size+needSize, bs, errInner)
			bs = append(bs, "\n    "...)
			return append(bs, cache...)
		}
		bs = tryGrow(bs, size+needSize)
		bs = append(bs, cache...)
	default:
		cache := e.Error()
		needSize := len(cache) + 5 + 1
		if errInner != nil {
			bs = marshalText(size+needSize, bs, errInner)
			bs = append(bs, "\n    "...)
			return append(bs, cache...)
		}
		bs = tryGrow(bs, size+needSize)
		bs = append(bs, cache...)
	}
	return bs
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

func skipFile(f string) bool {
	for _, skipPkg := range skipPrefixFiles {
		if strings.HasPrefix(f, skipPkg) {
			return true
		}
	}
	return false
}
