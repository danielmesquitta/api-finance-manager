package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"golang.org/x/sync/errgroup"

	root "github.com/danielmesquitta/api-finance-manager"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/internal/provider/oauth/mockoauth"
)

func main() {
	v := validator.New()
	e := config.LoadConfig(v)

	var baseURL string
	if e.Port == "" || e.Port == "80" {
		baseURL = fmt.Sprintf("%s/api", e.Host)
	} else {
		baseURL = fmt.Sprintf("%s:%s/api", e.Host, e.Port)
	}

	client := resty.New().
		SetBaseURL(baseURL).
		SetDebug(true)

	res, err := client.R().Get("/health")
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(string(res.Body()))
	}

	signInRes := dto.SignInResponse{}

	res, err = client.R().
		SetHeader("Authorization", mockoauth.MockToken).
		SetBody(dto.SignInRequest{SignInInput: usecase.SignInInput{
			Provider: entity.ProviderMock,
		}}).
		SetResult(&signInRes).
		Post("/v1/auth/sign-in")
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(string(res.Body()))
	}

	client.
		SetBasicAuth(e.BasicAuthUsername, e.BasicAuthPassword)

	res, err = client.R().
		Post("/v1/admin/categories/sync")
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(string(res.Body()))
	}

	res, err = client.R().
		Post("/v1/admin/institutions/sync")
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(string(res.Body()))
	}

	data, err := root.TestData.ReadFile("test/data/pluggy/items.json")
	if err != nil {
		panic(err)
	}

	syncAccountsReqs := []dto.CreateAccountsRequest{}
	if err := json.Unmarshal(data, &syncAccountsReqs); err != nil {
		panic(err)
	}

	for i := range syncAccountsReqs {
		syncAccountsReqs[i].ClientUserID = signInRes.User.ID
	}

	g, ctx := errgroup.WithContext(context.Background())

	for _, syncAccountsReq := range syncAccountsReqs {
		g.Go(func() error {
			res, err := client.R().
				SetContext(ctx).
				SetBody(syncAccountsReq).
				Post("/v1/admin/accounts/sync")
			if err != nil {
				return err
			}
			if res.IsError() {
				panic(string(res.Body()))
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		panic(err)
	}
}
