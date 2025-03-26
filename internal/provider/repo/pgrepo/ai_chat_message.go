package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AIChatMessageRepo struct {
	db *db.DB
}

func NewAIChatMessageRepo(
	db *db.DB,
) *AIChatMessageRepo {
	return &AIChatMessageRepo{
		db: db,
	}
}

func (r *AIChatMessageRepo) CreateAIChatMessage(
	ctx context.Context,
	params repo.CreateAIChatMessageParams,
) (*entity.AIChatMessage, error) {
	dbParams := sqlc.CreateAIChatMessageParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	aiChat, err := tx.CreateAIChatMessage(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	var result entity.AIChatMessage
	if err := copier.Copy(&result, aiChat); err != nil {
		return nil, errs.New(err)
	}

	return &result, nil
}

func (r *AIChatMessageRepo) DeleteAIChatMessages(
	ctx context.Context,
	aiChatID uuid.UUID,
) error {
	return r.db.DeleteAIChatMessages(ctx, aiChatID)
}

var _ repo.AIChatMessageRepo = (*AIChatMessageRepo)(nil)
