package pqueue

import (
	"testing"
)

type Element struct {
	// priority int
	// name string
}

func benchmarkEnqueue(b *testing.B, queue *PriorityQueue[Element, int], size int) {
	for range b.N {
		for n := range size {
			queue.Put(Element{}, n)
		}
	}
}

func benchmarkDequeue(b *testing.B, queue *PriorityQueue[Element, int], size int) {
	for range b.N {
		for range size {
			queue.Get()
		}
	}
}

func BenchmarkBinaryQueueDequeue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueDequeue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkBinaryQueueEnqueue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := New[Element, int](MinHeap)

	for n := range size {
		queue.Put(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}
