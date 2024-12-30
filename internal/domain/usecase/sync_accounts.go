package usecase

import (
	"context"
	"fmt"
	"log"

	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/errs"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/openfinance"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/repo"
	"github.com/jinzhu/copier"
)

type SyncAccounts struct {
	v  *validator.Validator
	o  openfinance.Client
	ur repo.UserRepo
	ar repo.AccountRepo
	ir repo.InstitutionRepo
}

func NewSyncAccounts(
	v *validator.Validator,
	o openfinance.Client,
	ur repo.UserRepo,
	ar repo.AccountRepo,
	ir repo.InstitutionRepo,
) *SyncAccounts {
	return &SyncAccounts{
		v:  v,
		o:  o,
		ur: ur,
		ar: ar,
		ir: ir,
	}
}

type SyncAccountsInput struct {
	ItemID          string      `json:"id"              validate:"required"`
	Institution     Institution `json:"connector"       validate:"required"`
	ExecutionStatus string      `json:"executionStatus" validate:"required"`
	ClientUserID    string      `json:"clientUserId"    validate:"required"`
}

type Institution struct {
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
		log.Println(
			"sync-accounts: execution status is not SUCCESS:",
			in.ExecutionStatus,
		)
		return nil
	}

	var institution *entity.Institution
	var accounts []entity.Account
	var users []entity.User
	errCh := make(chan error, 3)
	defer close(errCh)

	institutionExternalID := fmt.Sprintf("%d", in.Institution.ID)
	go func() {
		var err error
		institution, err = uc.ir.GetInstitutionByExternalID(
			ctx,
			institutionExternalID,
		)
		errCh <- err
	}()

	go func() {
		var err error
		accounts, err = uc.o.ListAccounts(ctx, in.ItemID)
		errCh <- err
	}()

	go func() {
		var err error
		users, err = uc.ur.ListUsers(ctx)
		errCh <- err
	}()

	for i := 0; i < cap(errCh); i++ {
		if err := <-errCh; err != nil {
			return errs.New(err)
		}
	}

	if institution == nil {
		return errs.New("institution not found")
	}

	params := []repo.CreateAccountsParams{}
	for _, user := range users {
		for _, account := range accounts {
			param := repo.CreateAccountsParams{}
			if err := copier.Copy(&param, account); err != nil {
				return errs.New(err)
			}
			param.UserID = user.ID
			param.InstitutionID = institution.ID
			params = append(params, param)
		}
	}

	if err := uc.ar.CreateAccounts(ctx, params); err != nil {
		return errs.New(err)
	}

	return nil
}
