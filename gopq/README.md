# gopq

A small, type-safe priority queue for Go interviews. Works with integers,
strings, structs, pointers, and any other Go type you can supply an ordering for.
Go 1.18+, no external dependencies.

Each queue has one element type. For example, `Queue[int]` stores integers and
`Queue[Node]` stores nodes; the same implementation supports both.

## Install

Run in your Go project:

```sh
go get github.com/savil/golang/gopq
```

Confirm your interview environment allows external packages. If it does not, copy
`gopq.go` into your project and change its package declaration to match.

## Integers: min-heap

```go
package main

import (
    "fmt"
    "github.com/savil/golang/gopq"
)

func main() {
    q := gopq.New(func(a, b int) bool { return a < b }, 5, 1, 3)
    q.Push(2)

    for q.Len() > 0 {
        value, _ := q.Pop()
        fmt.Println(value) // 1, 2, 3, 5
    }
}
```

`less(a, b)` means **a should come out before b**. Reverse the comparison for a
max-heap:

```go
q := gopq.New(func(a, b int) bool { return a > b })
q.Push(10)
q.Push(20)
value, ok := q.Pop() // 20, true
_, _ = value, ok
```

## Strings

```go
q := gopq.New(func(a, b string) bool { return a < b }, "pear", "apple")
word, ok := q.Pop() // "apple", true (lexicographic ordering)
_, _ = word, ok
```

## Structs

```go
type Node struct {
    ID       int
    Distance int
}

q := gopq.New(func(a, b Node) bool {
    return a.Distance < b.Distance
})
q.Push(Node{ID: 1, Distance: 10})
q.Push(Node{ID: 2, Distance: 3})
node, ok := q.Pop() // Node{ID: 2, Distance: 3}, true
_, _ = node, ok
```

For pointers, use `func(a, b *Node) bool`. Define how nil pointers compare if
you plan to enqueue them. Structs may contain slices or maps: `T` has no
`comparable` constraint.

## API and complexity

| Operation | Behavior | Time |
| --- | --- | --- |
| `New(less, values...)` | Build from optional values; copies the input slice | O(n) |
| `Push(value)` | Add a value | O(log n), amortized |
| `Pop()` | Remove highest priority; return `(value, ok)` | O(log n) |
| `Peek()` | Read highest priority; return `(value, ok)` | O(1) |
| `Len()` | Number of queued values | O(1) |

Uses O(n) storage. An empty `Pop` or `Peek` returns the zero value and `false`,
so a queued zero is distinguishable from an empty queue. Always construct with
`New`; a nil comparison function panics.

The comparison must be consistent and transitive, and return false for equal
values. Equal priorities are not FIFO; add a sequence number to the comparison
if you need insertion-order tie-breaking. Handle NaN explicitly when comparing
floating-point values that may contain NaN.

Do not mutate the priorities of queued pointers, or change ordering through
captured variables. For Dijkstra's algorithm, push a new value when a distance
improves and skip stale entries when popping. The queue does not expose
decrease-key or arbitrary removal. It is not safe for concurrent use and should
not be copied after use.

## License

MIT licensed; see [LICENSE](../LICENSE).
