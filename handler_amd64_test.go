package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestHandlerCheck0(t *testing.T) {
	defer func() {
		fmt.Printf("1 -> ")
	}()
	handler, err := NewHandler() // 当 handler.Check(err) 时，跳转此处并返回 err1
	fmt.Printf("2 -> ")
	if err != nil {
		fmt.Printf("3 -> ")
		return
		// 空的时候，对比一下生成的代码有什么区别
		// 参考： https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md
		// 可能原因是：defer 在 loop 中，导致编译器对 defer 内联优化策略的改变！
		// 逃逸分析失效，被迫当做 defer 逃逸
		/*
			src/cmd/compile/internal/walk/stmt.go:116
			// If n.Esc is not EscNever, then this defer occurs in a loop,
			// so open-coded defers cannot be used in this function.

			#define FUNCDATA_OpenCodedDeferInfo 4 // info for func with open-coded defers

			参考 runOpenDeferFrame()、addOneOpenDeferFrame() 函数调用 open-coded defers 函数？

			搜索关键词： "open-coded defers"
		*/
		for false {
			defer func() {}()
		}
	}

	defer func() {
		fmt.Printf("4 -> ") // 由于的缺陷：这里 debug 下 defer 不内联，会执行；release 下 defer 内联，不会执行
	}()

	fmt.Printf("5 -> ")
	handler.Check(errors.New("err"))

	fmt.Printf("6 -> ")
	return
}

func TestHandlerCheck(t *testing.T) {
	t.Run("NewLine1", func(t *testing.T) {
	gotohandler:
		defer func() {
			fmt.Printf("2")
		}()
		handler, err1 := NewHandler() // 当 handler.Check(err) 时，跳转此处并返回 err1
		if err1 != nil {
			return
		}
		defer func() {
			fmt.Printf("1")
		}()
		err := errors.New("err")
		handler.Check(err)

		return
		goto gotohandler
	})
	return
	t.Run("NewLine", func(t *testing.T) {
		func() {
			defer func() {
				t.Log("outer defer")
			}()
			err := func() (err error) {
				defer func() {
					t.Log("inner defer")
				}()
				t.Log("1")
				// fJump, err1 := NewHandler()
				handler, err1 := NewHandler()
				if err1 != nil {
					t.Log("3")
					err = err1
					t.Log("Handler() get error:", err1)
					// handler.Check(err3)
					return
				}
				defer func() {
					t.Log("inner defer 2")
				}()
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				// GotoHandler(err3)
				// CheckJump(err3)
				// fJump(err3)
				handler.Check(nil)
				_ = func() {
					tryHandlerErr = func(err error) {
						t.Fatal(err)
					}
					handler.Check(err3)
				}
				handler.Check(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}

func BenchmarkHandler(b *testing.B) {
	err3 := fmt.Errorf("error 3")

	b.Run("defer&panic", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		_ = err
		for i := 0; i < b.N; i++ {
			func() {
				defer func() {
					e := recover()
					if e != nil {
						err = e.(error)
						return
					}
				}()
				if err3 != nil {
					panic(err3)
				}
			}()
		}
		b.StopTimer()
	})
	b.Run("defer&panic-nil", func(b *testing.B) {
		b.ReportAllocs()
		var err error
		_ = err
		defer func() {
			e := recover()
			if e != nil {
				err = e.(error)
				return
			}
		}()
		for i := 0; i < b.N; i++ {
			if false {
				panic(err3)
			}
		}
		b.StopTimer()
	})

	b.Run("NewHandler&Check", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			handler, err1 := NewHandler()
			if err1 != nil {
				continue
			}
			handler.Check(err3)
		}
		b.StopTimer()
	})
	b.Run("NewHandler&Check(nil)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			handler, err1 := NewHandler()
			if err1 != nil {
				continue
			}
			handler.Check(nil)
		}
		b.StopTimer()
	})
	b.Run("Check(nil)", func(b *testing.B) {
		b.ReportAllocs()
		handler, err1 := NewHandler()
		count := 0
		if err1 != nil {
			b.Fatal("never goto here")
		}
		if count++; count > 1 {
			b.Fatal("never goto here")
		}
		for i := 0; i < b.N; i++ {
			handler.Check(nil)
		}
		b.StopTimer()
	})

	b.Run("NewHandler&Check-defer", func(b *testing.B) {
		b.ReportAllocs()
		defer func() {}()
		for i := 0; i < b.N; i++ {
			handler, err1 := NewHandler()
			if err1 != nil {
				continue
			}
			handler.Check(nil)
		}
		b.StopTimer()
	})

	b.Run("NewHandler&Check-defer-notinline", func(b *testing.B) {
		b.ReportAllocs()
		defer func() {}()
		for i := 0; i < b.N; i++ {
			handler, err1 := NewHandler()
			if err1 != nil {
				for false {
					defer func() {}()
				}
				continue
			}
			handler.Check(nil)
		}
		b.StopTimer()
	})
}
