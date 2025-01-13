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
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
)

type SyncTransactions struct {
	v  *validator.Validator
	o  openfinance.Client
	tx tx.TX
	ur repo.UserRepo
	tr repo.TransactionRepo
	cr repo.CategoryRepo
}

func NewSyncTransactions(
	v *validator.Validator,
	o openfinance.Client,
	tx tx.TX,
	ur repo.UserRepo,
	tr repo.TransactionRepo,
	cr repo.CategoryRepo,
) *SyncTransactions {
	return &SyncTransactions{
		v:  v,
		o:  o,
		tx: tx,
		ur: ur,
		tr: tr,
		cr: cr,
	}
}

func (uc *SyncTransactions) Execute(
	ctx context.Context,
) error {
	users, accounts := []entity.User{}, []entity.Account{}
	categories := []entity.Category{}
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		users, accounts, err = uc.ur.ListPremiumActiveUsersWithAccounts(ctx)
		return err
	})

	g.Go(func() error {
		var err error
		categories, err = uc.cr.ListCategories(ctx)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if len(users) == 0 || len(accounts) == 0 {
		return nil
	}

	categoriesByExternalID := uc.getCategoriesByExternalID(categories)

	accountsByUserID := uc.getAccountsByUserID(accounts)

	for _, user := range users {
		accounts := accountsByUserID[user.ID]
		if len(accounts) == 0 {
			continue
		}

		if err := uc.syncUserTransactions(
			ctx,
			user,
			accounts,
			categoriesByExternalID,
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

func (uc *SyncTransactions) getCategoriesByExternalID(
	categories []entity.Category,
) map[string]entity.Category {
	categoriesByExternalID := map[string]entity.Category{}
	for _, category := range categories {
		categoriesByExternalID[category.ExternalID] = category
	}
	return categoriesByExternalID
}

func (uc *SyncTransactions) getAccountsByUserID(
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

func (uc *SyncTransactions) syncUserTransactions(
	ctx context.Context,
	user entity.User,
	accounts []entity.Account,
	categoriesByExternalID map[string]entity.Category,
) error {
	var lastSynchronizedAt time.Time
	if user.SynchronizedAt != nil {
		lastSynchronizedAt = uc.getStartOfDay(*user.SynchronizedAt)
	}

	transactions, err := uc.getRepoTransactions(
		ctx,
		user.ID,
		lastSynchronizedAt,
	)
	if err != nil {
		return errs.New(err)
	}

	transactionsByExternalID := uc.getTransactionsByExternalID(transactions)

	params := []repo.CreateTransactionsParams{}
	for _, account := range accounts {
		openFinanceTransactions, err := uc.getOpenFinanceTransactions(
			ctx,
			account.ExternalID,
			lastSynchronizedAt,
		)
		if err != nil {
			return errs.New(err)
		}
		if len(openFinanceTransactions) == 0 {
			continue
		}

		uc.setCreateTransactionsParams(
			&params,
			user,
			account,
			openFinanceTransactions,
			transactionsByExternalID,
			categoriesByExternalID,
		)
	}

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

func (uc *SyncTransactions) setCreateTransactionsParams(
	params *[]repo.CreateTransactionsParams,
	user entity.User,
	account entity.Account,
	openFinanceTransactions []openfinance.Transaction,
	transactionsByExternalID map[string]entity.Transaction,
	categoriesByExternalID map[string]entity.Category,
) {
	for _, openFinanceTransaction := range openFinanceTransactions {
		_, transactionAlreadyRegistered := transactionsByExternalID[openFinanceTransaction.ExternalID]
		if transactionAlreadyRegistered {
			continue
		}

		categoryID := uc.getCategoryID(
			openFinanceTransaction.CategoryExternalID,
			categoriesByExternalID,
		)

		param := repo.CreateTransactionsParams{
			ExternalID:    openFinanceTransaction.ExternalID,
			Name:          openFinanceTransaction.Name,
			Amount:        openFinanceTransaction.Amount,
			PaymentMethod: openFinanceTransaction.PaymentMethod,
			Date:          openFinanceTransaction.Date,
			UserID:        user.ID,
			AccountID:     &account.ID,
			InstitutionID: &account.InstitutionID,
			CategoryID:    categoryID,
		}

		*params = append(*params, param)
	}
}

func (uc *SyncTransactions) getOpenFinanceTransactions(
	ctx context.Context,
	accountExternalID string,
	lastSynchronizedAt time.Time,
) ([]openfinance.Transaction, error) {
	opts := []openfinance.ListTransactionsOption{}
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

func (uc *SyncTransactions) getTransactionsByExternalID(
	transactions []entity.Transaction,
) map[string]entity.Transaction {
	transactionsByExternalID := map[string]entity.Transaction{}
	for _, transaction := range transactions {
		transactionsByExternalID[transaction.ExternalID] = transaction
	}
	return transactionsByExternalID
}

func (uc *SyncTransactions) getRepoTransactions(
	ctx context.Context,
	userID uuid.UUID,
	lastSynchronizedAt time.Time,
) ([]entity.Transaction, error) {
	opts := []repo.ListTransactionsOption{}
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
	startOfYesterday := uc.getStartOfDay(yesterday)
	updateUserParams.SynchronizedAt = &startOfYesterday

	_, err := uc.ur.UpdateUser(ctx, updateUserParams)
	if err != nil {
		return errs.New(err)
	}

	return nil
}

func (uc *SyncTransactions) getStartOfDay(date time.Time) time.Time {
	startOfDay := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0,
		0,
		0,
		0,
		time.Local,
	)

	return startOfDay
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
