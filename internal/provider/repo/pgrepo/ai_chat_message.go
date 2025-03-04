package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type AIChatMessageRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewAIChatMessageRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *AIChatMessageRepo {
	return &AIChatMessageRepo{
		db: db,
		qb: qb,
	}
}

func (r *AIChatMessageRepo) ListAIChatMessages(
	ctx context.Context,
	opts ...repo.AIChatMessageOption,
) ([]entity.AIChatMessage, error) {
	aiChats, err := r.qb.ListAIChatMessages(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return aiChats, nil
}

func (r *AIChatMessageRepo) CountAIChatMessages(
	ctx context.Context,
	opts ...repo.AIChatMessageOption,
) (int64, error) {
	return r.qb.CountAIChatMessages(ctx, opts...)
}

var _ repo.AIChatMessageRepo = (*AIChatMessageRepo)(nil)
