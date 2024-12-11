package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type CreateAccountUseCase struct {
	ur repo.UserRepo
	ar repo.AccountRepo
	oc openfinance.Client
}

func (uc *CreateAccountUseCase) NewCreateAccountUseCase(
	ur repo.UserRepo,
	ar repo.AccountRepo,
	oc openfinance.Client,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		ur: ur,
		ar: ar,
		oc: oc,
	}
}

func (uc *CreateAccountUseCase) Execute(
	ctx context.Context,
	userID string,
	accountIDs []string,
) (*entity.Account, error) {
	return nil, nil
}
