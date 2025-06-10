package btree

import (
	"flag"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"sync"
	"testing"
)

var btreeDegree = flag.Int("degree", 32, "B-Tree degree")

func intRange(s int, reverse bool) []int {
	out := make([]int, s)

	for i := range s {
		v := i
		if reverse {
			v = s - i - 1
		}

		out[i] = v
	}

	return out
}

func intAll(t *BTreeG[int]) (out []int) {
	t.Ascend(func(a int) bool {
		out = append(out, a)

		return true
	})

	return
}

func intAllRev(t *BTreeG[int]) (out []int) {
	t.Descend(func(a int) bool {
		out = append(out, a)

		return true
	})

	return
}

func TestBTreeG(t *testing.T) {
	tr := NewOrderedG[int](*btreeDegree)

	const treeSize = 10000

	for range 10 {
		if min, ok := tr.Min(); ok || min != 0 {
			t.Fatalf("empty min, got %+v", min)
		}

		if max, ok := tr.Max(); ok || max != 0 {
			t.Fatalf("empty max, got %+v", max)
		}

		for _, item := range rand.Perm(treeSize) {
			if x, ok := tr.ReplaceOrInsert(item); ok || x != 0 {
				t.Fatal("insert found item", item)
			}
		}

		for _, item := range rand.Perm(treeSize) {
			if x, ok := tr.ReplaceOrInsert(item); !ok || x != item {
				t.Fatal("insert didn't find item", item)
			}
		}

		want := 0
		if min, ok := tr.Min(); !ok || min != want {
			t.Fatalf("min: ok %v want %+v, got %+v", ok, want, min)
		}

		want = treeSize - 1
		if max, ok := tr.Max(); !ok || max != want {
			t.Fatalf("max: ok %v want %+v, got %+v", ok, want, max)
		}

		got := intAll(tr)
		wantRange := intRange(treeSize, false)

		if !reflect.DeepEqual(got, wantRange) {
			t.Fatalf("mismatch:\n got: %v\nwant: %v", got, wantRange)
		}

		gotrev := intAllRev(tr)
		wantrev := intRange(treeSize, true)

		if !reflect.DeepEqual(gotrev, wantrev) {
			t.Fatalf("mismatch:\n got: %v\nwant: %v", gotrev, wantrev)
		}

		for _, item := range rand.Perm(treeSize) {
			if x, ok := tr.Delete(item); !ok || x != item {
				t.Fatalf("didn't find %v", item)
			}
		}

		if got = intAll(tr); len(got) > 0 {
			t.Fatalf("some left!: %v", got)
		}

		if got = intAllRev(tr); len(got) > 0 {
			t.Fatalf("some left!: %v", got)
		}
	}
}

func ExampleBTreeG() {
	tr := NewOrderedG[int](*btreeDegree)
	for i := range 10 {
		tr.ReplaceOrInsert(i)
	}

	fmt.Println("len:       ", tr.Len())
	v, ok := tr.Get(3)
	fmt.Println("get3:      ", v, ok)
	v, ok = tr.Get(100)
	fmt.Println("get100:    ", v, ok)
	v, ok = tr.Delete(4)
	fmt.Println("del4:      ", v, ok)
	v, ok = tr.Delete(100)
	fmt.Println("del100:    ", v, ok)
	v, ok = tr.ReplaceOrInsert(5)
	fmt.Println("replace5:  ", v, ok)
	v, ok = tr.ReplaceOrInsert(100)
	fmt.Println("replace100:", v, ok)
	v, ok = tr.Min()
	fmt.Println("min:       ", v, ok)
	v, ok = tr.DeleteMin()
	fmt.Println("delmin:    ", v, ok)
	v, ok = tr.Max()
	fmt.Println("max:       ", v, ok)
	v, ok = tr.DeleteMax()
	fmt.Println("delmax:    ", v, ok)
	fmt.Println("len:       ", tr.Len())
	// Output:
	// len:        10
	// get3:       3 true
	// get100:     0 false
	// del4:       4 true
	// del100:     0 false
	// replace5:   5 true
	// replace100: 0 false
	// min:        0 true
	// delmin:     0 true
	// max:        100 true
	// delmax:     100 true
	// len:        8
}

