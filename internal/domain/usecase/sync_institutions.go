package usecase

import (
	"context"
	"strconv"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
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
	var openFinanceInstitutions []openfinance.Institution
	var institutions []entity.Institution

	var errCh = make(chan error, 2)
	defer close(errCh)

	go func() {
		var err error
		openFinanceInstitutions, err = uc.o.ListInstitutions(ctx)
		errCh <- errs.New(err)
	}()

	go func() {
		var err error
		institutions, err = uc.ir.ListInstitutions(ctx)
		errCh <- errs.New(err)
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

	params := []repo.CreateManyInstitutionsParams{}
	for _, i := range openFinanceInstitutions {
		externalID := strconv.Itoa(i.ID)
		if _, ok := institutionsByExternalID[externalID]; ok {
			continue
		}

		var logo *string
		if i.ImageURL != "" {
			logo = &i.ImageURL
		}

		params = append(params, repo.CreateManyInstitutionsParams{
			ExternalID: externalID,
			Name:       i.Name,
			Logo:       logo,
		})
	}

	if err := uc.ir.CreateManyInstitutions(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
