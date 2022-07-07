package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var pc uintptr

func TestJump1(t *testing.T) {
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
				tag, err1 := NewTag2()
				if err1 != nil {
					t.Log("3")
					err = err1
					t.Log("Tag() get error:", err1)
					tag.Try(err3)
					return
				}
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				// GotoTag(err3)
				// TryJump(err3)
				// fJump(err3)
				tag.Try(nil)
				func() {
					tryTagErr = func(err error) {
						t.Fatal(err)
					}
					tag.Try(err3)
				}()
				tag.Try(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}
func TestJump(t *testing.T) {
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
				if err1 := Tag(); err1 != nil {
					t.Log("3")
					err = err1
					t.Log("Tag() get error:", err1)
					return
				}
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				GotoTag(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}

func TestTag(t *testing.T) {
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
				if err1 := Tag(); err1 != nil {
					t.Log("3")
					err = err1
					t.Log("Tag() get error:", err1)
					return
				}
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				GotoTag(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}

var err3 = fmt.Errorf("error 3")

func TagTest() (str string) {
	if err1 := Tag(); err1 != nil {
		str += "Tagx"
		return
	}
	GotoTag(err3)
	str += "GotoTagx(err3)"
	return
}

func TestTag2(t *testing.T) {
	t.Run("NewLine", func(t *testing.T) {
		str := TagTest()
		t.Log(str)
		assert.Panics(t, func() {
			GotoTag(nil)
		})
	})
	t.Run("NewLine", func(t *testing.T) {
		for i := 0; i < 3; i++ {
			func() {
				defer func() {
					t.Log("outer defer")
				}()
				err := func() (err error) {
					defer func() {
						t.Log("inner defer")
					}()
					t.Log("1")
					if err1 := Tag(); err1 != nil {
						// if err1 := tag2(tag); err1 != nil {
						t.Log("3")
						err = err1
						t.Log("Tag() get error:", err1)
						return
					}
					// }
					t.Log("2")
					err3 := fmt.Errorf("error 3")
					GotoTag(err3)
					t.Log("4")
					return
				}()
				t.Log("outer err:", err)
			}()
			assert.Panics(t, func() {
				GotoTag(nil)
			})
			err2 := func() (err error) {
				defer func() {
					t.Log("inner defer")
				}()
				t.Log("1")
				// tag := newTag()
				if i == 0 {
					if err1 := Tag(); err1 != nil {
						// if err1 := tag2(tag); err1 != nil {
						t.Log("3")
						err = err1
						t.Log("Tag() get error:", err1)
						return
					}
				}
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				// GotoTag2(tag, err3)
				GotoTag(err3)
				t.Log("4")
				return
			}()
			t.Log("outer err:", err2)
		}
	})
}

func BenchmarkTag(b *testing.B) {
	err3 := fmt.Errorf("error 3")
	b.Run("NewTag2&Try", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag2()
			if err1 != nil {
				continue
			}
			tag.Try(err3)
		}
		b.StopTimer()
	})
	b.Run("NewTag2&Try(nil)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			tag, err1 := NewTag2()
			if err1 != nil {
				continue
			}
			tag.Try(nil)
		}
		b.StopTimer()
	})
	b.Run("Try(nil)", func(b *testing.B) {
		b.ReportAllocs()
		tag, err1 := NewTag2()
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
	b.Run("Tagx&GotoTagx", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			if err1 := Tag(); err1 != nil {
				continue
			}
			GotoTag(err3)
		}
		b.StopTimer()
	})
	b.Run("Tagx", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Tag()
		}
		b.StopTimer()
	})
	b.Run("GotoTagx(nil)", func(b *testing.B) {
		b.ReportAllocs()
		count := 0
		if err := Tag(); err != nil {
			b.Fatal("never goto here")
		}
		if count++; count > 1 {
			b.Fatal("never goto here")
		}

		for i := 0; i < b.N; i++ {
			GotoTag(nil)
		}
		b.StopTimer()
	})
}
