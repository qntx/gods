package slicedeque_test

import (
	"testing"

	"github.com/qntx/gods/slicedeque"
)

func benchmarkPushBack(b *testing.B, queue *slicedeque.Deque[int], size int) {
	b.Helper()

	for b.N > 0 {
		for n := range size {
			queue.PushBack(n)
		}
	}
}

func benchmarkPopFront(b *testing.B, queue *slicedeque.Deque[int], size int) {
	b.Helper()

	for b.N > 0 {
		for range size {
			queue.PopFront()
		}
	}
}

func BenchmarkArrayQueuePopFront100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePopFront100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPopFront(b, queue, size)
}

func BenchmarkArrayQueuePushBack100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := slicedeque.New[int](3)

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}

func BenchmarkArrayQueuePushBack100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := slicedeque.New[int](3)

	for n := range size {
		queue.PushBack(n)
	}

	b.StartTimer()
	benchmarkPushBack(b, queue, size)
}
