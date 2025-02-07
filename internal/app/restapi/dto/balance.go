package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type GetBalanceResponse struct {
	usecase.GetBalanceOutput
}
