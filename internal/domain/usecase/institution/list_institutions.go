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
	offset := usecase.PreparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var institutions []entity.Institution
	var count int64

	options := []repo.InstitutionOption{}
	if len(in.UserIDs) > 0 {
		options = append(options, repo.WithInstitutionUsers(in.UserIDs))
	}

	if in.Search != "" {
		options = append(options, repo.WithInstitutionSearch(in.Search))
	}

	g.Go(func() error {
		var err error
		count, err = uc.ir.CountInstitutions(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithInstitutionPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		institutions, err = uc.ir.ListInstitutions(
			gCtx,
			options...,
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
