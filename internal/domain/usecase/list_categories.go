package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListCategoriesUseCase struct {
	cr repo.CategoryRepo
}

func NewListCategoriesUseCase(
	cr repo.CategoryRepo,
) *ListCategoriesUseCase {
	return &ListCategoriesUseCase{
		cr: cr,
	}
}

type ListCategoriesUseCaseInput struct {
	PaginationInput
}

func (uc *ListCategoriesUseCase) Execute(
	ctx context.Context,
	in ListCategoriesUseCaseInput,
) (*entity.PaginatedList[entity.Category], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	var errCh = make(chan error, 2)
	defer close(errCh)
	var categories []entity.Category
	var count int64

	go func() {
		var err error
		categories, err = uc.cr.SearchCategories(
			ctx,
			repo.SearchCategoriesParams{
				Search: in.Search,
				Offset: uint(offset),
				Limit:  uint(in.PageSize),
			},
		)
		errCh <- errs.New(err)
	}()

	go func() {
		var err error
		count, err = uc.cr.CountSearchCategories(
			ctx,
			in.Search,
		)
		errCh <- errs.New(err)
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
