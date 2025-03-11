package db

import (
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
)

type TestDB struct {
	*db.DB
}

func NewTestDB(
	db *db.DB,
) *TestDB {
	return &TestDB{
		DB: db,
	}
}
