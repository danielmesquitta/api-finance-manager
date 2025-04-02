package transactioncategory

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListTransactionCategoriesUseCase struct {
	cr repo.TransactionCategoryRepo
}

func NewListTransactionCategoriesUseCase(
	cr repo.TransactionCategoryRepo,
) *ListTransactionCategoriesUseCase {
	return &ListTransactionCategoriesUseCase{
		cr: cr,
	}
}

type ListCategoriesInput struct {
	repo.TransactionCategoryOptions
	usecase.PaginationInput
}

func (uc *ListTransactionCategoriesUseCase) Execute(
	ctx context.Context,
	in ListCategoriesInput,
) (*entity.PaginatedList[entity.TransactionCategory], error) {
	g, gCtx := errgroup.WithContext(ctx)
	var categories []entity.TransactionCategory
	var count int64

	g.Go(func() (err error) {
		count, err = uc.cr.CountTransactionCategories(
			gCtx,
			in.TransactionCategoryOptions,
		)
		return err
	})

	in.Limit, in.Offset = usecase.PreparePaginationInput(
		in.PaginationInput,
	)

	g.Go(func() (err error) {
		categories, err = uc.cr.ListTransactionCategories(
			gCtx,
			in.TransactionCategoryOptions,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.TransactionCategory]{
		Items: categories,
	}

	usecase.PreparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
