package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type SignInRequest struct {
	Provider entity.Provider `json:"provider"`
}

type SignInResponse struct {
	usecase.SignInUseCaseOutput
}
