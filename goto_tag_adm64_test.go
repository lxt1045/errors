package errors

import (
	"fmt"
	"testing"
)

var pc uintptr

func TestTagTry(t *testing.T) {
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

var err3 = fmt.Errorf("error 3")

func BenchmarkTag(b *testing.B) {
	err3 := fmt.Errorf("error 3")
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
}
