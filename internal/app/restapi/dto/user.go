package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
)

type UserProfileResponse struct {
	entity.User
}

type UpdateProfileRequest struct {
	usecase.UpdateUserInput
}
