package ptr

import (
	"reflect"
)

// New returns a pointer to the value passed as argument.
func New[T any](value T) *T {
	return &value
}

// Value returns the value pointed to by ptr. If ptr is nil, it returns the first non-zero default value,
// or the zero value if none exist.
func Value[T any](ptr *T, defaultValues ...T) T {
	if ptr != nil {
		return *ptr
	}
	return or(defaultValues...)
}

// or returns the first of its arguments that is not the zero value.
// If no argument is non-zero, it returns the zero value.
func or[T any](vals ...T) T {
	for _, val := range vals {
		if !reflect.ValueOf(val).IsZero() {
			return val
		}
	}
	if len(vals) > 0 {
		return vals[len(vals)-1]
	}
	var zero T
	return zero
}
