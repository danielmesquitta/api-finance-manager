package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/auth"
)

type SignInRequest struct {
	auth.SignInUseCaseInput
}

type SignInResponse struct {
	auth.SignInUseCaseOutput
}
