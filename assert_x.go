package errors

import (
	"context"
	"log"
	"sync"

	"github.com/petermattis/goid"
)

/*
	通过ctx待err的方式监减少检查err
  if ctx.Get("err")!=nil{
	return
  }
*/

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile) //log.Llongfile
}

var (
	lockRoutineDefer  sync.RWMutex
	mRoutineLastDefer = make(map[int64]struct{})

	// 当goroutine没有defer时,会调用此函数,建议用于通知告警
	hookBeforePanic = func(ctx context.Context, err error) {
		<-ctx.Done()
	}
)

func SetHookBeforePanic(f func(context.Context, error)) {
	hookBeforePanic = f
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

func OKx(ctx context.Context, ok bool, err *Cause) {
	if ok {
		return
	}
	if err == nil {
		err = NewCause(1, DefaultCode, "not ok")
	}
	maybePanic(ctx, err)
}

func maybePanic(ctx context.Context, err *Cause) {
	lockRoutineDefer.RLock()
	_, ok := mRoutineLastDefer[goid.Get()]
	lockRoutineDefer.RUnlock()
	if ok {
		panic(err)
	}

	_defer := goid.GetDefer()
	if _defer != 0 {
		panic(err)
	}

	hookBeforePanic(ctx, err)
	return
}
