package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type SignInRequest struct {
	usecase.SignInInput
}

type SignInResponse struct {
	usecase.SignInOutput
}
