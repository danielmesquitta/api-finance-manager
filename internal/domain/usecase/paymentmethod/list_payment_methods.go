package paymentmethod

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"golang.org/x/sync/errgroup"
)

type ListPaymentMethodsUseCase struct {
	pmr repo.PaymentMethodRepo
}

func NewListPaymentMethodsUseCase(
	pmr repo.PaymentMethodRepo,
) *ListPaymentMethodsUseCase {
	return &ListPaymentMethodsUseCase{
		pmr: pmr,
	}
}

type ListPaymentMethodsUseCaseInput struct {
	usecase.PaginationInput
	repo.PaymentMethodOptions
}

func (uc *ListPaymentMethodsUseCase) Execute(
	ctx context.Context,
	in ListPaymentMethodsUseCaseInput,
) (*entity.PaginatedList[entity.PaymentMethod], error) {
	offset := usecase.PreparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var paymentMethods []entity.PaymentMethod
	var count int64

	options := []repo.PaymentMethodOption{}

	if in.Search != "" {
		options = append(
			options,
			repo.WithPaymentMethodSearch(in.Search),
		)
	}

	g.Go(func() error {
		var err error
		count, err = uc.pmr.CountPaymentMethods(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithPaymentMethodPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		paymentMethods, err = uc.pmr.ListPaymentMethods(
			gCtx,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.PaymentMethod]{
		Items: paymentMethods,
	}

	usecase.PreparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
