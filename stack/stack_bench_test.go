package stack

import (
	"testing"

	lxtstack "github.com/lxt1045/errors"
)

func BenchmarkNewStack1(b *testing.B) {
	b.Run("NewStack", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			s := NewStack(2, 16)
			s.ReclaimCache()
		}
		b.StopTimer()
	})
	b.Run("lxtstack.NewStack", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			lxtstack.NewStack(2, 16)
		}
		b.StopTimer()
	})
}
