package errors

import (
	"runtime"
	"strconv"
	"strings"
	"sync"
)

const DefaultLevel = 16

var (
	rootDir = "/" //file 会从rootDir开始截断
)

//用于描述调用栈
type stack struct {
	nrpc     int                   //runtime 的调用栈
	rpcCache [DefaultLevel]uintptr //runtime 的调用栈缓存,避免一次内存分配

	callers []caller //解析后的缓存
}

//NewStack 锚定调用栈
func NewStack(skip int, level int) (s stack) {
	return newStack(skip+1, level)
}
func newStack(skip, depth int) (s stack) {
	if depth > len(s.rpcCache) {
		depth = len(s.rpcCache)
	} else if depth <= 0 {
		depth = DefaultLevel
	}
	s.nrpc = runtime.Callers(skip+2, s.rpcCache[:depth])
	return
}

func (s *stack) Callers() (callers []string) {
	s.parse()
	for _, f := range s.callers {
		callers = append(callers, f.String())
	}
	return
}

var (
	mCallers     = make(map[[DefaultLevel]uintptr][]caller)
	mCallersLock sync.RWMutex
)

func (s *stack) parse() {
	if len(s.callers) <= 0 {
		mCallersLock.RLock()
		ok := false
		s.callers, ok = mCallers[s.rpcCache]
		mCallersLock.RUnlock()

		if !ok {
			s.callers = s.parseSlow() //这步放在Lock()外虽然可能会造成重复计算,但是极大减少了锁争抢
			mCallersLock.Lock()
			if _, ok = mCallers[s.rpcCache]; !ok {
				mCallers[s.rpcCache] = s.callers
			}
			mCallersLock.Unlock()
		}
		return
	}
	return
}

func (s *stack) parseSlow() (callers []caller) {
	if len(s.callers) <= 0 {
		frames, more, f := runtime.CallersFrames(s.rpcCache[:s.nrpc]), true, runtime.Frame{}
		for more {
			f, more = frames.Next()
			if strings.Contains(f.Function, "github.com/cloudwego/kitex") && len(s.callers) > 0 {
				break
			}
			//TODO 这个特化逻辑不能省略,将来可以作为初始化参数
			if strings.HasSuffix(f.Function, "testing.tRunner") && len(s.callers) > 0 {
				break
			}
			s.callers = append(s.callers, newCaller(f))
			if strings.HasSuffix(f.Function, "main.main") && len(s.callers) > 0 {
				break
			}
		}

		if len(s.callers) <= 0 {
			s.callers = []caller{{File: "nil", FuncName: "nil"}}
		}
		// else if len(s.callers) > 1 {
		// 	//合并最长前缀
		// 	last := strings.Split(s.callers[0].File, "/")
		// 	for i := 1; i < len(s.callers); i++ {
		// 		cur := strings.Split(s.callers[i].File, "/")
		// 		j := 0
		// 		for j = range cur {
		// 			if last[j] != cur[j] {
		// 				break
		// 			}
		// 		}
		// 		s.callers[i].File = strings.Join(cur[j:], "/")
		// 		last = last[:j]
		// 		if len(cur) > j {
		// 			last = append(last, cur[j:]...)
		// 		}
		// 	}
		// }
	}
	return s.callers
}

func (s stack) json(buf *strings.Builder) *strings.Builder {
	callers := s.Callers()

	if len(callers) == 1 {
		buf.Write([]byte(`"caller":`))
		buf.WriteByte('"')
		buf.WriteString(callers[0])
		buf.WriteByte('"')
	} else {
		buf.Write([]byte(`"callers":[`))
		for i, caller := range callers {
			if i != 0 {
				buf.Write([]byte(","))
			}
			buf.WriteByte('"')
			buf.WriteString(caller)
			buf.WriteByte('"')
		}
		buf.WriteByte(']')
	}

	return buf
}
func (s stack) text(buf *strings.Builder) *strings.Builder {
	callers := s.Callers()
	for i, caller := range callers {
		if i != 0 {
			buf.Write([]byte(", \n"))
		}
		buf.Write([]byte("    "))
		buf.WriteString(caller)
	}

	return buf
}

func (s stack) String() string {
	buf := new(strings.Builder)
	if Layout == LayoutTypeJSON {
		return s.json(buf).String()
	}
	return s.text(buf).String()
}

func newCaller(f runtime.Frame) (c caller) {
	funcName := f.Function //.FuncForPC(f.PC).Name()

	// /@v0.0.3-0.20211019092134-6247f1f99488/...
	i := strings.Index(f.File, "@")
	if i > 0 {
		j := strings.Index(f.File[i+1:], "/")
		if j > 0 {
			f.File = f.File[:i] + f.File[i+1+j:]
		}
	}

	i = strings.LastIndex(funcName, "/")
	if i > 0 {
		rootDir := funcName[:i]
		funcName = funcName[i+1:]
		i = strings.Index(f.File, rootDir)
		if i > 0 {
			f.File = f.File[i:]
		}
	}
	if i <= 0 {
		rootDirs := []string{"src", "/pkg/mod/"}
		for _, rootDir := range rootDirs {
			i = strings.Index(f.File, rootDir)
			if i > 0 {
				i += len(rootDir)
				f.File = f.File[i:]
				break
			}
		}
	}

	c = caller{
		File:     f.File + ":" + strconv.Itoa(f.Line),
		FuncName: funcName, //获取函数名
	}
	return
}

type caller struct {
	File     string `json:"file"`
	FuncName string `json:"func"`
}

func (c caller) String() (s string) {
	buf := new(strings.Builder)
	c.serialize(buf)
	return buf.String()
}

func (c caller) serialize(buf *strings.Builder) *strings.Builder {
	buf.WriteString("(")
	buf.WriteString(c.File)
	buf.WriteString(") ")
	buf.WriteString(c.FuncName)
	return buf
}
