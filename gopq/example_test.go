package gopq_test

import (
	"fmt"

	"github.com/savil/gopq"
)

func ExampleNew() {
	q := gopq.New(func(a, b int) bool { return a < b }, 5, 1, 3)
	q.Push(2)
	for q.Len() > 0 {
		value, _ := q.Pop()
		fmt.Println(value)
	}
	// Output:
	// 1
	// 2
	// 3
	// 5
}

func ExampleNew_maxHeap() {
	q := gopq.New(func(a, b int) bool { return a > b }, 5, 1, 3)
	value, ok := q.Pop()
	fmt.Println(value, ok)
	// Output: 5 true
}

func ExampleNew_strings() {
	q := gopq.New(func(a, b string) bool { return a < b }, "pear", "apple")
	value, ok := q.Pop()
	fmt.Println(value, ok)
	// Output: apple true
}

func ExampleNew_structs() {
	type Node struct {
		Name     string
		Distance int
	}
	q := gopq.New(func(a, b Node) bool {
		if a.Distance != b.Distance {
			return a.Distance < b.Distance
		}
		return a.Name < b.Name // Optional deterministic tie-breaker.
	})
	q.Push(Node{Name: "B", Distance: 9})
	q.Push(Node{Name: "A", Distance: 2})
	value, ok := q.Pop()
	fmt.Println(value.Name, value.Distance, ok)
	// Output: A 2 true
}
