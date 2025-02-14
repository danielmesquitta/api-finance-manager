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
	var openFinanceCategories, categories []entity.TransactionCategory

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		openFinanceCategories, err = uc.o.ListTransactionCategories(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListTransactionCategories(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	categoriesByExternalID := make(map[string]entity.TransactionCategory)
	for _, i := range categories {
		categoriesByExternalID[i.ExternalID] = i
	}

	params := []repo.CreateTransactionCategoriesParams{}
	for _, i := range openFinanceCategories {
		if _, ok := categoriesByExternalID[i.ExternalID]; ok {
			continue
		}
		param := repo.CreateTransactionCategoriesParams{}
		if err := copier.Copy(&param, i); err != nil {
			return errs.New(err)
		}
		params = append(params, param)
	}

	if err := uc.cr.CreateTransactionCategories(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
