package errors

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/lxt1045/errors/jmp"
	"github.com/rs/zerolog"
)

func TestExample(t *testing.T) {
	ctx := context.TODO()

	err := func(ctx context.Context) (err error) {
		jump, err := Setjmp()
		if err != nil {
			return
		}
		err = func() (err error) {
			return nil
		}()
		jump.Longjmp(err)

		err = func() (err error) {
			return New("err")
		}()
		jump.Longjmp(err)

		return
	}(ctx)

	log := zerolog.New(os.Stdout)
	log.Info().Caller().Err(err).Msg("just for test")

}

func TestSetJMP(t *testing.T) {
	t.Run("Setjmp", func(t *testing.T) {
		// t.Error("err")
		func() {
			defer func() {
				t.Log("outer defer")
			}()
			err := func() (err error) {
				defer func() {
					t.Log("inner defer")
				}()
				t.Log("1")
				jump, err := Setjmp()
				t.Logf("defer:0x%x, defer_offset:%d, jump:%+v", jump._defer, defer_offset, jump)
				if err != nil {
					t.Log("3")
					t.Log("Setjmp() get error:", err)
					return
				}
				defer func() {
					t.Log("inner defer 2")
				}()
				t.Log("2")
				err3 := fmt.Errorf("error 3")
				jump.Longjmp(nil)
				t.Log("4")
				jump.Longjmp(err3)
				t.Log("5")
				return
			}()
			t.Log("outer err:", err)
		}()
	})
}

func BenchmarkSetJMP1(b *testing.B) {
	err := fmt.Errorf("error 0")

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
	b.Run("Check(err)", func(b *testing.B) {
		b.ReportAllocs()
		handler, err1 := NewHandler()
		for i := 0; i < b.N; i++ {
			if i == 0 {
				handler, err1 = NewHandler()
				if err1 != nil {
					continue
				}
			}
			handler.Check(err)
		}
		b.StopTimer()
	})

	b.Run("Longjmp(nil)", func(b *testing.B) {
		b.ReportAllocs()
		handler, err1 := Setjmp()
		count := 0
		if err1 != nil {
			b.Fatal("never goto here")
		}
		if count++; count > 1 {
			b.Fatal("never goto here")
		}
		for i := 0; i < b.N; i++ {
			handler.Longjmp(nil)
		}
		b.StopTimer()
	})
	b.Run("Longjmp(err)", func(b *testing.B) {
		b.ReportAllocs()
		handler, err1 := Setjmp()
		for i := 0; i < b.N; i++ {
			if i == 0 {
				handler, err1 = Setjmp()
				if err1 != nil {
					continue
				}
			}
			handler.Longjmp(err)
		}
		b.StopTimer()
	})
	b.Run("Setjmp()", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			Setjmp()
		}
		b.StopTimer()
	})
	b.Run("Setjmp(err)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			handler, err1 := Setjmp()
			if err1 != nil {
				continue
			}
			handler.Longjmp(err)
		}
		b.StopTimer()
	})

	b.Run("jmp.Try(nil)", func(b *testing.B) {
		b.ReportAllocs()
		pc, err1 := jmp.Set()
		count := 0
		if err1 != nil {
			b.Fatal("never goto here")
		}
		if count++; count > 1 {
			b.Fatal("never goto here")
		}
		for i := 0; i < b.N; i++ {
			jmp.Try(pc, nil)
		}
		b.StopTimer()
	})
	b.Run("jmp.Try(err)", func(b *testing.B) {
		b.ReportAllocs()
		pc, err1 := jmp.Set()
		for i := 0; i < b.N; i++ {
			if i == 0 {
				pc, err1 = jmp.Set()
				if err1 != nil {
					continue
				}
			}
			jmp.Try(pc, err)
		}
		b.StopTimer()
	})
	b.Run("jmp.Set()", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			jmp.Set()
		}
		b.StopTimer()
	})
	b.Run("jmp.Set(err)", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			pc, err1 := jmp.Set()
			if err1 != nil {
				continue
			}
			jmp.Try(pc, err)
		}
		b.StopTimer()
	})
}

func BenchmarkSetJMP(b *testing.B) {
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
