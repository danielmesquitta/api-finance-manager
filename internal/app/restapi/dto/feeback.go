package dto

import "github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"

type CreateFeedbackRequest struct {
	usecase.CreateFeedbackInput
}
