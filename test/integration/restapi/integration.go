package restapi

import (
	"context"
	"sync"
	"testing"

	"github.com/danielmesquitta/api-finance-manager/internal/app/restapi"
	"github.com/danielmesquitta/api-finance-manager/internal/config"
	"github.com/danielmesquitta/api-finance-manager/internal/pkg/validator"
	"github.com/danielmesquitta/api-finance-manager/test/container"
	"golang.org/x/sync/errgroup"
)

func Setup(
	t *testing.T,
) (app *restapi.App, cleanUp func(context.Context) error) {
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

	app = restapi.NewTest(v, e, t)

	return app, cleanUp
}
