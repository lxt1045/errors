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

package errors

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	pathSeparator = string(os.PathSeparator)
)

var (
	rootDirs = []string{"src/", "/pkg/mod/"} // file 会从rootDir开始截断

	// skipPkgs里的pkg会被忽略
	skipPrefixFiles = []string{
		"github.com/cloudwego/kitex",
		"testing/benchmark.go",
		"testing/testing.go",
	}
)

//MarshalJSON 将err序列化为json格式
func MarshalJSON(err error) (bs []byte) {
	buf := &writeBuffer{}
	marshalJSON(2, buf, err)
	bs = buf.Bytes()
	if bs[len(bs)-1] == ',' {
		bs = bs[:len(bs)-1]
	}
	return append(bs, `]}`...)
}

//marshalJSON 递归 Unwrap 并序列化为 JSON 格式
func marshalJSON(size int, buf *writeBuffer, err error) {
	switch e := err.(type) {
	//如果将 *wrapper 和 *Code 合成一个 interface{} 分支, 将导致性能退化
	case *Code:
		cache := e.fmt()
		buf.Grow(size + cache.jsonSize() + len(`{"cause":,"wrapper":[`))
		buf.WriteString(`{"cause":`)
		cache.json(buf)
		buf.WriteString(`,"wrapper":[`)
	case *wrapper:
		cache := e.fmt()
		needSize := cache.jsonSize() + 1
		marshalJSON(size+needSize, buf, e.Unwrap())
		cache.json(buf)
		buf.WriteByte(',')
		return
	case fmt.Formatter:
		cache := fmt.Sprintf("%+v", err)
		cacheSize, escape := countEscape(cache)
		buf.Grow(size + cacheSize + len(`{"cause":"","wrapper":[`))
		buf.WriteString(`{"cause":"`)
		if !escape {
			buf.WriteString(cache)
		} else {
			buf.WriteEscape(cache)
		}
		buf.WriteString(`","wrapper":[`)
	default:
		if err == nil {
			buf.Grow(size)
			return
		}
		cache := e.Error()
		cacheSize, escape := countEscape(cache)
		buf.Grow(size + cacheSize + len(`{"cause":"","wrapper":[`))
		buf.WriteString(`{"cause":"`)
		if !escape {
			buf.WriteString(cache)
		} else {
			buf.WriteEscape(cache)
		}
		buf.WriteString(`","wrapper":[`)
	}
}

func MarshalText(err error) (bs []byte) {
	buf := &writeBuffer{}
	marshalText(0, buf, err)
	return buf.Bytes()
}

func marshalText(size int, buf *writeBuffer, err error) {
	switch e := err.(type) {
	case *Code:
		cache := e.fmt()
		needSize := cache.textSize()
		buf.Grow(size + needSize)
		cache.text(buf)
	case *wrapper:
		cache := e.fmt()
		needSize := cache.textSize() + 1
		marshalText(size+needSize, buf, errors.Unwrap(err))
		buf.WriteByte('\n')
		cache.text(buf)
	case fmt.Formatter:
		cache := fmt.Sprintf("%+v", err)
		buf.Grow(size + len(cache) + 1)
		buf.WriteString(cache)
		buf.WriteByte(';')
	default:
		if err == nil {
			buf.Grow(size)
			return
		}
		cache := e.Error()
		buf.Grow(size + len(cache) + 1)
		buf.WriteString(cache)
		buf.WriteByte(';')
	}
}

type caller struct {
	File string
	Func string
}

func toCaller(f runtime.Frame) caller { // nolint:gocritic
	funcName, file, line := f.Function, f.File, f.Line

	// // /xxx@v0.0.3-0.20211019092134-6247f1f99488/...
	// i := strings.Index(file, "@")
	// if i > 0 {
	// 	j := strings.Index(file[i+1:], pathSeparator)
	// 	if j > 0 {
	// 		file = file[:i] + file[i+1+j:]
	// 	}
	// }

	i := strings.LastIndex(funcName, pathSeparator)
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
		File: file + ":" + strconv.Itoa(line),
		Func: funcName,
	}
}

func toCallers(pcs []uintptr) (callers []caller) {
	for f, i := runtime.CallersFrames(pcs), 0; ; i++ {
		ff, more := f.Next()
		callers = append(callers, toCaller(ff))
		if !more {
			break
		}
	}
	return
}

func (c caller) String() (s string) {
	if c.File == "" || c.Func == "" {
		return ""
	}
	return "(" + c.File + ") " + c.Func
}

func skipFile(f string) bool {
	for _, skipPkg := range skipPrefixFiles {
		if strings.HasPrefix(f, skipPkg) {
			return true
		}
	}
	return false
}

func MarshalJSON2(err error) (bs []byte) {
	bs = marshalJSON2(1, bs, err)
	bs[len(bs)-1] = ']'
	bs = append(bs, '}')

	return
}

func marshalJSON2(size int, bs []byte, err error) []byte {
	errInner := errors.Unwrap(err)
	switch e := err.(type) {
	case *wrapper:
		cache := e.fmt()
		if errInner != nil {
			needSize := cache.jsonSize() + 1
			bs = marshalJSON2(size+needSize, bs, errInner)
			bs = cache.json2(bs)
			bs = append(bs, ',')
			return bs
		}
		bs = tryGrow(bs, size+cache.jsonSize()+9+12)
		bs = append(bs, `{"cause":`...)
		bs = cache.json2(bs)
		bs = append(bs, `,"wrapper":[`...)
	case *Code:
		cache := e.fmt()
		if errInner != nil {
			needSize := cache.jsonSize() + 1
			bs = marshalJSON2(size+needSize, bs, errInner)
			bs = cache.json2(bs)
			bs = append(bs, ',')
			return bs
		}
		bs = tryGrow(bs, size+cache.jsonSize()+9+12)
		bs = append(bs, `{"cause":`...)
		bs = cache.json2(bs)
		bs = append(bs, `,"wrapper":[`...)
	case fmt.Formatter:
		cache := fmt.Sprintf("%+v", err)
		if errInner != nil {
			needSize := len(cache) + 10 + 3
			bs = marshalJSON2(size+needSize, bs, errInner)
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
			bs = marshalJSON2(size+needSize, bs, errInner)
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
func tryGrow(bs []byte, l int) []byte {
	if cap(bs) < l {
		bs2 := make([]byte, len(bs), l+len(bs))
		copy(bs2, bs)
		// fmt.Printf("tryGrow:%d\n", l)
		return bs2
	}
	return bs
}

func (f *fmtCode) json2(bs []byte) []byte {
	bs = append(bs, `{"code":`...)
	bs = append(bs, f.code...)
	bs = append(bs, `,"msg":"`...)
	bs = append(bs, f.msg...)
	bs = append(bs, '"')
	if len(f.stack) > 0 {
		bs = append(bs, `,"stack":`...)
		for i, str := range f.stack {
			if i != 0 {
				bs = append(bs, ',')
			}
			bs = append(bs, '"')
			if f.attr&(1<<i) == 0 {
				bs = append(bs, str...)
			} else {
				bs = append(bs, str...)
			}
			bs = append(bs, '"')
		}
	}
	bs = append(bs, '}')
	return bs
}

func (f *fmtWrapper) json2(bs []byte) []byte {
	bs = append(bs, `{"trace":"`...)
	bs = append(bs, f.trace...)
	bs = append(bs, `","caller":"`...)
	if !f.traceEscape {
		bs = append(bs, f.trace...)
	} else {
		bs = append(bs, f.trace...)
	}
	bs = append(bs, `"}`...)
	return bs
}
