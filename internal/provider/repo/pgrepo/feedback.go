package pgrepo

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/db/sqlc"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type FeedbackRepo struct {
	db *db.DB
}

func NewFeedbackRepo(db *db.DB) *FeedbackRepo {
	return &FeedbackRepo{
		db: db,
	}
}

func (r *FeedbackRepo) CreateFeedback(
	ctx context.Context,
	params repo.CreateFeedbackParams,
) error {
	dbParams := sqlc.CreateFeedbackParams{}
	if err := copier.Copy(&dbParams, params); err != nil {
		return errs.New(err)
	}

	tx := r.db.UseTx(ctx)
	if err := tx.CreateFeedback(ctx, dbParams); err != nil {
		return errs.New(err)
	}

	return nil
}

var _ repo.FeedbackRepo = (*FeedbackRepo)(nil)
