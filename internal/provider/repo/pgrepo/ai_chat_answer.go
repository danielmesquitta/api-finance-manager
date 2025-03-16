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

type AIChatAnswerRepo struct {
	db *db.DB
}

func NewAIChatAnswerRepo(
	db *db.DB,
) *AIChatAnswerRepo {
	return &AIChatAnswerRepo{
		db: db,
	}
}

func (r *AIChatAnswerRepo) CreateAIChatAnswer(
	ctx context.Context,
	params repo.CreateAIChatAnswerParams,
) (*entity.AIChatAnswer, error) {
	dbParams := sqlc.CreateAIChatAnswerParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	aiChatAnswer, err := tx.CreateAIChatAnswer(ctx, dbParams)
	if err != nil {
		return nil, errs.New(err)
	}

	response := entity.AIChatAnswer{}
	if err := copier.Copy(&response, aiChatAnswer); err != nil {
		return nil, errs.New(err)
	}

	return &response, nil
}

func (r *AIChatAnswerRepo) DeleteAIChatAnswers(
	ctx context.Context,
	aiChatID uuid.UUID,
) error {
	return r.db.DeleteAIChatAnswers(ctx, aiChatID)
}

func (r *AIChatAnswerRepo) UpdateAIChatAnswer(
	ctx context.Context,
	params repo.UpdateAIChatAnswerParams,
) error {
	dbParams := sqlc.UpdateAIChatAnswerParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	if err := r.db.UpdateAIChatAnswer(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.AIChatAnswerRepo = (*AIChatAnswerRepo)(nil)
