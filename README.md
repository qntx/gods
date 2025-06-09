# Go Data Structures

A collection of data structures implemented in Go.

## Containers

| **Data** | **Structure**   | **Ordered** | **Iterator** | **Enumerable** | **Referenced by** | **Implemented** |
| -------- | --------------- | ----------- | ------------ | -------------- | ----------------- | --------------- |
| Set      |                 |             |              |                |                   |                 |
|          | `HashSet`       | N           | N            | N              | Index             | Y               |
|          | `BTreeSet`      | Y           | Y            | Y              | Index             | N               |
|          | `RBTreeSet`     | Y           | Y            | Y              | Index             | Y               |
|          | `LinkedHashSet` | Y           | Y            | Y              | Index             | Y               |
| Map      |                 |             |              |                |                   |                 |
|          | `HashMap`       | N           | N            | N              | Key               | Y               |
|          | `HashBidiMap`   | N           | N            | N              | Key               | Y               |
|          | `BTreeBidiMap`  | Y           | Y            | Y              | Key               | N               |
|          | `RBTreeBidiMap` | Y           | Y            | Y              | Key               | Y               |
|          | `LinkedHashMap` | Y           | Y            | Y              | Key               | Y               |
| Tree     |                 |             |              |                |                   |                 |
|          | `AVLTree`       | Y           | Y            | N              | Key               | Y               |
|          | `BTree`         | Y           | Y            | N              | Key               | N               |
|          | `RBTree`        | Y           | Y            | N              | Key               | Y               |
| Queue    |                 |             |              |                |                   |                 |
|          | `SliceDeque`    | Y           | Y            | N              | Index             | Y               |
|          | `PriorityQueue` | Y           | Y            | N              | Index             | Y               |

## Benchmarks

```shell
go test -run=NO_TEST -bench . -benchmem  -benchtime 1s ./...
```

## License

MIT

## Acknowledgments

- [google/btree](https://github.com/google/btree/tree/master)
- [emirpasic/gods](https://github.com/emirpasic/gods)
- [gammazero/deque](https://github.com/gammazero/deque)
- [dnaeon/go-priorityqueue](https://github.com/dnaeon/go-priorityqueue)
- [deckarep/golang-set](https://github.com/deckarep/golang-set)
- [huandu/skiplist](https://github.com/huandu/skiplist)
