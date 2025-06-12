# Go Data Structures

A collection of data structures implemented in Go.

## Containers

| **Data** | **Structure**   | **Ordered** | **Iterator** | **Referenced by** | **Implemented** |
| -------- | --------------- | ----------- | ------------ | ----------------- | --------------- |
| Tree     |                 |             |              |                   |                 |
|          | `BTree`         | Y           | Y            | Key               | Y               |
|          | `RBTree`        | Y           | Y            | Key               | Y               |
|          | `AVLTree`       | Y           | Y            | Key               | Y               |
| Map      |                 |             |              |                   |                 |
|          | `HashMap`       | N           | Y            | Key               | Y               |
|          | `HashBiMap`     | N           | Y            | Key               | Y               |
|          | `BTreeBiMap`    | Y           | Y            | Key               | Y               |
|          | `RBTreeBiMap`   | Y           | Y            | Key               | Y               |
|          | `LinkedHashMap` | Y           | Y            | Key               | Y               |
| Set      |                 |             |              |                   |                 |
|          | `HashSet`       | N           | Y            | Index             | Y               |
|          | `BTreeSet`      | Y           | Y            | Index             | Y               |
|          | `RBTreeSet`     | Y           | Y            | Index             | Y               |
|          | `LinkedHashSet` | Y           | Y            | Index             | Y               |
| Queue    |                 |             |              |                   |                 |
|          | `SliceDeque`    | Y           | Y            | Index             | Y               |
|          | `PriorityQueue` | Y           | Y            | Index             | Y               |
| Stack    |                 |             |              |                   |                 |
|          | `SliceStack`    | Y           | Y            | Index             | N               |

## Benchmarks

```shell
go test -run=NO_TEST -bench . -benchmem  -benchtime 1s ./...
```

## License

MIT

## Acknowledgments

- [emirpasic/gods](https://github.com/emirpasic/gods)
- [google/btree](https://github.com/google/btree/tree/master)
- [huandu/skiplist](https://github.com/huandu/skiplist)
- [deckarep/golang-set](https://github.com/deckarep/golang-set)
- [gammazero/deque](https://github.com/gammazero/deque)
- [dnaeon/go-priorityqueue](https://github.com/dnaeon/go-priorityqueue)
