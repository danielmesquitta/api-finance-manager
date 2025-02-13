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

	pool, err := pgxpool.New(ctx, e.PostgresDatabaseURL)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
		return nil
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("could not ping the database: %v", err)
		return nil
	}

	return pool
}

type DB struct {
	*sqlc.Queries
	*pgxpool.Pool
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		Queries: sqlc.New(pool),
		Pool:    pool,
	}
}

func (db *DB) UseTx(
	ctx context.Context,
) *DB {
	t, ok := ctx.Value(tx.Key).(pgx.Tx)
	if ok {
		return &DB{
			Queries: db.WithTx(t),
			Pool:    db.Pool,
		}
	}
	return db
}
