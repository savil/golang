# golang

Reusable Go data structures for coding interviews and small projects. Requires
Go 1.18 or later and has no external dependencies.

## Packages

| Package | Purpose | Import path |
| --- | --- | --- |
| [gopq](gopq/README.md) | Generic min/max priority queue | `github.com/savil/golang/gopq` |

## Install

Run in your Go project:

```sh
go get github.com/savil/golang/gopq
```

## Usage

Create a priority queue with a comparison function. This example uses a min-heap
of integers:

```go
package main

import (
    "fmt"

    "github.com/savil/golang/gopq"
)

func main() {
    q := gopq.New(func(a, b int) bool { return a < b }, 5, 1, 3)
    q.Push(2)

    fmt.Println(q.Pop())  // 1 true
    fmt.Println(q.Peek()) // 2 true
    fmt.Println(q.Len())  // 3
}
```

The comparison returns true when `a` should be popped before `b`. Use `a > b`
for a max-heap. The same queue implementation supports strings, structs,
pointers, and other types through a comparison function for that type.

`Pop` and `Peek` return `(value, ok)`; `ok` is false when the queue is empty.
Push and pop take O(log n) time, with amortized slice growth for push; peek and
length take O(1). Equal priorities have no guaranteed FIFO order. Keep queued
priorities unchanged, including fields accessed through pointers.

See the [gopq documentation](gopq/README.md) for more examples, the full API,
and ordering rules.

## License

MIT licensed; see [LICENSE](LICENSE).
