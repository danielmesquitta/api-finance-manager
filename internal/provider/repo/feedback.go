package repo

import (
	"context"
)

type FeedbackRepo interface {
	CreateFeedback(
		ctx context.Context,
		params CreateFeedbackParams,
	) error
}
