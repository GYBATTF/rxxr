package rxxr

import (
	"math"
	"testing"
)

func BenchmarkPipe_Publish(b *testing.B) {
	p := New[int](new(PipeConfig[int]).SetPanicLogger(func(r any) {
		b.Fail()
		b.Log("benchmark panicked")
	}))

	subscribers := math.MaxInt64
	for i := 0; i < subscribers; i++ {
		p.Subscribe(bencharkDoSomething(b))
	}

	for i := 0; i < b.N; i++ {
		p.Publish(i)
	}
}

func bencharkDoSomething(b *testing.B) func(int) {
	return func(i int) {
		i++
	}
}
