package slicedeque_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/slicedeque"
)

func benchmarkPushBack(b *testing.B, queue *slicedeque.Deque[int], keys []int) {
	b.Helper()

	for range b.N {
		for key := range keys {
			queue.PushBack(key)
		}
	}
}

func benchmarkPopFront(b *testing.B, queue *slicedeque.Deque[int], keys []int) {
	b.Helper()

	for range b.N {
		for range keys {
			queue.PopFront()
		}
	}
}

func BenchmarkArrayQueuePopFront100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		queue.PushBack(key)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, keys)
}

func BenchmarkArrayQueuePopFront1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		queue.PushBack(key)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, keys)
}

func BenchmarkArrayQueuePopFront10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		queue.PushBack(key)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, keys)
}

func BenchmarkArrayQueuePopFront100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	for key := range keys {
		queue.PushBack(key)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, keys)
}

func BenchmarkArrayQueuePushBack100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPushBack(b, queue, keys)
}

func BenchmarkArrayQueuePushBack1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPushBack(b, queue, keys)
}

func BenchmarkArrayQueuePushBack10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPushBack(b, queue, keys)
}

func BenchmarkArrayQueuePushBack100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := slicedeque.New[int](3)
	keys := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkPushBack(b, queue, keys)
}
