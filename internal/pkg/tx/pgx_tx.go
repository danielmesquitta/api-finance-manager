package tx

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxTX struct {
	db *pgxpool.Pool
	l  *slog.Logger
}

func NewPgxTX(
	db *pgxpool.Pool,
	l *slog.Logger,
) *PgxTX {
	return &PgxTX{
		db: db,
	}
}

func (t *PgxTX) Do(
	ctx context.Context,
	fn func(context.Context) error,
) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			t.l.Error("failed to rollback transaction", "error", err)
		}
	}()

	ctx = context.WithValue(ctx, Key, tx)
	if err := fn(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
