package account

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/ptr"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type SyncAccountsBalancesUseCase struct {
	e   *env.Env
	tx  tx.TX
	o   openfinance.Client
	c   cache.Cache
	ar  repo.AccountRepo
	abr repo.AccountBalanceRepo
}

func NewSyncAccountsBalancesUseCase(
	e *env.Env,
	tx tx.TX,
	o openfinance.Client,
	c cache.Cache,
	ar repo.AccountRepo,
	abr repo.AccountBalanceRepo,
) *SyncAccountsBalancesUseCase {
	return &SyncAccountsBalancesUseCase{
		e:   e,
		tx:  tx,
		o:   o,
		c:   c,
		ar:  ar,
		abr: abr,
	}
}

func (uc *SyncAccountsBalancesUseCase) Execute(ctx context.Context) error {
	offset := 0
	cacheExp := time.Hour * 12
	if _, err := uc.c.Scan(ctx, cache.KeySyncBalancesOffset, &offset); err != nil {
		return errs.New(err)
	}

	if offset == -1 {
		slog.Info("sync-balances: already completed")
		return nil
	}

	accounts, err := uc.ar.ListFullAccounts(
		ctx,
		repo.AccountOptions{
			Limit:                uint(uc.e.SyncBalancesMaxAccounts),
			Offset:               uint(offset),
			IsSubscriptionActive: ptr.New(true),
		},
	)
	if err != nil {
		return errs.New(err)
	}

	if len(accounts) == 0 {
		if err := uc.c.Set(ctx, cache.KeySyncBalancesOffset, -1, cacheExp); err != nil {
			return errs.New(err)
		}

		slog.Info("sync-balances: completed")
		return nil
	}

	accountsByUserID := make(map[uuid.UUID][]entity.FullAccount)
	accountsByExternalIDs := make(map[string]entity.FullAccount)
	for _, account := range accounts {
		if account.UserID == nil {
			slog.Error(
				"sync-balances: account without user id",
				"account", account,
			)
			continue
		}

		accountsByUserID[*account.UserID] = append(
			accountsByUserID[*account.UserID],
			account,
		)
		accountsByExternalIDs[account.ExternalID] = account
	}

	g := errgroup.Group{}
	for userID, userAccounts := range accountsByUserID {
		g.Go(func() error {
			if err := uc.syncUserBalance(
				ctx, userID, userAccounts, accountsByExternalIDs,
			); err != nil {
				return errs.New(err)
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if len(accounts) < uc.e.SyncBalancesMaxAccounts {
		if err := uc.c.Set(ctx, cache.KeySyncBalancesOffset, -1, cacheExp); err != nil {
			return errs.New(err)
		}
		slog.Info("sync-balances: completed")
		return nil
	}

	offset += uc.e.SyncBalancesMaxAccounts
	if err := uc.c.Set(ctx, cache.KeySyncBalancesOffset, offset, cacheExp); err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncAccountsBalancesUseCase) syncUserBalance(
	ctx context.Context,
	userID uuid.UUID,
	userAccounts []entity.FullAccount,
	accountsByExternalIDs map[string]entity.FullAccount,
) error {
	userInstitutionExternalIDs := map[string]struct{}{}
	for _, account := range userAccounts {
		if account.UserInstitutionExternalID == nil {
			slog.Error(
				"sync-balances: account without user institution external id",
				"account",
				account,
			)
			continue
		}

		userInstitutionExternalIDs[*account.UserInstitutionExternalID] = struct{}{}
	}

	err := uc.tx.Do(ctx, func(ctx context.Context) error {
		g, gCtx := errgroup.WithContext(ctx)
		for userInstitutionExternalID := range userInstitutionExternalIDs {
			g.Go(func() error {
				openFinanceAccounts, err := uc.o.ListAccounts(
					gCtx,
					userInstitutionExternalID,
				)
				if err != nil {
					return errs.New(err)
				}

				if err := uc.createAccountBalances(
					gCtx,
					openFinanceAccounts,
					accountsByExternalIDs,
				); err != nil {
					slog.Error(
						"sync-balances: error creating account balances",
						"user_id", userID,
						"error", err,
					)
					return errs.New(err)
				}

				return nil
			})
		}

		if err := g.Wait(); err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncAccountsBalancesUseCase) createAccountBalances(
	ctx context.Context,
	openFinanceAccounts []openfinance.Account,
	accountsByExternalIDs map[string]entity.FullAccount,
) error {
	var params []repo.CreateAccountBalancesParams
	for _, openFinanceAccount := range openFinanceAccounts {
		account, ok := accountsByExternalIDs[openFinanceAccount.ExternalID]
		if !ok {
			slog.Error(
				"sync-balances: account not found",
				"external_id", openFinanceAccount.ExternalID,
			)
			continue
		}
		params = append(params, repo.CreateAccountBalancesParams{
			AccountID: account.ID,
			Amount:    openFinanceAccount.Balance,
		})
	}

	if err := uc.abr.CreateAccountBalances(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
