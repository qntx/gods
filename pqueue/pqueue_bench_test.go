package pqueue_test

import (
	"testing"

	"github.com/qntx/gods/pqueue"
)

func benchmarkEnqueue(b *testing.B, queue *pqueue.PriorityQueue[Element, int], size int) {
	for range b.N {
		for n := range size {
			queue.Enqueue(Element{}, n)
		}
	}
}

func benchmarkDequeue(b *testing.B, queue *pqueue.PriorityQueue[Element, int], size int) {
	for range b.N {
		for range size {
			queue.Dequeue()
		}
	}
}

func BenchmarkPriorityQueueDequeue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkPriorityQueueDequeue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkPriorityQueueDequeue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkPriorityQueueDequeue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, size)
}

func BenchmarkPriorityQueueEnqueue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkPriorityQueueEnqueue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkPriorityQueueEnqueue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}

func BenchmarkPriorityQueueEnqueue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := pqueue.New[Element, int](pqueue.MinHeap)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkEnqueue(b, queue, size)
}
