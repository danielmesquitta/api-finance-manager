package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type RefreshTokenUseCase struct {
	si *SignInUseCase
}

func NewRefreshTokenUseCase(
	si *SignInUseCase,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		si: si,
	}
}

func (uc *RefreshTokenUseCase) Execute(
	ctx context.Context,
	userID uuid.UUID,
) (*SignInUseCaseOutput, error) {
	return uc.si.Execute(ctx, SignInUseCaseInput{
		Provider: entity.ProviderRefresh,
		UserID:   userID,
	})
}
