package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListCategories struct {
	cr repo.CategoryRepo
}

func NewListCategories(
	cr repo.CategoryRepo,
) *ListCategories {
	return &ListCategories{
		cr: cr,
	}
}

type ListCategoriesInput struct {
	repo.ListCategoriesOptions
	PaginationInput
}

func (uc *ListCategories) Execute(
	ctx context.Context,
	in ListCategoriesInput,
) (*entity.PaginatedList[entity.Category], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var categories []entity.Category
	var count int64

	options := []repo.ListCategoriesOption{}

	if in.Search != "" {
		options = append(options, repo.WithCategoriesSearch(in.Search))
	}

	g.Go(func() error {
		var err error
		count, err = uc.cr.CountCategories(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithCategoriesPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListCategories(
			gCtx,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.Category]{
		Items: categories,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
