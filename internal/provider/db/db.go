package db

import (
	"context"
	"log"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/jackc/pgx/v5"
)

func NewPGXConn(e *config.Env) *pgx.Conn {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, e.DatabaseURL)
	if err != nil {
		log.Fatalf("could not connect to the database: %v", err)
		return nil
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatalf("could not ping the database: %v", err)
		return nil
	}

	return conn
}

type Queries struct {
	*sqlc.Queries
	*pgx.Conn
}

func NewQueries(conn *pgx.Conn) *Queries {
	return &Queries{
		Queries: sqlc.New(conn),
		Conn:    conn,
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
