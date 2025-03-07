package container

import (
	"context"

	"github.com/testcontainers/testcontainers-go/modules/redis"
)

func NewRedisContainer(
	ctx context.Context,
) (connectionString string, cleanUp func(context.Context) error) {
	redisCont, err := redis.Run(ctx, "redis:alpine")
	if err != nil {
		panic(err)
	}
	connStr, err := redisCont.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	cleanUp = func(ctx context.Context) error {
		return redisCont.Terminate(ctx)
	}

	return connStr, cleanUp
}
