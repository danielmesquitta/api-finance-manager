package query

import (
	"math"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
)

type QueryBuilder struct {
	e  *config.Env
	db *db.DB
}

func NewQueryBuilder(
	e *config.Env,
	db *db.DB,
) *QueryBuilder {
	return &QueryBuilder{
		e:  e,
		db: db,
	}
}

func (qb *QueryBuilder) calculateMaxLevenshteinDistance(search string) int {
	searchLength := float64(len(search))
	maxLevenshteinDistance := int(
		math.Floor(qb.e.MaxLevenshteinDistancePercentage * searchLength),
	)
	return maxLevenshteinDistance
}
