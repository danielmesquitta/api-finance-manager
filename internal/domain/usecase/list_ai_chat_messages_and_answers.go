package usecase

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListAIChatMessagesAndAnswers struct {
	acr repo.AIChatRepo
}

func NewListAIChatMessagesAndAnswers(
	acr repo.AIChatRepo,
) *ListAIChatMessagesAndAnswers {
	return &ListAIChatMessagesAndAnswers{
		acr: acr,
	}
}

type ListAIChatMessagesAndAnswersInput struct {
	PaginationInput
	AIChatID uuid.UUID `validate:"required"`
}

func (uc *ListAIChatMessagesAndAnswers) Execute(
	ctx context.Context,
	in ListAIChatMessagesAndAnswersInput,
) (*entity.PaginatedList[entity.AIChatMessageAndAnswer], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var (
		messagesAndAnswers []entity.AIChatMessageAndAnswer
		count              int64
	)

	g.Go(func() error {
		var err error
		count, err = uc.acr.CountAIChatMessagesAndAnswers(
			gCtx,
			in.AIChatID,
		)
		return err
	})

	g.Go(func() error {
		var err error
		messagesAndAnswers, err = uc.acr.ListAIChatMessagesAndAnswers(
			gCtx,
			repo.ListAIChatMessagesAndAnswersParams{
				AiChatID: in.AIChatID,
				Limit:    int32(in.PageSize),
				Offset:   int32(offset),
			},
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.AIChatMessageAndAnswer]{
		Items: messagesAndAnswers,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
