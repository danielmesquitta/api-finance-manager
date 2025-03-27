package ptr_test

import (
	"fmt"

	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
)

func ExampleNew() {
	strptr1 := ptr.New("meaning of life")
	strptr2 := ptr.New("meaning of life")
	fmt.Println(strptr1 != strptr2)
	fmt.Println(*strptr1 == *strptr2)

	intp1 := ptr.New(42)
	intp2 := ptr.New(42)
	fmt.Println(intp1 != intp2)
	fmt.Println(*intp1 == *intp2)

	type MyFloat float64
	fp := ptr.New[MyFloat](42)
	fmt.Println(fp != nil)
	fmt.Println(*fp == 42)

	// Output:
	// true
	// true
	// true
	// true
	// true
	// true
}
