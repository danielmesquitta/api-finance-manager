package server

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/server"
	"github.com/danielmesquitta/api-finance-manager/internal/app/server/dto"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/config/env"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/entity"
	"github.com/danielmesquitta/api-finance-manager/internal/domain/usecase/auth"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/test/container"
	"github.com/danielmesquitta/api-finance-manager/test/db"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

var (
	wg sync.WaitGroup
	ev *env.Env
	vl *validator.Validator
)

type TestApp struct {
	t  *testing.T
	a  *server.App
	db *db.TestDB
}

func init() {
	wg.Add(1)
	defer wg.Done()

	vl = validator.New()
	ev = config.LoadConfig(vl)
}

func NewTestApp(
	t *testing.T,
) (app *TestApp, cleanUp func(context.Context) error) {
	wg.Wait()
	v := *vl
	e := *ev

	mu := sync.Mutex{}
	g, gCtx := errgroup.WithContext(context.Background())
	cleanUps := []func(context.Context) error{}

	var postgresDatabaseURL string
	g.Go(func() error {
		connectionString, cleanUp := container.NewPostgresContainer(gCtx)

		mu.Lock()
		postgresDatabaseURL = connectionString
		cleanUps = append(cleanUps, cleanUp)
		mu.Unlock()

		return nil
	})

	var redisDatabaseURL string
	g.Go(func() error {
		connectionString, cleanUp := container.NewRedisContainer(gCtx)

		mu.Lock()
		redisDatabaseURL = connectionString
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

	e.PostgresDatabaseURL = postgresDatabaseURL
	e.RedisDatabaseURL = redisDatabaseURL

	restAPI := server.NewTest(&v, &e)

	app = &TestApp{
		t:  t,
		a:  restAPI,
		db: db.NewTestDB(restAPI.DB),
	}

	return app, cleanUp
}

// RequestOption represents an option for the MakeRequest function
type RequestOption func(*requestOptions)

type requestOptions struct {
	body          any
	token         string
	bearerToken   string
	headers       map[string]string
	queryParams   map[string]string
	response      any
	errorResponse any
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
func WithBearerToken(bearerToken string) RequestOption {
	return func(o *requestOptions) {
		o.bearerToken = bearerToken
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

// WithError sets the error response object to unmarshal into
func WithError(errorResponse any) RequestOption {
	return func(o *requestOptions) {
		o.errorResponse = errorResponse
	}
}

// WithQueryParams sets the query parameters for the request
func WithQueryParams(queryParams map[string]string) RequestOption {
	return func(o *requestOptions) {
		o.queryParams = queryParams
	}
}

func (ta *TestApp) MakeRequest(
	method string,
	path string,
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

	u := &url.URL{Path: path}
	if len(options.queryParams) > 0 {
		v := url.Values{}
		for k, val := range options.queryParams {
			v.Add(k, val)
		}
		u.RawQuery = v.Encode()
	}

	req, _ := http.NewRequest(
		method,
		u.String(),
		bytes.NewReader(jsonBody),
	)

	req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
	req.Header.Set(fiber.HeaderAccept, fiber.MIMEApplicationJSON)

	if options.token != "" {
		req.Header.Set(fiber.HeaderAuthorization, options.token)
	}

	if options.bearerToken != "" {
		req.Header.Set(fiber.HeaderAuthorization, "Bearer "+options.bearerToken)
	}

	for k, v := range options.headers {
		req.Header.Set(k, v)
	}

	res, err := ta.a.Test(req, -1)
	assert.Nil(ta.t, err)

	bytesBody, _ := io.ReadAll(res.Body)
	if len(bytesBody) == 0 {
		return res.StatusCode, "", nil
	}

	if options.response != nil && res.StatusCode >= 200 &&
		res.StatusCode < 300 {
		_ = json.Unmarshal(bytesBody, options.response)
	} else if options.errorResponse != nil && res.StatusCode >= 400 {
		_ = json.Unmarshal(bytesBody, options.errorResponse)
	}

	return res.StatusCode, string(bytesBody), nil
}

func (ta *TestApp) SignIn(token string) *dto.SignInResponse {
	body := dto.SignInRequest{
		SignInUseCaseInput: auth.SignInUseCaseInput{
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

	return &out
}
