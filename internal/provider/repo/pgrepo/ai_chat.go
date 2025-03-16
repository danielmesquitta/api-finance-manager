package pgrepo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type AIChatRepo struct {
	db *db.DB
}

func NewAIChatRepo(
	db *db.DB,
) *AIChatRepo {
	return &AIChatRepo{
		db: db,
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
	aiChats, err := r.db.ListAIChats(ctx, opts...)
	if err != nil {
		return nil, errs.New(err)
	}

	return aiChats, nil
}

func (r *AIChatRepo) CountAIChats(
	ctx context.Context,
	opts ...repo.AIChatOption,
) (int64, error) {
	return r.db.CountAIChats(ctx, opts...)
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

func (r *AIChatRepo) ListAIChatMessagesAndAnswers(
	ctx context.Context,
	params repo.ListAIChatMessagesAndAnswersParams,
) ([]entity.AIChatMessageAndAnswer, error) {
	dbParams := sqlc.ListAIChatMessagesAndAnswersParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return nil, errs.New(err)
	}

	messages, err := r.db.ListAIChatMessagesAndAnswers(ctx, dbParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errs.New(err)
	}

	result := []entity.AIChatMessageAndAnswer{}
	if err := copier.Copy(&result, messages); err != nil {
		return nil, errs.New(err)
	}

	return result, nil
}

func (a *AIChatRepo) CountAIChatMessagesAndAnswers(
	ctx context.Context,
	aiChatID uuid.UUID,
) (int64, error) {
	count, err := a.db.CountAIChatMessagesAndAnswers(ctx, aiChatID)
	if err != nil {
		return 0, errs.New(err)
	}
	return int64(count), nil
}

var _ repo.AIChatRepo = (*AIChatRepo)(nil)
