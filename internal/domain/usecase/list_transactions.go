package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListTransactions struct {
	cr repo.TransactionRepo
}

func NewListTransactions(
	cr repo.TransactionRepo,
) *ListTransactions {
	return &ListTransactions{
		cr: cr,
	}
}

type ListTransactionsInput struct {
	PaginationInput
	repo.ListTransactionsOptions
	UserID uuid.UUID `json:"-"`
}

func (uc *ListTransactions) Execute(
	ctx context.Context,
	in ListTransactionsInput,
) (*entity.PaginatedList[entity.Transaction], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var transactions []entity.Transaction
	var count int64

	options := []repo.ListTransactionsOption{}

	if in.Search != "" {
		options = append(options, repo.WithTransactionsSearch(in.Search))
	}

	if in.CategoryID != uuid.Nil {
		options = append(
			options,
			repo.WithTransactionCategory(in.CategoryID),
		)
	}

	if in.InstitutionID != uuid.Nil {
		options = append(
			options,
			repo.WithTransactionInstitution(in.InstitutionID),
		)
	}

	if !in.StartDate.IsZero() {
		options = append(
			options,
			repo.WithTransactionDateAfter(in.StartDate),
		)
	}

	if !in.EndDate.IsZero() {
		options = append(
			options,
			repo.WithTransactionDateBefore(in.EndDate),
		)
	}

	if in.IsExpense {
		options = append(
			options,
			repo.WithTransactionIsExpense(in.IsExpense),
		)
	}

	if in.IsIncome {
		options = append(
			options,
			repo.WithTransactionIsIncome(in.IsIncome),
		)
	}

	if in.PaymentMethod != "" {
		options = append(
			options,
			repo.WithTransactionPaymentMethod(in.PaymentMethod),
		)
	}

	g.Go(func() error {
		var err error
		count, err = uc.cr.CountTransactions(
			gCtx,
			in.UserID,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithTransactionsPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		transactions, err = uc.cr.ListTransactions(
			gCtx,
			in.UserID,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.Transaction]{
		Items: transactions,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
