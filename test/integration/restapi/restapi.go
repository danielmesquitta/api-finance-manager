package restapi

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/test/container"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

type TestApp struct {
	t   *testing.T
	App *restapi.App
}

func NewTestApp(
	t *testing.T,
) (app *TestApp, cleanUp func(context.Context) error) {
	v := validator.New()
	e := config.LoadConfig(v)

	mu := sync.Mutex{}
	g, gCtx := errgroup.WithContext(context.Background())
	cleanUps := []func(context.Context) error{}

	g.Go(func() error {
		connectionString, cleanUp := container.NewPostgresContainer(gCtx)

		mu.Lock()
		e.PostgresDatabaseURL = connectionString
		cleanUps = append(cleanUps, cleanUp)
		mu.Unlock()

		return nil
	})

	g.Go(func() error {
		connectionString, cleanUp := container.NewRedisContainer(gCtx)

		mu.Lock()
		e.RedisDatabaseURL = connectionString
		cleanUps = append(cleanUps, cleanUp)
		mu.Unlock()

		return nil
	})

	if err := g.Wait(); err != nil {
		panic(err)
	}

	cleanUp = func(ctx context.Context) error {
		for _, c := range cleanUps {
			if err := c(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	app = &TestApp{
		t:   t,
		App: restapi.NewTest(v, e, t),
	}

	return app, cleanUp
}

func (ta *TestApp) SignIn(token string) (accessToken string) {
	body := dto.SignInRequest{
		SignInInput: usecase.SignInInput{
			Provider: entity.ProviderMock,
		},
	}

	var out dto.SignInResponse
	statusCode, _, err := ta.MakeRequest(
		http.MethodPost,
		"/api/v1/auth/sign-in",
		WithBody(body),
		WithToken(token),
		WithResponse(&out),
	)
	assert.Nil(ta.t, err)
	assert.Equal(ta.t, http.StatusOK, statusCode)

	return out.AccessToken
}

// RequestOption represents an option for the MakeRequest function
type RequestOption func(*requestOptions)

type requestOptions struct {
	body        any
	token       string
	bearerToken string
	headers     map[string]string
	response    any
}

// WithBody sets the body for the request
func WithBody(body any) RequestOption {
	return func(o *requestOptions) {
		o.body = body
	}
}

// WithToken sets the authorization token for the request
func WithToken(token string) RequestOption {
	return func(o *requestOptions) {
		o.token = token
	}
}

// WithToken sets the authorization token for the request
func WithBearerToken(token string) RequestOption {
	return func(o *requestOptions) {
		o.token = token
	}
}

// WithHeaders sets additional headers for the request
func WithHeaders(headers map[string]string) RequestOption {
	return func(o *requestOptions) {
		o.headers = headers
	}
}

// WithResponse sets the response object to unmarshal into
func WithResponse(response any) RequestOption {
	return func(o *requestOptions) {
		o.response = response
	}
}

func (ta *TestApp) MakeRequest(
	method string,
	url string,
	opts ...RequestOption,
) (statusCode int, rawBody string, err error) {
	options := &requestOptions{}
	for _, opt := range opts {
		opt(options)
	}

	var jsonBody []byte
	if options.body != nil {
		jsonBody, err = json.Marshal(options.body)
		if err != nil {
			assert.Nil(ta.t, err)
			return 0, "", err
		}
	}

	req, _ := http.NewRequest(
		method,
		url,
		bytes.NewReader(jsonBody),
	)

	req.Header.Set("Content-Type", "application/json")

	if options.token != "" {
		req.Header.Set("Authorization", options.token)
	}

	if options.bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+options.bearerToken)
	}

	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	res, err := ta.App.Test(req, -1)
	assert.Nil(ta.t, err)

	bytesBody, err := io.ReadAll(res.Body)
	assert.Nil(ta.t, err)

	if options.response != nil {
		err = json.Unmarshal(bytesBody, options.response)
		assert.Nil(ta.t, err)
	}

	return res.StatusCode, string(bytesBody), nil
}
