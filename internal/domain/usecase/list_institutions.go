package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListInstitutions struct {
	ir repo.InstitutionRepo
}

func NewListInstitutions(
	ir repo.InstitutionRepo,
) *ListInstitutions {
	return &ListInstitutions{
		ir: ir,
	}
}

type ListInstitutionsInput struct {
	repo.InstitutionOptions
	PaginationInput
}

func (uc *ListInstitutions) Execute(
	ctx context.Context,
	in ListInstitutionsInput,
) (*entity.PaginatedList[entity.Institution], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var institutions []entity.Institution
	var count int64

	options := []repo.InstitutionOption{}

	if in.Search != "" {
		options = append(options, repo.WithInstitutionsSearch(in.Search))
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
		repo.WithInstitutionsPagination(in.PageSize, offset),
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

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
