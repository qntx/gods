# Go Data Structures

A collection of data structures implemented in Go.

## Containers

| **Data** | **Structure**                         | **Ordered** | **Iterator** | **Enumerable** | **Referenced by** |
| :--- |:--------------------------------------| :---: | :---: | :---: | :---: |
| Set ||||||
|   | `HashSet`                 | no | no | no | index |
|   | `BTreeSet`                 | yes | yes | yes | index |
|   | `RBTreeSet`                | yes | yes | yes | index |
| Map ||||||
|   | `HashMap`                 | no | no | no | key |
|   | `BTreeMap`                | yes | yes | yes | key |
|   | `RBTreeMap`                | yes | yes | yes | key |
| Tree ||||||
|   | `BTree`         | yes | yes | no | key |
|   | `RBTree`                    | yes | yes | no | key |
| Queue ||||||
|   | `ArrayDeque`           | yes | yes | no | index |
|   | `PriorityQueue`     | yes | yes | no | index |

## License

MIT

## Acknowledgments

- [google/btree](https://github.com/google/btree/tree/master)
- [emirpasic/gods](https://github.com/emirpasic/gods)
- [gammazero/deque](https://github.com/gammazero/deque)
- [dnaeon/go-priorityqueue](https://github.com/dnaeon/go-priorityqueue)
- [deckarep/golang-set](https://github.com/deckarep/golang-set)
- [huandu/skiplist](https://github.com/huandu/skiplist)
