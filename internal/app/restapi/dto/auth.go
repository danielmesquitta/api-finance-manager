package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type SignInRequestDTO struct {
	Provider entity.Provider `json:"provider"`
}

type SignInResponseDTO struct {
	usecase.SignInUseCaseOutput
}
