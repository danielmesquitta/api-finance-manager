package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type SyncInstitutionsUseCase struct {
	o  openfinance.Client
	ir repo.InstitutionRepo
}

func NewSyncInstitutionsUseCase(
	o openfinance.Client,
	ir repo.InstitutionRepo,
) *SyncInstitutionsUseCase {
	return &SyncInstitutionsUseCase{
		o:  o,
		ir: ir,
	}
}

func (uc *SyncInstitutionsUseCase) Execute(ctx context.Context) error {
	var openFinanceInstitutions, institutions []entity.Institution

	var errCh = make(chan error, 2)
	defer close(errCh)

	go func() {
		var err error
		openFinanceInstitutions, err = uc.o.ListInstitutions(ctx)
		errCh <- err
	}()

	go func() {
		var err error
		institutions, err = uc.ir.ListInstitutions(ctx)
		errCh <- err
	}()

	for i := 0; i < cap(errCh); i++ {
		if err := <-errCh; err != nil {
			return errs.New(err)
		}
	}

	institutionsByExternalID := make(map[string]entity.Institution)
	for _, i := range institutions {
		institutionsByExternalID[i.ExternalID] = i
	}

	params := []repo.CreateInstitutionsParams{}
	for _, i := range openFinanceInstitutions {
		if _, ok := institutionsByExternalID[i.ExternalID]; ok {
			continue
		}
		param := repo.CreateInstitutionsParams{}
		if err := copier.Copy(&param, i); err != nil {
			return errs.New(err)
		}
		params = append(params, param)
	}

	if err := uc.ir.CreateInstitutions(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
