package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	shortCodeExpiration       = 10
	shortCodeKey              = "short_code:"
	shortCodeIdKey            = "short-code-id"
	minimumShortCodeId  int64 = 1000
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

func (c *Cache) GetNextId(ctx context.Context) (int64, error) {
	pipe := c.rdb.Pipeline()
	pipe.SetNX(ctx, shortCodeIdKey, minimumShortCodeId, 0)
	incr := pipe.Incr(ctx, shortCodeIdKey)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return incr.Result()
}

// cached short_code -> long_url
func (c *Cache) GetByShortCode(ctx context.Context, shortCode string) (string, error) {
	long_url, err := c.rdb.Get(ctx, shortCodeKey+shortCode).Result()
	switch {
	case err == redis.Nil:
		return "", nil
	case err != nil:
		return "", err
	default:
		return long_url, nil
	}
}

func (c *Cache) PutShortCode(ctx context.Context, shortCode string, longUrl string) error {
	_, err := c.rdb.SetEx(ctx, shortCodeKey+shortCode, longUrl, time.Duration(shortCodeExpiration)*time.Minute).Result()
	if err != nil {
		return err
	}
	return nil
}
