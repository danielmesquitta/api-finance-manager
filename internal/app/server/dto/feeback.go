package dto

import (
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/feedback"
)

type CreateFeedbackRequest struct {
	feedback.CreateFeedbackUseCaseInput
}
