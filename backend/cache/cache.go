package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(host string, port int64) *Cache {
	return &Cache{
		rdb: redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", host, port),
			Password: "",
			DB:       0,
		}),
	}
}

func (c *Cache) GetNextId(ctx context.Context, key string, minimum int64) (int64, error) {
	pipe := c.rdb.Pipeline()
	pipe.SetNX(ctx, key, minimum, 0)
	incr := pipe.Incr(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Result()
}
