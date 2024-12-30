package usecase

import (
	"context"

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

	var errCh = make(chan error, 2)
	defer close(errCh)
	var categories []entity.Category
	var count int64

	go func() {
		var err error
		categories, err = uc.cr.ListCategories(
			ctx,
			repo.WithCategoriesPagination(uint(in.PageSize), uint(offset)),
		)
		errCh <- err
	}()

	go func() {
		var err error
		count, err = uc.cr.CountCategories(
			ctx,
			repo.WithCategoriesSearch(in.Search),
		)
		errCh <- err
	}()

	for i := 0; i < cap(errCh); i++ {
		if err := <-errCh; err != nil {
			return nil, errs.New(err)
		}
	}

	return &entity.PaginatedList[entity.Category]{
		Items:      categories,
		TotalItems: int(count),
		Page:       in.Page,
		PageSize:   in.PageSize,
		TotalPages: int(count) / in.PageSize,
	}, nil
}
