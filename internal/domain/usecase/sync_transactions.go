package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type SyncTransactions struct {
	e   *config.Env
	o   openfinance.Client
	c   cache.Cache
	tx  tx.TX
	ar  repo.AccountRepo
	ur  repo.UserRepo
	tr  repo.TransactionRepo
	cr  repo.CategoryRepo
	pmr repo.PaymentMethodRepo
}

func NewSyncTransactions(
	e *config.Env,
	o openfinance.Client,
	c cache.Cache,
	tx tx.TX,
	ar repo.AccountRepo,
	ur repo.UserRepo,
	tr repo.TransactionRepo,
	cr repo.CategoryRepo,
	pmr repo.PaymentMethodRepo,
) *SyncTransactions {
	return &SyncTransactions{
		e:   e,
		o:   o,
		c:   c,
		tx:  tx,
		ar:  ar,
		ur:  ur,
		tr:  tr,
		cr:  cr,
		pmr: pmr,
	}
}

type SyncTransactionsInput struct {
	UserIDs []uuid.UUID `json:"user_ids"`
}

func (uc *SyncTransactions) Execute(
	ctx context.Context,
	in SyncTransactionsInput,
) error {
	isSyncingAllUsers := len(in.UserIDs) == 0
	cacheExp := time.Hour * 12
	offset := 0

	accountOpts := []repo.AccountOption{
		repo.WithAccountSubscriptionActive(true),
	}

	if isSyncingAllUsers {
		_, err := uc.c.Scan(ctx, cache.KeySyncTransactionsOffset, &offset)
		if err != nil {
			return errs.New(err)
		}

		if offset == -1 {
			slog.Info("sync transactions already completed")
			return nil
		}

		accountOpts = append(
			accountOpts,
			repo.WithAccountPagination(
				uint(uc.e.SyncTransactionsMaxAccounts),
				uint(offset),
			),
		)
	}

	if len(in.UserIDs) > 0 {
		accountOpts = append(
			accountOpts,
			repo.WithAccountUserIDs(in.UserIDs),
		)
	}

	g, gCtx := errgroup.WithContext(ctx)

	var (
		accounts       []entity.FullAccount
		categories     []entity.TransactionCategory
		paymentMethods []entity.PaymentMethod
	)

	g.Go(func() error {
		var err error
		accounts, err = uc.ar.ListFullAccounts(gCtx, accountOpts...)
		return err
	})

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListTransactionCategories(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		paymentMethods, err = uc.pmr.ListPaymentMethods(gCtx)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if len(accounts) == 0 {
		if !isSyncingAllUsers {
			return nil
		}

		err := uc.c.Set(
			ctx,
			cache.KeySyncTransactionsOffset,
			-1,
			cacheExp,
		)
		if err != nil {
			return errs.New(err)
		}

		slog.Info("sync transactions completed")
		return nil
	}

	accountsByUserID := make(map[uuid.UUID][]entity.FullAccount)
	accountsByID := make(map[uuid.UUID]entity.FullAccount)

	for _, account := range accounts {
		accountsByUserID[account.UserID] = append(
			accountsByUserID[account.UserID],
			account,
		)
		accountsByID[account.ID] = account
	}

	categoriesByExternalID := make(map[string]entity.TransactionCategory)
	for _, category := range categories {
		categoriesByExternalID[category.ExternalID] = category
	}

	openFinanceTransactionsByAccountID := make(
		map[uuid.UUID][]openfinance.Transaction,
	)
	for userID, userAccounts := range accountsByUserID {
		if len(userAccounts) == 0 {
			continue
		}

		lastSynchronizedAt := uc.calculateLastSynchronizedAt(
			userAccounts[0].SynchronizedAt,
		)

		for _, account := range userAccounts {
			ofTransactions, err := uc.listOpenFinanceTransactions(
				ctx,
				account.ExternalID,
				lastSynchronizedAt,
			)
			if err != nil {
				slog.Error(
					"error getting open finance transactions",
					"user", userID,
					"account", account,
					"err", err,
				)
				continue
			}
			openFinanceTransactionsByAccountID[account.ID] = ofTransactions
		}
	}

	paymentMethods, err := uc.syncPaymentMethods(
		ctx,
		openFinanceTransactionsByAccountID,
		paymentMethods,
	)
	if err != nil {
		return errs.New(err)
	}

	paymentMethodsByExternalID := make(map[string]entity.PaymentMethod)
	for _, pm := range paymentMethods {
		paymentMethodsByExternalID[pm.ExternalID] = pm
	}

	for userID, userAccounts := range accountsByUserID {
		if len(userAccounts) == 0 {
			continue
		}

		lastSynchronizedAt := uc.calculateLastSynchronizedAt(
			userAccounts[0].SynchronizedAt,
		)

		if err := uc.syncUserTransactions(
			ctx,
			userID,
			lastSynchronizedAt,
			accountsByID,
			categoriesByExternalID,
			paymentMethodsByExternalID,
			openFinanceTransactionsByAccountID,
		); err != nil {
			slog.Error(
				"error syncing user transactions",
				"user_id", userID,
				"accounts", userAccounts,
				"categories", categories,
				"err", err,
			)
			continue
		}
	}

	if isSyncingAllUsers {
		if len(accounts) < uc.e.SyncTransactionsMaxAccounts {
			err := uc.c.Set(
				ctx,
				cache.KeySyncTransactionsOffset,
				-1,
				cacheExp,
			)
			if err != nil {
				return errs.New(err)
			}

			slog.Info("sync transactions completed")
			return nil
		}

		offset += uc.e.SyncTransactionsMaxAccounts
		err := uc.c.Set(
			ctx,
			cache.KeySyncTransactionsOffset,
			offset,
			cacheExp,
		)
		if err != nil {
			return errs.New(err)
		}
	}

	return nil
}

func (uc *SyncTransactions) syncPaymentMethods(
	ctx context.Context,
	openFinanceTransactionsByAccountID map[uuid.UUID][]openfinance.Transaction,
	paymentMethods []entity.PaymentMethod,
) ([]entity.PaymentMethod, error) {
	paymentMethodNamesByExternalID := map[string]string{
		string(pluggy.PaymentMethodBOLETO):     "Boleto",
		string(pluggy.PaymentMethodCreditCard): "Cartão de crédito",
		string(pluggy.PaymentMethodDEBIT):      "Cartão de débito",
	}

	uniqueExternalIDs := make(map[string]struct{})
	for _, ofTransactions := range openFinanceTransactionsByAccountID {
		for _, ofTrans := range ofTransactions {
			uniqueExternalIDs[ofTrans.PaymentMethodExternalID] = struct{}{}
		}
	}

	for _, pm := range paymentMethods {
		delete(uniqueExternalIDs, pm.ExternalID)
	}

	params := make([]repo.CreatePaymentMethodsParams, 0, len(uniqueExternalIDs))
	for externalID := range uniqueExternalIDs {
		name, ok := paymentMethodNamesByExternalID[externalID]
		if !ok {
			name = externalID
		}
		params = append(params, repo.CreatePaymentMethodsParams{
			Name:       name,
			ExternalID: externalID,
		})
	}

	if err := uc.pmr.CreatePaymentMethods(ctx, params); err != nil {
		return nil, errs.New(err)
	}

	updated, err := uc.pmr.ListPaymentMethods(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	return updated, nil
}

func (uc *SyncTransactions) syncUserTransactions(
	ctx context.Context,
	userID uuid.UUID,
	lastSynchronizedAt time.Time,
	accountsByID map[uuid.UUID]entity.FullAccount,
	categoriesByExternalID map[string]entity.TransactionCategory,
	paymentMethodsByExternalID map[string]entity.PaymentMethod,
	openFinanceTransactionsByAccountID map[uuid.UUID][]openfinance.Transaction,
) error {
	transactions, err := uc.listRepoTransactions(
		ctx,
		userID,
		lastSynchronizedAt,
	)
	if err != nil {
		return errs.New(err)
	}

	transactionsByExternalID := make(
		map[string]entity.Transaction,
		len(transactions),
	)
	for _, t := range transactions {
		if t.ExternalID == nil {
			continue
		}
		transactionsByExternalID[*t.ExternalID] = t
	}

	params := uc.buildCreateTransactionsParams(
		userID,
		accountsByID,
		categoriesByExternalID,
		paymentMethodsByExternalID,
		openFinanceTransactionsByAccountID,
		transactionsByExternalID,
	)

	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if err := uc.tr.CreateTransactions(ctx, params); err != nil {
			return errs.New(err)
		}
		return uc.updateUserSynchronizedAt(ctx, userID)
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncTransactions) buildCreateTransactionsParams(
	userID uuid.UUID,
	accountsByID map[uuid.UUID]entity.FullAccount,
	categoriesByExternalID map[string]entity.TransactionCategory,
	paymentMethodsByExternalID map[string]entity.PaymentMethod,
	openFinanceTransactionsByAccountID map[uuid.UUID][]openfinance.Transaction,
	transactionsByExternalID map[string]entity.Transaction,
) []repo.CreateTransactionsParams {
	var params []repo.CreateTransactionsParams

	for accountID, ofTransactions := range openFinanceTransactionsByAccountID {
		account := accountsByID[accountID]

		for _, ofTrans := range ofTransactions {
			if ofTrans.ExternalID == nil {
				continue
			}

			if _, ok := transactionsByExternalID[*ofTrans.ExternalID]; ok {
				continue
			}

			categoryID := uc.getCategoryID(
				ofTrans.CategoryExternalID,
				categoriesByExternalID,
			)
			pm := paymentMethodsByExternalID[ofTrans.PaymentMethodExternalID]

			params = append(params, repo.CreateTransactionsParams{
				ExternalID:      ofTrans.ExternalID,
				Name:            ofTrans.Name,
				Amount:          ofTrans.Amount,
				PaymentMethodID: pm.ID,
				Date:            ofTrans.Date,
				UserID:          userID,
				AccountID:       &account.ID,
				InstitutionID:   &account.InstitutionID,
				CategoryID:      categoryID,
			})
		}
	}

	return params
}

func (uc *SyncTransactions) calculateLastSynchronizedAt(
	userSynchronizedAt *time.Time,
) time.Time {
	if userSynchronizedAt == nil {
		return time.Time{}
	}
	return getStartOfDay(*userSynchronizedAt)
}

func (uc *SyncTransactions) listOpenFinanceTransactions(
	ctx context.Context,
	accountExternalID string,
	lastSynchronizedAt time.Time,
) ([]openfinance.Transaction, error) {
	var opts []openfinance.TransactionOption
	if !lastSynchronizedAt.IsZero() {
		opts = append(
			opts,
			openfinance.WithTransactionDateAfter(lastSynchronizedAt),
		)
	}

	ofTrans, err := uc.o.ListTransactions(ctx, accountExternalID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return ofTrans, nil
}

func (uc *SyncTransactions) listRepoTransactions(
	ctx context.Context,
	userID uuid.UUID,
	lastSynchronizedAt time.Time,
) ([]entity.Transaction, error) {
	var opts []repo.TransactionOption
	if !lastSynchronizedAt.IsZero() {
		opts = append(opts, repo.WithTransactionDateAfter(lastSynchronizedAt))
	}

	txs, err := uc.tr.ListTransactions(ctx, userID, opts...)
	if err != nil {
		return nil, errs.New(err)
	}
	return txs, nil
}

func (uc *SyncTransactions) updateUserSynchronizedAt(
	ctx context.Context,
	userID uuid.UUID,
) error {
	yesterday := time.Now().AddDate(0, 0, -1)
	startOfYesterday := getStartOfDay(yesterday)

	err := uc.ur.UpdateUserSynchronizedAt(
		ctx,
		repo.UpdateUserSynchronizedAtParams{
			ID:             userID,
			SynchronizedAt: &startOfYesterday,
		},
	)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncTransactions) getCategoryID(
	categoryExternalID string,
	categoriesByExternalID map[string]entity.TransactionCategory,
) *uuid.UUID {
	parentCategoryExternalID, ok := uc.o.GetParentCategoryExternalID(
		categoryExternalID,
		categoriesByExternalID,
	)
	if !ok {
		return nil
	}

	category, ok := categoriesByExternalID[parentCategoryExternalID]
	if !ok {
		return nil
	}

	return &category.ID
}
