package ptr_test

import (
	"fmt"

	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
)

func ExampleFirst() {
	type config struct{ string }
	userInput := func() *config {
		return nil
	}
	someConfig := ptr.First(
		userInput(),
		&config{"default config"},
	)
	fmt.Println(someConfig)
	// Output:
	// &{default config}
}
