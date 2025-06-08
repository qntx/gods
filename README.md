# Go Data Structures

A collection of data structures implemented in Go.

## Containers

| **Data** | **Structure**   | **Ordered** | **Iterator** | **Enumerable** | **Referenced by** | **Implemented** |
| -------- | --------------- | ----------- | ------------ | -------------- | ----------------- | --------------- |
| Set      |                 |             |              |                |                   |                 |
|          | `HashSet`       | N           | N            | N              | Index             | Y               |
|          | `BTreeSet`      | Y           | Y            | Y              | Index             | N               |
|          | `RBTreeSet`     | Y           | Y            | Y              | Index             | N               |
| Map      |                 |             |              |                |                   |                 |
|          | `HashMap`       | N           | N            | N              | Key               | N               |
|          | `HashBidiMap`   | N           | N            | N              | Key               | N               |
|          | `BTreeMap`      | Y           | Y            | Y              | Key               | N               |
|          | `RBTreeMap`     | Y           | Y            | Y              | Key               | N               |
|          | `RBTreeBidiMap` | Y           | Y            | Y              | Key               | N               |
| Tree     |                 |             |              |                |                   |                 |
|          | `BTree`         | Y           | Y            | N              | Key               | N               |
|          | `RBTree`        | Y           | Y            | N              | Key               | N               |
| Queue    |                 |             |              |                |                   |                 |
|          | `ArrayDeque`    | Y           | Y            | N              | Index             | N               |
|          | `PriorityQueue` | Y           | Y            | N              | Index             | N               |

## License

MIT

## Acknowledgments

- [google/btree](https://github.com/google/btree/tree/master)
- [emirpasic/gods](https://github.com/emirpasic/gods)
- [gammazero/deque](https://github.com/gammazero/deque)
- [dnaeon/go-priorityqueue](https://github.com/dnaeon/go-priorityqueue)
- [deckarep/golang-set](https://github.com/deckarep/golang-set)
- [huandu/skiplist](https://github.com/huandu/skiplist)
