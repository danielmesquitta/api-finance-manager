package usecase

import (
	"context"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type SyncCategories struct {
	o  openfinance.Client
	cr repo.CategoryRepo
}

func NewSyncCategories(
	o openfinance.Client,
	cr repo.CategoryRepo,
) *SyncCategories {
	return &SyncCategories{
		o:  o,
		cr: cr,
	}
}

func (uc *SyncCategories) Execute(ctx context.Context) error {
	var openFinanceCategories, institutions []entity.Category

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		openFinanceCategories, err = uc.o.ListCategories(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		institutions, err = uc.cr.ListCategories(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	institutionsByExternalID := make(map[string]entity.Category)
	for _, i := range institutions {
		institutionsByExternalID[i.ExternalID] = i
	}

	params := []repo.CreateCategoriesParams{}
	for _, i := range openFinanceCategories {
		if _, ok := institutionsByExternalID[i.ExternalID]; ok {
			continue
		}
		param := repo.CreateCategoriesParams{}
		if err := copier.Copy(&param, i); err != nil {
			return errs.New(err)
		}
		params = append(params, param)
	}

	if err := uc.cr.CreateCategories(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
