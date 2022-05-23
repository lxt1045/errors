package errors

import (
	"context"
	"log"
	"runtime"
	"sync"

	"github.com/petermattis/goid"
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile) //log.Llongfile
}

var (
	lockRoutineDefer  sync.RWMutex
	mRoutineLastDefer = make(map[int64]struct{})

	// HookNorRecoverPanic 如果主动 panic 前,本 goroutine 没有 recover(),
	// 则调用此方法建议在此处一直等待
	notRecoverPanicHook = func(ctx context.Context, err error) {
		<-ctx.Done()
		log.Printf("not recover panic: %+v", err)
		runtime.Goexit()
	}
)

func SetNotRecoverPanicHook(f func(context.Context, error)) {
	notRecoverPanicHook = f
}

func TryCatchx(fCatch func(interface{}) bool) func() {
	gid := goid.Get()
	lockRoutineDefer.Lock()
	_, ok := mRoutineLastDefer[gid]
	if !ok {
		mRoutineLastDefer[gid] = struct{}{}
	}
	lockRoutineDefer.Unlock()

	if ok {
		return func() {
			e := recover()
			if e == nil {
				return
			}
			if fCatch != nil && fCatch(e) {
				return
			}
			panic(e)
		}
	}

	return func() {
		lockRoutineDefer.Lock()
		delete(mRoutineLastDefer, gid)
		lockRoutineDefer.Unlock()

		e := recover()
		if e == nil {
			return
		}
		if fCatch != nil && fCatch(e) {
			return
		}
		panic(e)
	}
}

func OKx(ok bool, err *Cause) {
	if ok {
		return
	}
	if err == nil {
		err = buildCause(DefaultCode, "not ok", buildStack(1))
	}
	maybePanic(err)
}

func maybePanic(err *Cause) {
	lockRoutineDefer.RLock()
	_, ok := mRoutineLastDefer[goid.Get()]
	lockRoutineDefer.RUnlock()
	if ok {
		panic(err)
	}
	notRecoverPanicHook(context.Background(), err)
	panic(err)
}
