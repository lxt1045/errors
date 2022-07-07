package errors

import (
	"sync"
)

func init() {
	// log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile) //log.Llongfile
}

var (
	mRoutineLastDefer = map[int64]struct{}{}
	lockRoutineDefer  sync.RWMutex
)

type guard struct {
	gid int64
	own bool
	noCopy
}

//go:noinline
func NewGuard() guard {

	gid := GetG().goid
	lockRoutineDefer.Lock()
	_, own := mRoutineLastDefer[gid]
	if !own {
		mRoutineLastDefer[gid] = struct{}{}
	}
	lockRoutineDefer.Unlock()

	return guard{
		gid: gid,
		own: own,
	}
}

func Catcher(g guard, f func(err interface{}) bool) {
	lockRoutineDefer.Lock()
	delete(mRoutineLastDefer, g.gid)
	lockRoutineDefer.Unlock()

	e := recover()
	if e == nil {
		return
	}

	// err, ok := e.(error)
	if f != nil && f(e) {
		return
	}
	panic(e)
}

func TryEscape(err *Cause) {
	gid := GetG().goid
	lockRoutineDefer.Lock()
	_, own := mRoutineLastDefer[gid]
	lockRoutineDefer.Unlock()
	if !own {
		panic("should call defer Catcher(NewGuard(),func()bool before call TryEscape(err))")
	}
	if err != nil {
		panic(err)
	}

	return
}
