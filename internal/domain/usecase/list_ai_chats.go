package usecase

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"golang.org/x/sync/errgroup"
)

type ListAIChats struct {
	pmr repo.AIChatRepo
}

func NewListAIChats(
	pmr repo.AIChatRepo,
) *ListAIChats {
	return &ListAIChats{
		pmr: pmr,
	}
}

type ListAIChatsInput struct {
	PaginationInput
	repo.AIChatOptions
}

func (uc *ListAIChats) Execute(
	ctx context.Context,
	in ListAIChatsInput,
) (*entity.PaginatedList[entity.AIChat], error) {
	offset := preparePaginationInput(&in.PaginationInput)

	g, gCtx := errgroup.WithContext(ctx)
	var aiChats []entity.AIChat
	var count int64

	options := []repo.AIChatOption{}

	if in.Search != "" {
		options = append(
			options,
			repo.WithAIChatSearch(in.Search),
		)
	}

	g.Go(func() error {
		var err error
		count, err = uc.pmr.CountAIChats(
			gCtx,
			options...,
		)
		return err
	})

	options = append(
		options,
		repo.WithAIChatPagination(in.PageSize, offset),
	)

	g.Go(func() error {
		var err error
		aiChats, err = uc.pmr.ListAIChats(
			gCtx,
			options...,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.AIChat]{
		Items: aiChats,
	}

	preparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
