package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListTransactionCategories struct {
	cr repo.CategoryRepo
}

func NewListTransactionCategories(
	cr repo.CategoryRepo,
) *ListTransactionCategories {
	return &ListTransactionCategories{
		cr: cr,
	}
}

type ListCategoriesInput struct {
	repo.TransactionCategoryOptions
	PaginationInput
}

func (uc *ListTransactionCategories) Execute(
	ctx context.Context,
	in ListCategoriesInput,
) (*entity.PaginatedList[entity.TransactionCategory], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var categories []entity.TransactionCategory
	var count int64

	options := []repo.TransactionCategoryOption{}

	if in.Search != "" {
		options = append(
			options,
			repo.WithTransactionCategorySearch(in.Search),
		)
	}

	g.Go(func() error {
		var err error
		count, err = uc.cr.CountTransactionCategories(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithTransactionCategoryPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListTransactionCategories(
			gCtx,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.TransactionCategory]{
		Items: categories,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
