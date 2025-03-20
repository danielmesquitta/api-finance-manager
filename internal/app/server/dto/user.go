package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/user"
)

type GetUserProfileResponse struct {
	entity.User
}

type UpdateProfileRequest struct {
	user.UpdateUserUseCaseInput
}
