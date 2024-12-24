package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/google/uuid"
)

type RefreshToken struct {
	si *SignIn
}

func NewRefreshToken(
	si *SignIn,
) *RefreshToken {
	return &RefreshToken{
		si: si,
	}
}

func (uc *RefreshToken) Execute(
	ctx context.Context,
	userID uuid.UUID,
) (*SignInOutput, error) {
	return uc.si.Execute(ctx, SignInInput{
		Provider: entity.ProviderRefresh,
		UserID:   userID,
	})
}
