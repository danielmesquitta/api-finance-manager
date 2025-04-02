package institution

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListInstitutionsUseCase struct {
	ir repo.InstitutionRepo
}

func NewListInstitutionsUseCase(
	ir repo.InstitutionRepo,
) *ListInstitutionsUseCase {
	return &ListInstitutionsUseCase{
		ir: ir,
	}
}

type ListInstitutionsUseCaseInput struct {
	repo.InstitutionOptions
	usecase.PaginationInput
}

func (uc *ListInstitutionsUseCase) Execute(
	ctx context.Context,
	in ListInstitutionsUseCaseInput,
) (*entity.PaginatedList[entity.Institution], error) {
	g, gCtx := errgroup.WithContext(ctx)
	var institutions []entity.Institution
	var count int64

	g.Go(func() error {
		var err error
		count, err = uc.ir.CountInstitutions(
			gCtx,
			in.InstitutionOptions,
		)
		return err
	})

	in.Limit, in.Offset = usecase.PreparePaginationInput(
		in.PaginationInput,
	)

	g.Go(func() error {
		var err error
		institutions, err = uc.ir.ListInstitutions(
			gCtx,
			in.InstitutionOptions,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.Institution]{
		Items: institutions,
	}

	usecase.PreparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
