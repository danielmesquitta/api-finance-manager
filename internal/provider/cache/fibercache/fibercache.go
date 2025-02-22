package fibercache

import (
	"context"
	"time"

	"github.com/danielmesquitta/api-finance-manager/internal/provider/cache"
	"github.com/gofiber/fiber/v2"
)

type FiberCache struct {
	c cache.Cache
}

func NewFiberCache(c cache.Cache) *FiberCache {
	return &FiberCache{
		c: c,
	}
}

func (f *FiberCache) Get(key string) ([]byte, error) {
	var val []byte
	ok, err := f.c.Scan(context.Background(), cache.Key(key), &val)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (f *FiberCache) Set(key string, val []byte, exp time.Duration) error {
	return f.c.Set(context.Background(), cache.Key(key), val, exp)
}

func (f *FiberCache) Delete(key string) error {
	return f.c.Delete(context.Background(), cache.Key(key))
}

func (f *FiberCache) Reset() error {
	return nil
}

func (f *FiberCache) Close() error {
	return nil
}

var _ fiber.Storage = (*FiberCache)(nil)
