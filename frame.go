package errors

import (
	"bytes"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"unsafe"
)

var (
	// 调用栈名字的缓存
	mFuncs     = make(map[uintptr]caller)
	mFuncsLock sync.RWMutex
)

// 用于描述调用栈
type frame struct {
	trace string // 足迹

	rpc    [1]uintptr // runtime 的调用栈
	caller caller     // 解析后的缓存
}

// NewFrame 锚定调用栈
func NewFrame(skip int, trace string) *frame {
	f := buildFrame(skip+1, trace)
	return &f
}

func buildFrame(skip int, trace string) (s frame) {
	runtime.Callers(skip+1+baseSkip, s.rpc[:]) //nolint

	s.trace = trace
	return
}
func (s *frame) Caller() (caller string) {
	s.parse()
	return s.caller.String()
}

func (s *frame) parse() {
	if s.caller.FuncName != "" {
		return
	}
	mFuncsLock.RLock()
	ok := false
	s.caller, ok = mFuncs[s.rpc[0]]
	mFuncsLock.RUnlock()
	if ok {
		return
	}

	s.parseSlow()
	mFuncsLock.Lock()
	if _, ok := mFuncs[s.rpc[0]]; !ok {
		mFuncs[s.rpc[0]] = s.caller
	}
	mFuncsLock.Unlock()
}

func (s *frame) parseSlow() {
	f, _ := runtime.CallersFrames(s.rpc[:]).Next()
	s.caller = toCaller(f)
}

func (s *frame) json(buf *bytes.Buffer) {
	buf.Write([]byte(`{"trace":`))
	buf.WriteByte('"')
	buf.WriteString(s.trace)
	buf.WriteByte('"')
	buf.WriteByte(',')
	buf.Write([]byte(`"caller":`))
	buf.WriteByte('"')
	buf.WriteString(s.Caller())
	buf.WriteByte('"')

	buf.WriteByte('}')
}
func (s *frame) text(buf *bytes.Buffer) {
	if s.trace == "" {
		s.trace = "-"
	}
	buf.WriteString(s.trace)
	buf.Write([]byte(",\n    "))
	buf.WriteString(s.Caller())
}

func (s *frame) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	s.json(buf)
	return buf.Bytes(), nil
}
func (s *frame) String() string {
	buf := bytes.NewBuffer(nil)
	s.text(buf)
	return bufToString(buf)
}

func toCaller(f runtime.Frame) caller { // nolint:gocritic
	funcName := f.Function

	// /xxx@v0.0.3-0.20211019092134-6247f1f99488/...
	i := strings.Index(f.File, "@")
	if i > 0 {
		j := strings.Index(f.File[i+1:], pathSeparator)
		if j > 0 {
			f.File = f.File[:i] + f.File[i+1+j:]
		}
	}

	i = strings.LastIndex(funcName, pathSeparator)
	if i > 0 {
		rootDir := funcName[:i]
		funcName = funcName[i+1:]
		i = strings.Index(f.File, rootDir)
		if i > 0 {
			f.File = f.File[i:]
		}
	}
	if i <= 0 {
		for _, rootDir := range rootDirs {
			i = strings.Index(f.File, rootDir)
			if i > 0 {
				i += len(rootDir)
				f.File = f.File[i:]
				break
			}
		}
	}

	return caller{
		File:     f.File + ":" + strconv.Itoa(f.Line),
		FuncName: funcName, // 获取函数名
	}
}

type caller struct {
	File     string
	FuncName string
}

func (c caller) String() (s string) {
	buf := bytes.NewBuffer(nil)
	c.serialize(buf)
	return buf.String()
}

func (c caller) serialize(buf *bytes.Buffer) *bytes.Buffer {
	buf.WriteString("(")
	buf.WriteString(c.File)
	buf.WriteString(") ")
	buf.WriteString(c.FuncName)
	return buf
}

func bufToString(buf *bytes.Buffer) string {
	bs := buf.Bytes()
	return *(*string)(unsafe.Pointer(&bs))
}

func skipFunc(f string) bool {
	for _, skipPkg := range skipPkgs {
		if strings.Contains(f, skipPkg) {
			return true
		}
	}

	if strings.HasSuffix(f, "testing.tRunner") {
		return true
	}
	return false
}
