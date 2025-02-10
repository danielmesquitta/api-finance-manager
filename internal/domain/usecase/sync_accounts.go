package usecase

import (
	"context"
	"fmt"
	"log/slog"

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

type SyncAccounts struct {
	v   *validator.Validator
	o   openfinance.Client
	tx  tx.TX
	ur  repo.UserRepo
	ar  repo.AccountRepo
	abr repo.AccountBalanceRepo
	ir  repo.InstitutionRepo
}

func NewSyncAccounts(
	v *validator.Validator,
	o openfinance.Client,
	tx tx.TX,
	ur repo.UserRepo,
	ar repo.AccountRepo,
	abr repo.AccountBalanceRepo,
	ir repo.InstitutionRepo,
) *SyncAccounts {
	return &SyncAccounts{
		v:   v,
		o:   o,
		tx:  tx,
		ur:  ur,
		ar:  ar,
		abr: abr,
		ir:  ir,
	}
}

type SyncAccountsInput struct {
	ItemID          string                  `json:"id"              validate:"required"`
	Institution     SyncAccountsInstitution `json:"connector"       validate:"required"`
	ExecutionStatus string                  `json:"executionStatus" validate:"required"`
	ClientUserID    uuid.UUID               `json:"clientUserId"    validate:"required"`
}

type SyncAccountsInstitution struct {
	ID int `json:"id" validate:"required"`
}

func (uc *SyncAccounts) Execute(
	ctx context.Context,
	in SyncAccountsInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	if in.ExecutionStatus != "SUCCESS" {
		slog.Info(
			"sync-accounts: execution status is not SUCCESS",
			"execution_status", in.ExecutionStatus,
		)
		return nil
	}

	var institution *entity.Institution
	var openFinanceAccounts []openfinance.Account
	var user *entity.User
	g, gCtx := errgroup.WithContext(ctx)

	institutionExternalID := fmt.Sprintf("%d", in.Institution.ID)
	g.Go(func() error {
		var err error
		institution, err = uc.ir.GetInstitutionByExternalID(
			gCtx,
			institutionExternalID,
		)
		return err
	})

	g.Go(func() error {
		var err error
		openFinanceAccounts, err = uc.o.ListAccounts(gCtx, in.ItemID)
		return err
	})

	g.Go(func() error {
		var err error
		user, err = uc.ur.GetUserByID(gCtx, in.ClientUserID)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if len(openFinanceAccounts) == 0 {
		return errs.New("no open finance accounts found")
	}

	if institution == nil {
		return errs.ErrInstitutionNotFound
	}

	if user == nil {
		return errs.ErrUserNotFound
	}

	createAccountsParams := []repo.CreateAccountsParams{}
	accountExternalIDs := []string{}
	for _, account := range openFinanceAccounts {
		accountExternalIDs = append(accountExternalIDs, account.ExternalID)

		param := repo.CreateAccountsParams{}
		if err := copier.Copy(&param, account); err != nil {
			return errs.New(err)
		}

		param.UserID = user.ID
		param.InstitutionID = institution.ID
		createAccountsParams = append(createAccountsParams, param)
	}

	accounts, err := uc.ar.ListAccountsByExternalIDs(
		ctx,
		accountExternalIDs,
	)
	if err != nil {
		return errs.New(err)
	}

	if len(accounts) == len(accountExternalIDs) {
		return nil
	}

	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if err := uc.ar.CreateAccounts(ctx, createAccountsParams); err != nil {
			return errs.New(err)
		}

		accounts, err = uc.ar.ListAccountsByExternalIDs(
			ctx,
			accountExternalIDs,
		)
		if err != nil {
			return errs.New(err)
		}

		if len(accountExternalIDs) != len(accounts) {
			return errs.New("failed to sync accounts")
		}

		accountIDByExternalID := map[string]uuid.UUID{}
		for _, account := range accounts {
			accountIDByExternalID[account.ExternalID] = account.ID
		}

		createAccountBalancesParams := []repo.CreateAccountBalancesParams{}
		for _, account := range openFinanceAccounts {
			accountID := accountIDByExternalID[account.ExternalID]
			accountBalance := repo.CreateAccountBalancesParams{
				UserID:    user.ID,
				AccountID: accountID,
				Amount:    account.Balance,
			}
			createAccountBalancesParams = append(
				createAccountBalancesParams,
				accountBalance,
			)
		}

		if err := uc.abr.CreateAccountBalances(ctx, createAccountBalancesParams); err != nil {
			return errs.New(err)
		}

		updateUserParams := repo.UpdateUserParams{}
		if err := copier.Copy(&updateUserParams, user); err != nil {
			return errs.New(err)
		}
		updateUserParams.OpenFinanceID = &in.ItemID

		if _, err := uc.ur.UpdateUser(ctx, updateUserParams); err != nil {
			return errs.New(err)
		}

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}
