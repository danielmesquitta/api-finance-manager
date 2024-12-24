package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type SyncCategories struct {
	o  openfinance.Client
	ir repo.CategoryRepo
}

func NewSyncCategories(
	o openfinance.Client,
	ir repo.CategoryRepo,
) *SyncCategories {
	return &SyncCategories{
		o:  o,
		ir: ir,
	}
}

func (uc *SyncCategories) Execute(ctx context.Context) error {
	var openFinanceCategories, institutions []entity.Category

	var errCh = make(chan error, 2)
	defer close(errCh)

	go func() {
		var err error
		openFinanceCategories, err = uc.o.ListCategories(ctx)
		errCh <- err
	}()

	go func() {
		var err error
		institutions, err = uc.ir.ListCategories(ctx)
		errCh <- err
	}()

	for i := 0; i < cap(errCh); i++ {
		if err := <-errCh; err != nil {
			return errs.New(err)
		}
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

	if err := uc.ir.CreateCategories(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
