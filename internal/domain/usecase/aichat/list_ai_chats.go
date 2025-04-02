package aichat

import (
	"context"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"golang.org/x/sync/errgroup"
)

type ListAIChatsUseCase struct {
	pmr repo.AIChatRepo
}

func NewListAIChatsUseCase(
	pmr repo.AIChatRepo,
) *ListAIChatsUseCase {
	return &ListAIChatsUseCase{
		pmr: pmr,
	}
}

type ListAIChatsUseCaseInput struct {
	usecase.PaginationInput
	repo.AIChatOptions
}

func (uc *ListAIChatsUseCase) Execute(
	ctx context.Context,
	in ListAIChatsUseCaseInput,
) (*entity.PaginatedList[entity.AIChat], error) {

	g, gCtx := errgroup.WithContext(ctx)
	var aiChats []entity.AIChat
	var count int64

	g.Go(func() error {
		var err error
		count, err = uc.pmr.CountAIChats(
			gCtx,
			in.AIChatOptions,
		)
		return err
	})

	in.Limit, in.Offset = usecase.PreparePaginationInput(
		in.PaginationInput,
	)

	g.Go(func() error {
		var err error
		aiChats, err = uc.pmr.ListAIChats(
			gCtx,
			in.AIChatOptions,
		)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, errs.New(err)
	}

	out := entity.PaginatedList[entity.AIChat]{
		Items: aiChats,
	}

	usecase.PreparePaginationOutput(&out, in.PaginationInput, count)

	return &out, nil
}
