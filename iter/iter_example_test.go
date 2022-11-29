package iter_test

import (
	"fmt"

	"go.acim.net/aim/iter"
)

var _ iter.Iterator[any] = (*iter.Slice[any])(nil)

func ExampleIterator_Next() {
	items := iter.FromSlice([]int{1, 2, 3})

	for items.HasNext() {
		fmt.Println(*items.Next())
	}

	// Output:
	// 1
	// 2
	// 3
}

func ExampleIterator_Chan() {
	items := iter.FromSlice([]string{"a", "b", "c"})

	for item := range items.Chan() {
		fmt.Println(*item)
	}

	// Output:
	// a
	// b
	// c
}
