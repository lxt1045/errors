package errors

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Test_getPCSlow(t *testing.T) {
	t.Run("getPCSlow", func(t *testing.T) {
		pcs2, pcs1 := fPC()
		t.Logf("getPCSlow:%+v, getPC:%+v", pcs1, pcs2)
		assert.Equal(t, pcs2, pcs1)
	})
	t.Run("getPCSlow0", func(t *testing.T) {
		pcs2, pcs1 := fPC0()
		t.Logf("getPCSlow:%+v, getPC:%+v", pcs1, pcs2)
		assert.Equal(t, pcs2, pcs1)
	})

	t.Run("GetPCSlow", func(t *testing.T) {
		pcs1, pcs2 := fPC2()
		t.Logf("GetPCSlow:%+v, GetPC:%+v", pcs1, pcs2)
		assert.Equal(t, pcs2, pcs1)

		c1 := CallerFrame(pcs1)
		c2, _ := runtime.CallersFrames([]uintptr{pcs2}).Next()
		t.Logf("GetPCSlow:%s:%d, GetPC:%s:%d", c1.File, c1.Line, c2.File, c2.Line)
	})

	t.Run("buildStack", func(t *testing.T) {
		pcs1, pcs2 := make([]uintptr, 3), make([]uintptr, 3)
		n1, n2 := fBuildStack(pcs1, pcs2)
		t.Logf("buildStack:%+v, buildStackSlow:%+v", pcs1, pcs2)
		printfLine := func(pcs []uintptr) {
			fs := runtime.CallersFrames(pcs)
			for {
				c1, more := fs.Next()
				t.Logf("line:%s:%d", c1.File, c1.Line)
				if !more {
					break
				}
			}
		}
		printfLine(pcs1)
		printfLine(pcs2)

		assert.Equal(t, pcs2, pcs1)
		assert.Equal(t, n1, n2)
	})
}

//go:noinline
func fPC() ([1]uintptr, [1]uintptr) {
	return getPC(), getPCSlow()
}
func fPC0() ([1]uintptr, [1]uintptr) {
	return getPCSlow(), getPC()
}

//go:noinline
func fPC2() (uintptr, uintptr) {
	var getPC2 func() [1]uintptr = getPCSlow
	var GetPC2 func() uintptr = *(*func() uintptr)(unsafe.Pointer(&getPC2))
	return GetPC2(), GetPC()
}

//go:noinline
func fBuildStack(in1, in2 []uintptr) (n1, n2 int) {
	return fBuildStack2(in1, in2)
}

//go:noinline
func fBuildStack2(in1, in2 []uintptr) (n1, n2 int) {
	return buildStack(in1), buildStackSlow(in2)
}

func Test_CallersSkip(t *testing.T) {
	t.Run("CallersSkip", func(t *testing.T) {
		func() {
			func() {
				func() {
					cs := CallersSkip(1)
					t.Logf("callers:%+v", cs)
				}()
			}()
		}()
	})
}

func BenchmarkCallersSkip(b *testing.B) {
	b.Run("CallersSkip", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			func() {
				func() {
					func() {
						CallersSkip(0)
					}()
				}()
			}()
		}
		b.StopTimer()
	})
	b.Run("runtime.Callers", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := pool.Get().(*[DefaultDepth]uintptr)
			func() {
				func() {
					func() {
						runtime.Callers(baseSkip, pcs[:DefaultDepth])
						// parseSlow(pcs[:n])
					}()
				}()
			}()
			pool.Put(pcs)
		}
		b.StopTimer()
	})
	b.Run("runtime.Callers & runtime.CallersFrames", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pcs := pool.Get().(*[DefaultDepth]uintptr)
			func() {
				func() {
					func() {
						n := runtime.Callers(baseSkip, pcs[:DefaultDepth])
						parseSlow(pcs[:n])
					}()
				}()
			}()
			pool.Put(pcs)
		}
		b.StopTimer()
	})
}
