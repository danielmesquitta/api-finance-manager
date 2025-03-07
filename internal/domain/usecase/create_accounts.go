package usecase

import (
	"context"
	"log/slog"
	"maps"
	"slices"
	"strconv"

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

type CreateAccounts struct {
	v   *validator.Validator
	o   openfinance.Client
	l   *slog.Logger
	tx  tx.TX
	ur  repo.UserRepo
	ar  repo.AccountRepo
	abr repo.AccountBalanceRepo
	ir  repo.InstitutionRepo
}

func NewCreateAccounts(
	v *validator.Validator,
	o openfinance.Client,
	l *slog.Logger,
	tx tx.TX,
	ur repo.UserRepo,
	ar repo.AccountRepo,
	abr repo.AccountBalanceRepo,
	ir repo.InstitutionRepo,
) *CreateAccounts {
	return &CreateAccounts{
		v:   v,
		o:   o,
		l:   l,
		tx:  tx,
		ur:  ur,
		ar:  ar,
		abr: abr,
		ir:  ir,
	}
}

type CreateAccountsInput struct {
	ItemID          string                    `json:"id"              validate:"required"`
	Institution     CreateAccountsInstitution `json:"connector"       validate:"required"`
	ExecutionStatus string                    `json:"executionStatus" validate:"required"`
	ClientUserID    string                    `json:"clientUserId"    validate:"required,uuid"`
}

type CreateAccountsInstitution struct {
	ID int `json:"id" validate:"required"`
}

func (uc *CreateAccounts) Execute(
	ctx context.Context,
	in CreateAccountsInput,
) error {
	if err := uc.v.Validate(in); err != nil {
		return errs.New(err)
	}

	if in.ExecutionStatus != "SUCCESS" {
		uc.l.Info(
			"sync-accounts: execution status is not SUCCESS",
			"execution_status", in.ExecutionStatus,
		)
		return nil
	}

	userID := uuid.MustParse(in.ClientUserID)

	var (
		user                *entity.User
		institution         *entity.Institution
		openFinanceAccounts []openfinance.Account
	)
	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		user, err = uc.ur.GetUserByID(gCtx, userID)
		return err
	})

	g.Go(func() error {
		var err error
		institution, err = uc.ir.GetInstitutionByExternalID(
			gCtx,
			strconv.Itoa(in.Institution.ID),
		)
		return err
	})

	g.Go(func() error {
		var err error
		openFinanceAccounts, err = uc.o.ListAccounts(gCtx, in.ItemID)
		return err
	})

	if err := g.Wait(); err != nil {
		return errs.New(err)
	}

	if user == nil {
		return errs.ErrUserNotFound
	}

	if institution == nil {
		return errs.ErrInstitutionNotFound
	}

	if len(openFinanceAccounts) == 0 {
		return errs.ErrOpenFinanceAccountsNotFound
	}

	createAccountsParamsByExternalID := map[string][]repo.CreateAccountsParams{}
	for _, account := range openFinanceAccounts {
		param := repo.CreateAccountsParams{}
		if err := copier.Copy(&param, account); err != nil {
			return errs.New(err)
		}

		param.ID = uuid.New()
		param.UserID = user.ID
		param.InstitutionID = institution.ID
		createAccountsParamsByExternalID[account.ExternalID] = append(
			createAccountsParamsByExternalID[account.ExternalID],
			param,
		)
	}

	accountExternalIDs := slices.Collect(
		maps.Keys(createAccountsParamsByExternalID),
	)

	registeredAccounts, err := uc.ar.ListAccounts(
		ctx,
		repo.WithAccountExternalIDs(accountExternalIDs),
	)
	if err != nil {
		return errs.New(err)
	}

	if len(registeredAccounts) == len(accountExternalIDs) {
		return nil
	}

	if len(registeredAccounts) > 0 {
		for _, registeredAccount := range registeredAccounts {
			delete(
				createAccountsParamsByExternalID,
				registeredAccount.ExternalID,
			)
		}
	}

	createAccountsParams := []repo.CreateAccountsParams{}
	for _, params := range createAccountsParamsByExternalID {
		createAccountsParams = append(createAccountsParams, params...)
	}

	err = uc.tx.Do(ctx, func(ctx context.Context) error {
		if err := uc.ar.CreateAccounts(ctx, createAccountsParams); err != nil {
			return errs.New(err)
		}

		accountIDByExternalID := map[string]uuid.UUID{}
		for _, account := range createAccountsParams {
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
