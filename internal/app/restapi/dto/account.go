package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/account"
)

type CreateAccountsRequest struct {
	account.CreateAccountsUseCaseInput
}

type GetAccountsBalanceResponse struct {
	account.GetAccountsBalanceUseCaseOutput
}
