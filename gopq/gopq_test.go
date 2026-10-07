package gopq

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestOrdering(t *testing.T) {
	for _, tc := range []struct {
		name   string
		values []int
	}{
		{name: "empty"},
		{name: "singleton", values: []int{0}},
		{name: "two", values: []int{2, 1}},
		{name: "ascending", values: []int{1, 2, 3, 4, 5}},
		{name: "descending", values: []int{5, 4, 3, 2, 1}},
		{name: "duplicates and negatives", values: []int{0, -3, 2, -3, 0, 7}},
	} {
		for _, maxHeap := range []bool{false, true} {
			name := tc.name + "/min"
			if maxHeap {
				name = tc.name + "/max"
			}
			t.Run(name, func(t *testing.T) {
				less := func(a, b int) bool { return a < b }
				if maxHeap {
					less = func(a, b int) bool { return a > b }
				}
				want := append([]int(nil), tc.values...)
				sort.Slice(want, func(i, j int) bool { return less(want[i], want[j]) })
				q := New(less, tc.values...)
				for i, expected := range want {
					if q.Len() != len(want)-i {
						t.Fatalf("Len() = %d, want %d", q.Len(), len(want)-i)
					}
					if got, ok := q.Peek(); !ok || got != expected {
						t.Fatalf("Peek() = %d, %v; want %d, true", got, ok, expected)
					}
					if got, ok := q.Pop(); !ok || got != expected {
						t.Fatalf("Pop() = %d, %v; want %d, true", got, ok, expected)
					}
				}
				if got, ok := q.Pop(); ok || got != 0 || q.Len() != 0 {
					t.Fatalf("empty Pop() = %d, %v; Len() = %d", got, ok, q.Len())
				}
				if got, ok := q.Peek(); ok || got != 0 {
					t.Fatalf("empty Peek() = %d, %v", got, ok)
				}
				q.Push(42)
				if got, ok := q.Pop(); !ok || got != 42 {
					t.Fatalf("reuse Pop() = %d, %v", got, ok)
				}
			})
		}
	}
}

func TestMixedOperations(t *testing.T) {
	for _, maxHeap := range []bool{false, true} {
		less := func(a, b int) bool { return a < b }
		if maxHeap {
			less = func(a, b int) bool { return a > b }
		}
		rng := rand.New(rand.NewSource(42))
		q := New(less)
		var model []int
		for step := 0; step < 10000; step++ {
			if len(model) == 0 || rng.Intn(2) == 0 {
				value := rng.Intn(101) - 50
				q.Push(value)
				model = append(model, value)
				sort.Slice(model, func(i, j int) bool { return less(model[i], model[j]) })
			} else {
				got, ok := q.Pop()
				if !ok || got != model[0] {
					t.Fatalf("step %d: Pop() = %d, %v; want %d", step, got, ok, model[0])
				}
				model = model[1:]
			}
			if q.Len() != len(model) {
				t.Fatalf("step %d: Len() = %d, want %d", step, q.Len(), len(model))
			}
			got, ok := q.Peek()
			if ok != (len(model) > 0) || (ok && got != model[0]) {
				t.Fatalf("step %d: Peek() = %d, %v; model = %v", step, got, ok, model)
			}
		}
		for _, expected := range model {
			if got, ok := q.Pop(); !ok || got != expected {
				t.Fatalf("drain Pop() = %d, %v; want %d", got, ok, expected)
			}
		}
	}
}

func TestAnyType(t *testing.T) {
	strings := New(func(a, b string) bool { return a < b }, "pear", "apple", "banana")
	for _, want := range []string{"apple", "banana", "pear"} {
		if got, ok := strings.Pop(); !ok || got != want {
			t.Fatalf("string Pop() = %q, %v; want %q", got, ok, want)
		}
	}
	// A struct containing a slice is not comparable; T need not be comparable.
	type node struct {
		cost int
		path []int
	}
	want := node{cost: 1, path: []int{0, 1}}
	nodes := New(func(a, b node) bool { return a.cost < b.cost },
		node{cost: 3, path: []int{0, 2}}, want)
	if got, ok := nodes.Pop(); !ok || !reflect.DeepEqual(got, want) {
		t.Fatalf("node Pop() = %v, %v; want %v", got, ok, want)
	}
}

func TestCopiesInitialSlice(t *testing.T) {
	values := []int{3, 1, 2}
	q := New(func(a, b int) bool { return a < b }, values...)
	if !reflect.DeepEqual(values, []int{3, 1, 2}) {
		t.Fatalf("New changed input: %v", values)
	}
	values[1] = -100
	if got, _ := q.Pop(); got != 1 {
		t.Fatalf("input mutation changed queue: Pop() = %d", got)
	}
}

func TestPopReleasesReferences(t *testing.T) {
	a, b := 1, 2
	q := New(func(a, b *int) bool { return *a < *b }, &a, &b)
	backing := q.items
	for q.Len() > 0 {
		q.Pop()
		if backing[q.Len()] != nil {
			t.Fatal("Pop retained a reference beyond the slice length")
		}
	}
}

func TestNilComparatorPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("New with nil comparator did not panic")
		}
	}()
	New[int](nil)
}

func TestZeroValuePushPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("zero-value Push did not panic")
		}
	}()
	var q Queue[int]
	q.Push(1)
}
