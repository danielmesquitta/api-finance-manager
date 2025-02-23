package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListTransactions struct {
	v  *validator.Validator
	tr repo.TransactionRepo
}

func NewListTransactions(
	v *validator.Validator,
	tr repo.TransactionRepo,
) *ListTransactions {
	return &ListTransactions{
		v:  v,
		tr: tr,
	}
}

type ListTransactionsInput struct {
	PaginationInput
	repo.TransactionOptions
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

func (uc *ListTransactions) Execute(
	ctx context.Context,
	in ListTransactionsInput,
) (*entity.PaginatedList[entity.FullTransaction], error) {
	if err := uc.v.Validate(in); err != nil {
		return nil, errs.New(err)
	}

	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var transactions []entity.FullTransaction
	var count int64

	opts := prepareTransactionOptions(in.TransactionOptions)

	g.Go(func() error {
		var err error
		count, err = uc.tr.CountTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	opts = append(
		opts,
		repo.WithTransactionPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		transactions, err = uc.tr.ListFullTransactions(
			gCtx,
			in.UserID,
			opts...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.FullTransaction]{
		Items: transactions,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
