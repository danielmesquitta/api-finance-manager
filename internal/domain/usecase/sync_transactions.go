package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"golang.org/x/sync/errgroup"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/tx"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance/pluggy"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type SyncTransactions struct {
	v   *validator.Validator
	o   openfinance.Client
	tx  tx.TX
	ur  repo.UserRepo
	tr  repo.TransactionRepo
	cr  repo.CategoryRepo
	pmr repo.PaymentMethodRepo
}

func NewSyncTransactions(
	v *validator.Validator,
	o openfinance.Client,
	tx tx.TX,
	ur repo.UserRepo,
	tr repo.TransactionRepo,
	cr repo.CategoryRepo,
	pmr repo.PaymentMethodRepo,
) *SyncTransactions {
	return &SyncTransactions{
		v:   v,
		o:   o,
		tx:  tx,
		ur:  ur,
		tr:  tr,
		cr:  cr,
		pmr: pmr,
	}
}

func (uc *SyncTransactions) Execute(
	ctx context.Context,
) error {
	users, accounts := []entity.User{}, []entity.Account{}
	categories := []entity.Category{}
	paymentMethods := []entity.PaymentMethod{}
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		users, accounts, err = uc.ur.ListPremiumActiveUsersWithAccounts(gCtx)
		return err
	})

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListCategories(gCtx)
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

	if len(users) == 0 || len(accounts) == 0 {
		return nil
	}

	accountsByUserID := uc.groupAccountsByUserID(accounts)

	categoriesByExternalID := uc.groupCategoriesByExternalID(categories)

	openFinanceTransactionsByAccountID := uc.groupOpenFinanceTransactionsByAccountID(
		ctx,
		users,
		accountsByUserID,
	)

	paymentMethods, err := uc.syncPaymentMethods(
		ctx,
		openFinanceTransactionsByAccountID,
		paymentMethods,
	)
	if err != nil {
		return errs.New(err)
	}

	paymentMethodsByExternalID := uc.groupPaymentMethodsByExternalID(
		paymentMethods,
	)

	for _, user := range users {
		accounts := accountsByUserID[user.ID]
		if len(accounts) == 0 {
			continue
		}

		accountsByID := uc.groupAccountsByID(accounts)

		if err := uc.syncUserTransactions(
			ctx,
			user,
			accountsByID,
			categoriesByExternalID,
			paymentMethodsByExternalID,
			openFinanceTransactionsByAccountID,
		); err != nil {
			slog.Error(
				"error syncing user transactions",
				"user", user,
				"accounts", accounts,
				"categories", categories,
				"err", err,
			)
			continue
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

	uniquePaymentMethodExternalIDs := map[string]struct{}{}
	for _, openFinanceTransactions := range openFinanceTransactionsByAccountID {
		for _, openFinanceTransaction := range openFinanceTransactions {
			uniquePaymentMethodExternalIDs[openFinanceTransaction.PaymentMethodExternalID] = struct{}{}
		}
	}

	for _, paymentMethod := range paymentMethods {
		delete(uniquePaymentMethodExternalIDs, paymentMethod.ExternalID)
	}

	params := []repo.CreatePaymentMethodsParams{}
	for paymentMethodExternalID := range uniquePaymentMethodExternalIDs {
		name, ok := paymentMethodNamesByExternalID[paymentMethodExternalID]
		if !ok {
			name = paymentMethodExternalID
		}

		params = append(
			params,
			repo.CreatePaymentMethodsParams{
				Name:       name,
				ExternalID: paymentMethodExternalID,
			},
		)
	}

	if err := uc.pmr.CreatePaymentMethods(ctx, params); err != nil {
		return nil, errs.New(err)
	}

	paymentMethods, err := uc.pmr.ListPaymentMethods(ctx)
	if err != nil {
		return nil, errs.New(err)
	}

	return paymentMethods, nil
}

func (uc *SyncTransactions) groupAccountsByID(
	accounts []entity.Account,
) map[uuid.UUID]entity.Account {
	accountsByID := map[uuid.UUID]entity.Account{}
	for _, account := range accounts {
		accountsByID[account.ID] = account
	}
	return accountsByID
}

func (uc *SyncTransactions) groupCategoriesByExternalID(
	categories []entity.Category,
) map[string]entity.Category {
	categoriesByExternalID := map[string]entity.Category{}
	for _, category := range categories {
		categoriesByExternalID[category.ExternalID] = category
	}
	return categoriesByExternalID
}

func (uc *SyncTransactions) groupPaymentMethodsByExternalID(
	paymentMethods []entity.PaymentMethod,
) map[string]entity.PaymentMethod {
	paymentMethodsByExternalID := map[string]entity.PaymentMethod{}
	for _, paymentMethod := range paymentMethods {
		paymentMethodsByExternalID[paymentMethod.ExternalID] = paymentMethod
	}
	return paymentMethodsByExternalID
}

func (uc *SyncTransactions) groupAccountsByUserID(
	accounts []entity.Account,
) map[uuid.UUID][]entity.Account {
	accountsByUserID := map[uuid.UUID][]entity.Account{}
	for _, account := range accounts {
		accountsByUserID[account.UserID] = append(
			accountsByUserID[account.UserID],
			account,
		)
	}
	return accountsByUserID
}

func (uc *SyncTransactions) groupOpenFinanceTransactionsByAccountID(
	ctx context.Context,
	users []entity.User,
	accountsByUserID map[uuid.UUID][]entity.Account,
) map[uuid.UUID][]openfinance.Transaction {
	openFinanceTransactionsByAccountID := map[uuid.UUID][]openfinance.Transaction{}
	for _, user := range users {
		accounts := accountsByUserID[user.ID]
		if len(accounts) == 0 {
			continue
		}

		lastSynchronizedAt := uc.calculateLastSynchronizedAt(
			user.SynchronizedAt,
		)

		for _, account := range accounts {
			openFinanceTransactions, err := uc.listOpenFinanceTransactions(
				ctx,
				account.ExternalID,
				lastSynchronizedAt,
			)
			if err != nil {
				slog.Error(
					"error getting open finance transactions",
					"user", user,
					"account", account,
					"err", err,
				)
				continue
			}
			openFinanceTransactionsByAccountID[account.ID] = openFinanceTransactions
		}
	}
	return openFinanceTransactionsByAccountID
}

func (uc *SyncTransactions) syncUserTransactions(
	ctx context.Context,
	user entity.User,
	accountsByID map[uuid.UUID]entity.Account,
	categoriesByExternalID map[string]entity.Category,
	paymentMethodsByExternalID map[string]entity.PaymentMethod,
	openFinanceTransactionsByAccountID map[uuid.UUID][]openfinance.Transaction,
) error {
	lastSynchronizedAt := uc.calculateLastSynchronizedAt(user.SynchronizedAt)

	transactions, err := uc.listRepoTransactions(
		ctx,
		user.ID,
		lastSynchronizedAt,
	)
	if err != nil {
		return errs.New(err)
	}

	transactionsByExternalID := uc.groupTransactionsByExternalID(transactions)

	params := uc.buildCreateTransactionsParams(
		user,
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

		if err := uc.updateUserSynchronizedAt(ctx, user); err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncTransactions) buildCreateTransactionsParams(
	user entity.User,
	accountsByID map[uuid.UUID]entity.Account,
	categoriesByExternalID map[string]entity.Category,
	paymentMethodsByExternalID map[string]entity.PaymentMethod,
	openFinanceTransactionsByAccountID map[uuid.UUID][]openfinance.Transaction,
	transactionsByExternalID map[string]entity.Transaction,
) []repo.CreateTransactionsParams {

	params := []repo.CreateTransactionsParams{}
	for accountID, openFinanceTransactions := range openFinanceTransactionsByAccountID {
		account := accountsByID[accountID]

		for _, openFinanceTransaction := range openFinanceTransactions {
			_, transactionAlreadyRegistered := transactionsByExternalID[openFinanceTransaction.ExternalID]
			if transactionAlreadyRegistered {
				continue
			}

			categoryID := uc.getCategoryID(
				openFinanceTransaction.CategoryExternalID,
				categoriesByExternalID,
			)

			paymentMethod := paymentMethodsByExternalID[openFinanceTransaction.PaymentMethodExternalID]

			param := repo.CreateTransactionsParams{
				ExternalID:      openFinanceTransaction.ExternalID,
				Name:            openFinanceTransaction.Name,
				Amount:          openFinanceTransaction.Amount,
				PaymentMethodID: paymentMethod.ID,
				Date:            openFinanceTransaction.Date,
				UserID:          user.ID,
				AccountID:       &account.ID,
				InstitutionID:   &account.InstitutionID,
				CategoryID:      categoryID,
			}

			params = append(params, param)
		}
	}

	return params
}

func (uc *SyncTransactions) calculateLastSynchronizedAt(
	userSynchronizedAt *time.Time,
) time.Time {
	var lastSynchronizedAt time.Time
	if userSynchronizedAt != nil {
		lastSynchronizedAt = getStartOfDay(*userSynchronizedAt)
	}
	return lastSynchronizedAt
}

func (uc *SyncTransactions) listOpenFinanceTransactions(
	ctx context.Context,
	accountExternalID string,
	lastSynchronizedAt time.Time,
) ([]openfinance.Transaction, error) {
	opts := []openfinance.TransactionOption{}
	if !lastSynchronizedAt.IsZero() {
		opts = append(
			opts,
			openfinance.WithTransactionDateAfter(lastSynchronizedAt),
		)
	}

	openFinanceTransactions, err := uc.o.ListTransactions(
		ctx,
		accountExternalID,
		opts...,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	return openFinanceTransactions, nil
}

func (uc *SyncTransactions) groupTransactionsByExternalID(
	transactions []entity.Transaction,
) map[string]entity.Transaction {
	transactionsByExternalID := map[string]entity.Transaction{}
	for _, transaction := range transactions {
		transactionsByExternalID[transaction.ExternalID] = transaction
	}
	return transactionsByExternalID
}

func (uc *SyncTransactions) listRepoTransactions(
	ctx context.Context,
	userID uuid.UUID,
	lastSynchronizedAt time.Time,
) ([]entity.Transaction, error) {
	opts := []repo.TransactionOption{}
	if !lastSynchronizedAt.IsZero() {
		opts = append(
			opts,
			repo.WithTransactionDateAfter(lastSynchronizedAt),
		)
	}

	transactions, err := uc.tr.ListTransactions(
		ctx,
		userID,
		opts...,
	)
	if err != nil {
		return nil, errs.New(err)
	}

	return transactions, nil
}

func (uc *SyncTransactions) updateUserSynchronizedAt(
	ctx context.Context,
	user entity.User,
) error {
	updateUserParams := repo.UpdateUserParams{}
	if err := copier.Copy(&updateUserParams, user); err != nil {
		return errs.New(err)
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	startOfYesterday := getStartOfDay(yesterday)
	updateUserParams.SynchronizedAt = &startOfYesterday

	_, err := uc.ur.UpdateUser(ctx, updateUserParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncTransactions) getCategoryID(
	categoryExternalID string,
	categoriesByExternalID map[string]entity.Category,
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
