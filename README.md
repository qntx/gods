# Go Data Structures

A collection of data structures implemented in Go.

## Containers

| **Data** | **Structure**   | **Ordered** | **Iterator** | **Serializable** | **Referenced by** | **Implemented** |
| -------- | --------------- | ----------- | ------------ | ---------------- | ----------------- | --------------- |
| Tree     |                 |             |              |                  |                   |                 |
|          | `BTree`         | Y           | Y            | Y                | Key               | √               |
|          | `RBTree`        | Y           | Y            | Y                | Key               | √               |
|          | `AVLTree`       | Y           | Y            | Y                | Key               | √               |
| Map      |                 |             |              |                  |                   |                 |
|          | `HashMap`       | Y (sort)    | Y            | Y                | Key               | √               |
|          | `HashBiMap`     | Y (sort)    | Y            | Y                | Key               | √               |
|          | `BTreeBiMap`    | Y           | Y            | Y                | Key               | √               |
|          | `RBTreeBiMap`   | Y           | Y            | Y                | Key               | √               |
|          | `LinkedHashMap` | Y           | Y            | Y                | Key               | √               |
| Set      |                 |             |              |                  |                   |                 |
|          | `HashSet`       | N           | Y            | Y                | Index             | √               |
|          | `BTreeSet`      | Y           | Y            | Y                | Index             | √               |
|          | `RBTreeSet`     | Y           | Y            | Y                | Index             | √               |
|          | `LinkedHashSet` | Y           | Y            | Y                | Index             | √               |
| Queue    |                 |             |              |                  |                   |                 |
|          | `SliceDeque`    | Y           | Y            | Y                | Index             | √               |
|          | `PriorityQueue` | Y           | Y            | N                | Index             | √               |
| Stack    |                 |             |              |                  |                   |                 |
|          | `SliceStack`    | Y           | Y            | Y                | Index             | √               |

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
