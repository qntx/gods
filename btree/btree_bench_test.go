package btree

import (
	"math/rand"
	"testing"
)

func BenchmarkDeleteAndRestoreG(b *testing.B) {
	items := rand.Perm(16392)

	b.ResetTimer()
	b.Run(`CopyBigFreeList`, func(b *testing.B) {
		fl := NewFreeListG[int](16392)

		tr := NewWithFreeListG(*btreeDegree, Less[int](), fl)
		for _, v := range items {
			tr.ReplaceOrInsert(v)
		}

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			dels := make([]int, 0, tr.Len())
			tr.Ascend(func(b int) bool {
				dels = append(dels, b)

				return true
			})

			for _, del := range dels {
				tr.Delete(del)
			}
			// tr is now empty, we make a new empty copy of it.
			tr = NewWithFreeListG(*btreeDegree, Less[int](), fl)
			for _, v := range items {
				tr.ReplaceOrInsert(v)
			}
		}
	})
	b.Run(`Copy`, func(b *testing.B) {
		tr := NewOrderedG[int](*btreeDegree)
		for _, v := range items {
			tr.ReplaceOrInsert(v)
		}

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			dels := make([]int, 0, tr.Len())
			tr.Ascend(func(b int) bool {
				dels = append(dels, b)

				return true
			})

			for _, del := range dels {
				tr.Delete(del)
			}
			// tr is now empty, we make a new empty copy of it.
			tr = NewOrderedG[int](*btreeDegree)
			for _, v := range items {
				tr.ReplaceOrInsert(v)
			}
		}
	})
	b.Run(`ClearBigFreelist`, func(b *testing.B) {
		fl := NewFreeListG[int](16392)

		tr := NewWithFreeListG(*btreeDegree, Less[int](), fl)
		for _, v := range items {
			tr.ReplaceOrInsert(v)
		}

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			tr.Clear(true)

			for _, v := range items {
				tr.ReplaceOrInsert(v)
			}
		}
	})
	b.Run(`Clear`, func(b *testing.B) {
		tr := NewOrderedG[int](*btreeDegree)
		for _, v := range items {
			tr.ReplaceOrInsert(v)
		}

		b.ReportAllocs()
		b.ResetTimer()

		for range b.N {
			tr.Clear(true)

			for _, v := range items {
				tr.ReplaceOrInsert(v)
			}
		}
	})
}