func TestDeleteMinG(t *testing.T) {
	tr := NewOrderedG[int](3)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int
	for v, ok := tr.DeleteMin(); ok; v, ok = tr.DeleteMin() {
		got = append(got, v)
	}

	if want := intRange(100, false); !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDeleteMaxG(t *testing.T) {
	tr := NewOrderedG[int](3)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int
	for v, ok := tr.DeleteMax(); ok; v, ok = tr.DeleteMax() {
		got = append(got, v)
	}

	if want := intRange(100, true); !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendRangeG(t *testing.T) {
	tr := NewOrderedG[int](2)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.AscendRange(40, 60, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[40:60]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.AscendRange(40, 60, func(a int) bool {
		if a > 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[40:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendRangeG(t *testing.T) {
	tr := NewOrderedG[int](2)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.DescendRange(60, 40, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[39:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.DescendRange(60, 40, func(a int) bool {
		if a < 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendLessThanG(t *testing.T) {
	tr := NewOrderedG[int](*btreeDegree)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.AscendLessThan(60, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[:60]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.AscendLessThan(60, func(a int) bool {
		if a > 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendLessOrEqualG(t *testing.T) {
	tr := NewOrderedG[int](*btreeDegree)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.DescendLessOrEqual(40, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[59:]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.DescendLessOrEqual(60, func(a int) bool {
		if a < 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[39:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendlessorequal:\n got: %v\nwant: %v", got, want)
	}
}

func TestAscendGreaterOrEqualG(t *testing.T) {
	tr := NewOrderedG[int](*btreeDegree)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.AscendGreaterOrEqual(40, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[40:]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.AscendGreaterOrEqual(40, func(a int) bool {
		if a > 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, false)[40:51]; !reflect.DeepEqual(got, want) {
		t.Fatalf("ascendrange:\n got: %v\nwant: %v", got, want)
	}
}

func TestDescendGreaterThanG(t *testing.T) {
	tr := NewOrderedG[int](*btreeDegree)
	for _, v := range rand.Perm(100) {
		tr.ReplaceOrInsert(v)
	}

	var got []int

	tr.DescendGreaterThan(40, func(a int) bool {
		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[:59]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}

	got = got[:0]

	tr.DescendGreaterThan(40, func(a int) bool {
		if a < 50 {
			return false
		}

		got = append(got, a)

		return true
	})

	if want := intRange(100, true)[:50]; !reflect.DeepEqual(got, want) {
		t.Fatalf("descendgreaterthan:\n got: %v\nwant: %v", got, want)
	}
}

const benchmarkTreeSize = 10000

func BenchmarkInsertG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)

	b.StartTimer()

	i := 0
	for i < b.N {
		tr := NewOrderedG[int](*btreeDegree)
		for _, item := range insertP {
			tr.ReplaceOrInsert(item)

			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkSeekG(b *testing.B) {
	b.StopTimer()

	size := 100000
	insertP := rand.Perm(size)
	tr := NewOrderedG[int](*btreeDegree)

	for _, item := range insertP {
		tr.ReplaceOrInsert(item)
	}

	b.StartTimer()

	for i := range b.N {
		tr.AscendGreaterOrEqual(i%size, func(i int) bool { return false })
	}
}

func BenchmarkDeleteInsertG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, item := range insertP {
		tr.ReplaceOrInsert(item)
	}

	b.StartTimer()

	for i := range b.N {
		tr.Delete(insertP[i%benchmarkTreeSize])
		tr.ReplaceOrInsert(insertP[i%benchmarkTreeSize])
	}
}

func BenchmarkDeleteInsertCloneOnceG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, item := range insertP {
		tr.ReplaceOrInsert(item)
	}

	tr = tr.Clone()

	b.StartTimer()

	for i := range b.N {
		tr.Delete(insertP[i%benchmarkTreeSize])
		tr.ReplaceOrInsert(insertP[i%benchmarkTreeSize])
	}
}

func BenchmarkDeleteInsertCloneEachTimeG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, item := range insertP {
		tr.ReplaceOrInsert(item)
	}

	b.StartTimer()

	for i := range b.N {
		tr = tr.Clone()
		tr.Delete(insertP[i%benchmarkTreeSize])
		tr.ReplaceOrInsert(insertP[i%benchmarkTreeSize])
	}
}

func BenchmarkDeleteG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	removeP := rand.Perm(benchmarkTreeSize)

	b.StartTimer()

	i := 0
	for i < b.N {
		b.StopTimer()

		tr := NewOrderedG[int](*btreeDegree)
		for _, v := range insertP {
			tr.ReplaceOrInsert(v)
		}

		b.StartTimer()

		for _, item := range removeP {
			tr.Delete(item)

			i++
			if i >= b.N {
				return
			}
		}

		if tr.Len() > 0 {
			panic(tr.Len())
		}
	}
}

func BenchmarkGetG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	removeP := rand.Perm(benchmarkTreeSize)

	b.StartTimer()

	i := 0
	for i < b.N {
		b.StopTimer()

		tr := NewOrderedG[int](*btreeDegree)
		for _, v := range insertP {
			tr.ReplaceOrInsert(v)
		}

		b.StartTimer()

		for _, item := range removeP {
			tr.Get(item)

			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkGetCloneEachTimeG(b *testing.B) {
	b.StopTimer()

	insertP := rand.Perm(benchmarkTreeSize)
	removeP := rand.Perm(benchmarkTreeSize)

	b.StartTimer()

	i := 0
	for i < b.N {
		b.StopTimer()

		tr := NewOrderedG[int](*btreeDegree)
		for _, v := range insertP {
			tr.ReplaceOrInsert(v)
		}

		b.StartTimer()

		for _, item := range removeP {
			tr = tr.Clone()
			tr.Get(item)

			i++
			if i >= b.N {
				return
			}
		}
	}
}

func BenchmarkAscendG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := 0

		tr.Ascend(func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j++

			return true
		})
	}
}

func BenchmarkDescendG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := len(arr) - 1

		tr.Descend(func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j--

			return true
		})
	}
}

func BenchmarkAscendRangeG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := 100

		tr.AscendRange(100, arr[len(arr)-100], func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j++

			return true
		})

		if j != len(arr)-100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, j)
		}
	}
}

func BenchmarkDescendRangeG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := len(arr) - 100
		tr.DescendRange(arr[len(arr)-100], 100, func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j--

			return true
		})

		if j != 100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, j)
		}
	}
}

func BenchmarkAscendGreaterOrEqualG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := 100
		k := 0

		tr.AscendGreaterOrEqual(100, func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j++
			k++

			return true
		})

		if j != len(arr) {
			b.Fatalf("expected: %v, got %v", len(arr), j)
		}

		if k != len(arr)-100 {
			b.Fatalf("expected: %v, got %v", len(arr)-100, k)
		}
	}
}

func BenchmarkDescendLessOrEqualG(b *testing.B) {
	arr := rand.Perm(benchmarkTreeSize)
	tr := NewOrderedG[int](*btreeDegree)

	for _, v := range arr {
		tr.ReplaceOrInsert(v)
	}

	sort.Ints(arr)
	b.ResetTimer()

	for range b.N {
		j := len(arr) - 100
		k := len(arr)
		tr.DescendLessOrEqual(arr[len(arr)-100], func(item int) bool {
			if item != arr[j] {
				b.Fatalf("mismatch: expected: %v, got %v", arr[j], item)
			}

			j--
			k--

			return true
		})

		if j != -1 {
			b.Fatalf("expected: %v, got %v", -1, j)
		}

		if k != 99 {
			b.Fatalf("expected: %v, got %v", 99, k)
		}
	}
}

const cloneTestSize = 10000

func cloneTestG(t *testing.T, b *BTreeG[int], start int, p []int, wg *sync.WaitGroup, trees *[]*BTreeG[int], lock *sync.Mutex) {
	t.Logf("Starting new clone at %v", start)
	lock.Lock()
	*trees = append(*trees, b)
	lock.Unlock()

	for i := start; i < cloneTestSize; i++ {
		b.ReplaceOrInsert(p[i])

		if i%(cloneTestSize/5) == 0 {
			wg.Add(1)

			go cloneTestG(t, b.Clone(), i+1, p, wg, trees, lock)
		}
	}

	wg.Done()
}

func TestCloneConcurrentOperationsG(t *testing.T) {
	b := NewOrderedG[int](*btreeDegree)
	trees := []*BTreeG[int]{}
	p := rand.Perm(cloneTestSize)

	var wg sync.WaitGroup

	wg.Add(1)

	go cloneTestG(t, b, 0, p, &wg, &trees, &sync.Mutex{})
	wg.Wait()

	want := intRange(cloneTestSize, false)

	t.Logf("Starting equality checks on %d trees", len(trees))

	for i, tree := range trees {
		if !reflect.DeepEqual(want, intAll(tree)) {
			t.Errorf("tree %v mismatch", i)
		}
	}

	t.Log("Removing half from first half")

	toRemove := intRange(cloneTestSize, false)[cloneTestSize/2:]

	for i := range len(trees) / 2 {
		tree := trees[i]

		wg.Add(1)

		go func() {
			for _, item := range toRemove {
				tree.Delete(item)
			}

			wg.Done()
		}()
	}

	wg.Wait()
	t.Log("Checking all values again")

	for i, tree := range trees {
		var wantpart []int
		if i < len(trees)/2 {
			wantpart = want[:cloneTestSize/2]
		} else {
			wantpart = want
		}

		if got := intAll(tree); !reflect.DeepEqual(wantpart, got) {
			t.Errorf("tree %v mismatch, want %v got %v", i, len(want), len(got))
		}
	}
}
