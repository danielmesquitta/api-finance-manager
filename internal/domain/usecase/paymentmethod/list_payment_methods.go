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
	g, gCtx := errgroup.WithContext(ctx)
	var paymentMethods []entity.PaymentMethod
	var count int64

	g.Go(func() error {
		var err error
		count, err = uc.pmr.CountPaymentMethods(
			gCtx,
			in.PaymentMethodOptions,
		)
		return err
	})

	in.PaymentMethodOptions.Limit, in.PaymentMethodOptions.Offset = usecase.PreparePaginationInput(
		in.PaginationInput,
	)

	g.Go(func() error {
		var err error
		paymentMethods, err = uc.pmr.ListPaymentMethods(
			gCtx,
			in.PaymentMethodOptions,
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
