package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"

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

	signInRes := dto.SignInResponse{}

	res, err := client.R().
		SetHeader(fiber.HeaderAuthorization, mockoauth.DefaultMockToken).
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
		Post("/v1/admin/transactions/categories/sync")
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

	jsonData, err := root.TestData.ReadFile("testdata/pluggy/items.json")
	if err != nil {
		panic(err)
	}

	createAccountsReqs := []dto.CreateAccountsRequest{}
	if err := json.Unmarshal(jsonData, &createAccountsReqs); err != nil {
		panic(err)
	}

	for i := range createAccountsReqs {
		createAccountsReqs[i].ClientUserID = signInRes.User.ID.String()
	}

	createAccountsReqs = createAccountsReqs[:2]

	for _, createAccountsReq := range createAccountsReqs {
		res, err := client.R().
			SetBody(createAccountsReq).
			Post("/v1/admin/accounts")
		if err != nil {
			panic(err)
		}
		if res.IsError() {
			panic(string(res.Body()))
		}
	}

	res, err = client.R().
		SetQueryParam("user_ids", signInRes.User.ID.String()).
		Post("/v1/admin/transactions/sync")
	if err != nil {
		panic(err)
	}
	if res.IsError() {
		panic(string(res.Body()))
	}
}
