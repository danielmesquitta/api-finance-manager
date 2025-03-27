package ptr_test

import (
	"fmt"

	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
)

func ExampleCoalesce() {
	var np *int
	fmt.Println(ptr.Coalesce(np, 1))
	np = new(int)
	fmt.Println(ptr.Coalesce(np, 1))
	// Output:
	// 1
	// 0
}

func ExampleDeref() {
	var np *int
	fmt.Println(ptr.Deref(np))
	one := 1
	np = &one
	fmt.Println(ptr.Deref(np))
	// Output:
	// 0
	// 1
}
