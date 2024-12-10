package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type SignInRequestDTO struct {
	Provider string `json:"provider"`
}

type SignInResponseDTO struct {
	usecase.SignInUseCaseOutput
}
