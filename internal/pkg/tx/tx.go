package tx

import (
	"context"
)

type key byte

var Key key

type TX interface {
	Do(ctx context.Context, fn func(context.Context) error) error
}
