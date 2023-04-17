package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestTagTry0(t *testing.T) {
	defer func() {
		fmt.Printf("1 -> ")
	}()
	tag, err := NewTag() // 当 tag.Try(err) 时，跳转此处并返回 err1
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
	tag.Try(errors.New("err"))

	fmt.Printf("6 -> ")
	return
}
func TestTagTry(t *testing.T) {
	t.Run("NewLine1", func(t *testing.T) {
	gototag:
		defer func() {
			fmt.Printf("2")
		}()
		tag, err1 := NewTag() // 当 tag.Try(err) 时，跳转此处并返回 err1
		if err1 != nil {
			return
		}
		defer func() {
			fmt.Printf("1")
		}()
		err := errors.New("err")
		tag.Try(err)

		return
		goto gototag
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
				// fJump, err1 := NewTag()
				tag, err1 := NewTag()
				if err1 != nil {
					t.Log("3")
					err = err1
					t.Log("Tag() get error:", err1)
					// tag.Try(err3)
					return
				}
				defer func() {
					t.Log("inner defer 2")
				}()
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				// GotoTag(err3)
				// TryJump(err3)
				// fJump(err3)
				tag.Try(nil)
				_ = func() {
					tryTagErr = func(err error) {
						t.Fatal(err)
					}
					tag.Try(err3)
				}
				tag.Try(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}

func BenchmarkTag(b *testing.B) {
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

	b.Run("NewTag&Try", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag()
			if err1 != nil {
				continue
			}
			tag.Try(err3)
		}
		b.StopTimer()
	})
	b.Run("NewTag&Try(nil)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag()
			if err1 != nil {
				continue
			}
			tag.Try(nil)
		}
		b.StopTimer()
	})
	b.Run("Try(nil)", func(b *testing.B) {
		b.ReportAllocs()
		tag, err1 := NewTag()
		count := 0
		if err1 != nil {
			b.Fatal("never goto here")
		}
		if count++; count > 1 {
			b.Fatal("never goto here")
		}
		for i := 0; i < b.N; i++ {
			tag.Try(nil)
		}
		b.StopTimer()
	})

	b.Run("NewTag&Try-defer", func(b *testing.B) {
		b.ReportAllocs()
		defer func() {}()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag()
			if err1 != nil {
				continue
			}
			tag.Try(nil)
		}
		b.StopTimer()
	})

	b.Run("NewTag&Try-defer-notinline", func(b *testing.B) {
		b.ReportAllocs()
		defer func() {}()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag()
			if err1 != nil {
				for false {
					defer func() {}()
				}
				continue
			}
			tag.Try(nil)
		}
		b.StopTimer()
	})
}
