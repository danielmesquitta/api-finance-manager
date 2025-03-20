package institution

import (
	"context"

	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

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
	var openFinanceInstitutions, institutions []entity.Institution

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		openFinanceInstitutions, err = uc.o.ListInstitutions(
			gCtx,
			openfinance.WithInstitutionTypes(
				[]string{"PERSONAL_BANK", "INVESTMENT"},
			),
		)
		return err
	})

	g.Go(func() error {
		var err error
		institutions, err = uc.ir.ListInstitutions(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
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
