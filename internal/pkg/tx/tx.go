package tx

import (
	"context"
)

type key byte

var Key key

type Tx interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
