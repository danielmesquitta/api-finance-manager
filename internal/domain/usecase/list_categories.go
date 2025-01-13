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
	PaginationInput
}

func (uc *ListCategories) Execute(
	ctx context.Context,
	in ListCategoriesInput,
) (*entity.PaginatedList[entity.Category], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, ctx := errgroup.WithContext(ctx)
	var categories []entity.Category
	var count int64

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListCategories(
			ctx,
			repo.WithCategoriesPagination(uint(in.PageSize), uint(offset)),
		)
		return err
	})

	g.Go(func() error {
		var err error
		count, err = uc.cr.CountCategories(
			ctx,
			repo.WithCategoriesSearch(in.Search),
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	return &entity.PaginatedList[entity.Category]{
		Items:      categories,
		TotalItems: int(count),
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalPages: int(count) / in.PageSize,
	}, nil
}
