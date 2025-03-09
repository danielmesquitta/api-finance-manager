package query

import "github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"

type TestQueryBuilder struct {
	*query.QueryBuilder
}

func NewTestQueryBuilder(
	qb *query.QueryBuilder,
) *TestQueryBuilder {
	return &TestQueryBuilder{
		QueryBuilder: qb,
	}
}
