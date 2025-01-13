package tx

import (
	"context"
)

type key byte

var Key key

type TX interface {
	// Do is not concurrent safe because all queries in a transaction must use the same connection.
	// Since each query blocks the connection, transactions enforce sequential execution.
	// Running queries concurrently with goroutines in a transaction will likely cause errors.
	Do(ctx context.Context, fn func(context.Context) error) error
}
