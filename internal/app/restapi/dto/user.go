package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type GetUserProfileResponse struct {
	entity.User
}

type UpdateProfileRequest struct {
	usecase.UpdateUserInput
}
