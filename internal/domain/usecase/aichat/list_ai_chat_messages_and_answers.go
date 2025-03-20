package aichat

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type ListAIChatMessagesAndAnswersUseCase struct {
	acr repo.AIChatRepo
}

func NewListAIChatMessagesAndAnswersUseCase(
	acr repo.AIChatRepo,
) *ListAIChatMessagesAndAnswersUseCase {
	return &ListAIChatMessagesAndAnswersUseCase{
		acr: acr,
	}
}

type ListAIChatMessagesAndAnswersUseCaseInput struct {
	usecase.PaginationInput
	AIChatID uuid.UUID `validate:"required"`
}

func (uc *ListAIChatMessagesAndAnswersUseCase) Execute(
	ctx context.Context,
	in ListAIChatMessagesAndAnswersUseCaseInput,
) (*entity.PaginatedList[entity.AIChatMessageAndAnswer], error) {
	offset := usecase.PreparePaginationInput(&in.PaginationInput)

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

	usecase.PreparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
