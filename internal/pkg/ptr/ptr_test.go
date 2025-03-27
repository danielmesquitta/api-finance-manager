package ptr

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {
	type args[T any] struct {
		ptr           *T
		defaultValues []T
	}

	type Test[T any] struct {
		name string
		args args[T]
		want T
	}

	intTests := []Test[int]{
		{
			name: "nil pointer returns first non-zero default",
			args: args[int]{
				ptr:           nil,
				defaultValues: []int{0, 17, 5},
			},
			want: 17,
		},
		{
			name: "nil pointer returns latest default value when all defaults are zero",
			args: args[int]{
				ptr:           nil,
				defaultValues: []int{0},
			},
			want: 0,
		},
		{
			name: "non-nil pointer returns its value",
			args: args[int]{
				ptr:           New(10),
				defaultValues: []int{0, 17, 5},
			},
			want: 10,
		},
		{
			name: "non-nil pointer returns its value ignoring defaults",
			args: args[int]{
				ptr:           New(42),
				defaultValues: []int{100, 200},
			},
			want: 42,
		},
		{
			name: "nil pointer and no defaults returns zero value",
			args: args[int]{
				ptr:           nil,
				defaultValues: []int{},
			},
			want: 0,
		},
		{
			name: "non-nil pointer returns its value",
			args: args[int]{
				ptr:           New(42),
				defaultValues: []int{},
			},
			want: 42,
		},
	}

	for _, tt := range intTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Value(tt.args.ptr, tt.args.defaultValues...); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}

	stringTests := []Test[string]{
		{
			name: "nil pointer returns first non-zero default string",
			args: args[string]{
				ptr:           nil,
				defaultValues: []string{"", "hello", "world"},
			},
			want: "hello",
		},
		{
			name: "non-nil pointer returns its value string",
			args: args[string]{
				ptr:           New("go"),
				defaultValues: []string{"", "hello", "world"},
			},
			want: "go",
		},
		{
			name: "nil pointer and no defaults returns empty string",
			args: args[string]{
				ptr:           nil,
				defaultValues: []string{},
			},
			want: "",
		},
	}

	for _, tt := range stringTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Value(tt.args.ptr, tt.args.defaultValues...); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}

	mapTests := []Test[map[string]int]{
		{
			name: "non-nil pointer returns its map",
			args: args[map[string]int]{
				ptr: New(map[string]int{"a": 1, "b": 2}),
				defaultValues: []map[string]int{
					{"a": 100},
				},
			},
			want: map[string]int{"a": 1, "b": 2},
		},
		{
			name: "nil pointer returns first non-nil default map",
			args: args[map[string]int]{
				ptr: nil,
				defaultValues: []map[string]int{
					nil,
					{"x": 10},
				},
			},
			want: map[string]int{"x": 10},
		},
		{
			name: "nil pointer and no defaults returns nil map",
			args: args[map[string]int]{
				ptr:           nil,
				defaultValues: []map[string]int{},
			},
			want: nil,
		},
		{
			name: "nil pointer with all nil defaults returns nil",
			args: args[map[string]int]{
				ptr: nil,
				defaultValues: []map[string]int{
					nil,
				},
			},
			want: nil,
		},
	}

	for _, tt := range mapTests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Value(tt.args.ptr, tt.args.defaultValues...); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
