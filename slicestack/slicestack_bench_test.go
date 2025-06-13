package slicestack_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/slicestack"
)

func benchmarkPush(b *testing.B, stack *slicestack.Stack[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			stack.Push(key)
		}
	}
}

func benchmarkPop(b *testing.B, stack *slicestack.Stack[int], keys []int) {
	b.Helper()

	for range b.N {
		for range keys {
			stack.Pop()
		}
	}
}

func BenchmarkArrayStackPop100(b *testing.B) {
	b.StopTimer()

	size := 100
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		stack.Push(key)
	}

	b.StartTimer()
	benchmarkPop(b, stack, keys)
}

func BenchmarkArrayStackPop1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		stack.Push(key)
	}

	b.StartTimer()
	benchmarkPop(b, stack, keys)
}

func BenchmarkArrayStackPop10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		stack.Push(key)
	}

	b.StartTimer()
	benchmarkPop(b, stack, keys)
}

func BenchmarkArrayStackPop100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		stack.Push(key)
	}

	b.StartTimer()
	benchmarkPop(b, stack, keys)
}

func BenchmarkArrayStackPush100(b *testing.B) {
	b.StopTimer()

	size := 100
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPush(b, stack, keys)
}

func BenchmarkArrayStackPush1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPush(b, stack, keys)
}

func BenchmarkArrayStackPush10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPush(b, stack, keys)
}

func BenchmarkArrayStackPush100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	stack := slicestack.New[int]()
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPush(b, stack, keys)
}
