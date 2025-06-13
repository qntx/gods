package pqueue_test

import (
	"testing"

	"github.com/qntx/gods/internal/testutil"
	"github.com/qntx/gods/pqueue"
)

func benchmarkEnqueue(b *testing.B, queue *pqueue.PriorityQueue[Element, int], priority []int) {
	b.Helper()

	for range b.N {
		for p := range priority {
			queue.Enqueue(Element{}, p)
		}
	}
}

func benchmarkDequeue(b *testing.B, queue *pqueue.PriorityQueue[Element, int], priority []int) {
	b.Helper()

	for range b.N {
		for range priority {
			queue.Dequeue()
		}
	}
}

func BenchmarkPriorityQueueDequeue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, priority)
}

func BenchmarkPriorityQueueDequeue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, priority)
}

func BenchmarkPriorityQueueDequeue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, priority)
}

func BenchmarkPriorityQueueDequeue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	for n := range size {
		queue.Enqueue(Element{}, n)
	}

	b.StartTimer()
	benchmarkDequeue(b, queue, priority)
}

func BenchmarkPriorityQueueEnqueue100(b *testing.B) {
	b.StopTimer()

	size := 100
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkEnqueue(b, queue, priority)
}

func BenchmarkPriorityQueueEnqueue1000(b *testing.B) {
	b.StopTimer()

	size := 1000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkEnqueue(b, queue, priority)
}

func BenchmarkPriorityQueueEnqueue10000(b *testing.B) {
	b.StopTimer()

	size := 10000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkEnqueue(b, queue, priority)
}

func BenchmarkPriorityQueueEnqueue100000(b *testing.B) {
	b.StopTimer()

	size := 100000
	queue := pqueue.New[Element, int](pqueue.MinHeap)
	priority := testutil.GeneratePermutedInts(size)

	b.StartTimer()
	benchmarkEnqueue(b, queue, priority)
}
