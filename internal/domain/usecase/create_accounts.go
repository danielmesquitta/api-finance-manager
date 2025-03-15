package usecase

import (
	"context"
	"log/slog"
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
	tx  tx.TX
	ar  repo.AccountRepo
	abr repo.AccountBalanceRepo
	ir  repo.InstitutionRepo
	uir repo.UserInstitutionRepo
}

func NewCreateAccounts(
	v *validator.Validator,
	o openfinance.Client,
	tx tx.TX,
	ar repo.AccountRepo,
	abr repo.AccountBalanceRepo,
	ir repo.InstitutionRepo,
	uir repo.UserInstitutionRepo,
) *CreateAccounts {
	return &CreateAccounts{
		v:   v,
		o:   o,
		tx:  tx,
		ar:  ar,
		abr: abr,
		ir:  ir,
		uir: uir,
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
		slog.Info(
			"sync-accounts: execution status is not SUCCESS",
			"execution_status", in.ExecutionStatus,
		)
		return nil
	}

	userID := uuid.MustParse(in.ClientUserID)

	var (
		institution         *entity.Institution
		userInstitution     *entity.UserInstitution
		openFinanceAccounts []openfinance.Account
	)
	g, gCtx := errgroup.WithContext(ctx)

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
		userInstitution, err = uc.uir.GetUserInstitutionByExternalID(
			gCtx,
			in.ItemID,
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

	if institution == nil {
		return errs.ErrInstitutionNotFound
	}

	if len(openFinanceAccounts) == 0 {
		return errs.ErrOpenFinanceAccountsNotFound
	}

	accountExternalIDs := []string{}
	accountsByExternalID := map[string]openfinance.Account{}
	for _, account := range openFinanceAccounts {
		accountsByExternalID[account.ExternalID] = account
		accountExternalIDs = append(accountExternalIDs, account.ExternalID)
	}

	if userInstitution != nil {
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
					accountsByExternalID,
					registeredAccount.ExternalID,
				)
			}
		}
	}

	err := uc.tx.Do(ctx, func(ctx context.Context) error {
		if userInstitution == nil {
			userInstitutionParams := repo.CreateUserInstitutionParams{
				UserID:        userID,
				InstitutionID: institution.ID,
				ExternalID:    in.ItemID,
			}

			var err error
			userInstitution, err = uc.uir.CreateUserInstitution(
				ctx,
				userInstitutionParams,
			)
			if err != nil {
				return errs.New(err)
			}
		}

		createAccountsParams := []repo.CreateAccountsParams{}
		for _, account := range accountsByExternalID {
			params := repo.CreateAccountsParams{}
			if err := copier.Copy(&params, account); err != nil {
				return errs.New(err)
			}
			params.ID = uuid.New()
			params.UserInstitutionID = userInstitution.ID
			createAccountsParams = append(
				createAccountsParams,
				params,
			)
		}

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
				UserID:    userID,
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

		return nil
	})
	if err != nil {
		return errs.New(err)
	}

	return nil
}
