package db

import (
	"context"
	"log"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPGXPool(e *config.Env) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, e.DatabaseURL)
	if err != nil {
		log.Fatalf("could not poolect to the database: %v", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("could not ping the database: %v", err)
		return nil
	}

	return pool
}

type Queries struct {
	*sqlc.Queries
	*pgxpool.Pool
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{
		Queries: sqlc.New(pool),
		Pool:    pool,
	}
}

func (q *Queries) UseTx(
	ctx context.Context,
) *Queries {
	t, ok := ctx.Value(tx.Key).(pgx.Tx)
	if ok {
		return &Queries{
			Queries: q.WithTx(t),
		}
	}
	return q
}
