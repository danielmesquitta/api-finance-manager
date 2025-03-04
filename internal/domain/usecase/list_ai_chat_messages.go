package usecase

import (
	"context"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type ListAIChatMessages struct {
	pmr repo.AIChatMessageRepo
}

func NewListAIChatMessages(
	pmr repo.AIChatMessageRepo,
) *ListAIChatMessages {
	return &ListAIChatMessages{
		pmr: pmr,
	}
}

type ListAIChatMessagesInput struct {
	PaginationInput
	repo.AIChatMessageOptions
}

func (uc *ListAIChatMessages) Execute(
	ctx context.Context,
	in ListAIChatMessagesInput,
) (*entity.PaginatedList[entity.AIChatMessage], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var aiChatMessages []entity.AIChatMessage
	var count int64

	options := []repo.AIChatMessageOption{}

	if in.AIChatID != uuid.Nil {
		options = append(
			options,
			repo.WithAIChatMessageAIChat(in.AIChatID),
		)
	}

	g.Go(func() error {
		var err error
		count, err = uc.pmr.CountAIChatMessages(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithAIChatMessagePagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		aiChatMessages, err = uc.pmr.ListAIChatMessages(
			gCtx,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.AIChatMessage]{
		Items: aiChatMessages,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
