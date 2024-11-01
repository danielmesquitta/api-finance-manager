package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/google/uuid"
)

type CreateAccountUseCase struct {
	ur repo.UserRepo
	ar repo.AccountRepo
	oc openfinance.Client
}

func (uc *CreateAccountUseCase) NewCreateAccountUseCase(
	ur repo.UserRepo,
	ar repo.AccountRepo,
	oc openfinance.Client,
) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		ur: ur,
		ar: ar,
		oc: oc,
	}
}

func (uc *CreateAccountUseCase) Execute(
	ctx context.Context,
	userID string,
	accountIDs []string,
) (*entity.Account, error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errs.ErrUnauthorized
	}

	user, err := uc.ur.GetUserByID(ctx, userUUID)
	if err != nil {
		return nil, errs.New(err)
	}
	if user == nil {
		return nil, errs.ErrUserNotFound
	}

	accounts, err := uc.ar.ListAccountsByUserID(ctx, userUUID)
	if err != nil {
		return nil, errs.New(err)
	}

	userAccountIDs := map[string]struct{}{}
	for _, account := range accounts {
		userAccountIDs[account.ID.String()] = struct{}{}
	}

	notRegisteredAccountIDs := []string{}
	for _, accountID := range accountIDs {
		if _, ok := userAccountIDs[accountID]; ok {
			continue
		}
		notRegisteredAccountIDs = append(notRegisteredAccountIDs, accountID)
	}

	if len(notRegisteredAccountIDs) == 0 {
		return nil, errs.ErrAccountsAlreadyRegistered
	}

	jobsCount := len(notRegisteredAccountIDs)
	openFinanceAccountsCh := make(
		chan *openfinance.Account,
		jobsCount,
	)
	errsCh := make(chan error, jobsCount)
	wg := sync.WaitGroup{}
	wg.Add(jobsCount)
	for _, accountID := range notRegisteredAccountIDs {
		go func() {
			defer wg.Done()

			account, err := uc.oc.GetAccount(ctx, accountID)
			if err != nil {
				errsCh <- err
				return
			}
			openFinanceAccountsCh <- account
		}()
	}

	go func() {
		for e := range errsCh {
			err = errors.Join(err, e)
		}
	}()

	openFinanceAccounts := []*openfinance.Account{}
	go func() {
		for openFinanceAccount := range openFinanceAccountsCh {
			openFinanceAccounts = append(
				openFinanceAccounts,
				openFinanceAccount,
			)
		}
	}()

	wg.Wait()
	close(openFinanceAccountsCh)
	close(errsCh)

	return nil, nil
}
