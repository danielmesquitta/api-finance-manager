package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/query"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AIChatRepo struct {
	db *db.DB
	qb *query.QueryBuilder
}

func NewAIChatRepo(
	db *db.DB,
	qb *query.QueryBuilder,
) *AIChatRepo {
	return &AIChatRepo{
		db: db,
		qb: qb,
	}
}

func (r *AIChatRepo) CreateAIChat(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.AIChat, error) {
	tx := r.db.UseTx(ctx)
	aiChat, err := tx.CreateAIChat(ctx, userID)
	if err != nil {
		return nil, errs.New(err)
	}

	var result entity.AIChat
	if err := copier.Copy(&result, aiChat); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *AIChatRepo) ListAIChats(
	ctx context.Context,
	opts ...repo.AIChatOption,
) ([]entity.AIChat, error) {
	aiChats, err := r.qb.ListAIChats(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return aiChats, nil
}

func (r *AIChatRepo) CountAIChats(
	ctx context.Context,
	opts ...repo.AIChatOption,
) (int64, error) {
	return r.qb.CountAIChats(ctx, opts...)
}

func (r *AIChatRepo) GetLatestAIChatByUserID(
	ctx context.Context,
	userID uuid.UUID,
) (*entity.FullAIChat, error) {
	aiChat, err := r.db.GetLatestAIChatByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := entity.FullAIChat{}
	if err := copier.Copy(&result, aiChat); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *AIChatRepo) DeleteAIChat(ctx context.Context, id uuid.UUID) error {
	return r.db.DeleteAIChat(ctx, id)
}

func (r *AIChatRepo) GetAIChat(
	ctx context.Context,
	id uuid.UUID,
) (*entity.AIChat, error) {
	aiChat, err := r.db.GetAIChat(ctx, id)
	if err != nil {
		return nil, errs.New(err)
	}

	var result entity.AIChat
	if err := copier.Copy(&result, aiChat); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *AIChatRepo) UpdateAIChat(
	ctx context.Context,
	params repo.UpdateAIChatParams,
) error {
	dbParams := sqlc.UpdateAIChatParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	if err := r.db.UpdateAIChat(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.AIChatRepo = (*AIChatRepo)(nil)
