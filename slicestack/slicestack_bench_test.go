package slicestack_test

import (
	"testing"

	"github.com/qntx/gods/slicestack"
)

func benchmarkPush(b *testing.B, stack *slicestack.Stack[int], size int) {
	for range b.N {
		for n := range size {
			stack.Push(n)
		}
	}
}

func benchmarkPop(b *testing.B, stack *slicestack.Stack[int], size int) {
	for range b.N {
		for range size {
			stack.Pop()
		}
	}
}

func BenchmarkArrayStackPop100(b *testing.B) {
	b.StopTimer()

	size := 100
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkArrayStackPop1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkArrayStackPop10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkArrayStackPop100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPop(b, stack, size)
}

func BenchmarkArrayStackPush100(b *testing.B) {
	b.StopTimer()

	size := 100
	stack := slicestack.New[int]()

	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkArrayStackPush1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkArrayStackPush10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPush(b, stack, size)
}

func BenchmarkArrayStackPush100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	stack := slicestack.New[int]()

	for n := range size {
		stack.Push(n)
	}

	b.StartTimer()
	benchmarkPush(b, stack, size)
}
