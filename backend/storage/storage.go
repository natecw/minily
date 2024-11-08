package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natecw/minily/cache"
	"github.com/natecw/minily/encoding"
	"github.com/natecw/minily/models"
)

const ()

type Storage struct {
	pool  *pgxpool.Pool
	cache *cache.Cache
}

func NewStorage(databaseUrl string, cache *cache.Cache) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}
	return &Storage{
		pool:  pool,
		cache: cache,
	}, nil
}

func (s *Storage) CreateMinily(ctx context.Context, m models.CreateRequest) (*models.Minily, error) {
	var code string
	s.pool.QueryRow(ctx, "select short_code from urls where long_url = $1 and coalesce(expiration, now() + '1 years'::interval) > $2", m.URL, time.Now()).
		Scan(&code)
	fmt.Fprintf(os.Stdout, "short_code: %v\n", code)
	if code != "" {
		return &models.Minily{
			ShortCode: code,
		}, nil
	}

	// todo: handle alias
	nextId, err := s.cache.GetNextId(ctx)
	fmt.Fprintf(os.Stdout, "nextId: %v\n", nextId)
	if err != nil {
		return nil, err
	}
	short_code := encoding.Encode(nextId)

	fmt.Fprintf(os.Stdout, "encoded short_code: %v\n", short_code)
	s.pool.QueryRow(ctx, "insert into urls(short_code, long_url, alias, expiration, created_by) values($1, $2, $3, $4, $5)",
		short_code, m.URL, m.Alias, m.ExpiresAt, m.CreatedBy)
	return &models.Minily{
		ShortCode: short_code,
	}, nil
}

func (s *Storage) GetOriginalUrl(ctx context.Context, short_code string) (string, error) {
	url, err := s.cache.GetByShortCode(ctx, short_code)
	fmt.Printf("short_code(%s)->url(%s) or error(%v)\n", short_code, url, err)
	if err != nil {
		return "", err
	}
	if url != "" {
		return url, nil
	}
	var long_url string
	// todo: not deleting expirations so cleanup required at some point
	s.pool.QueryRow(ctx, "select long_url from urls where short_code=$1 and coalesce(expiration, now() + '1 years'::interval) > $2", short_code, time.Now()).
		Scan(&long_url)
	if long_url == "" {
		return "", fmt.Errorf("unknown url %v", short_code)
	}

	fmt.Println("after query", long_url)
	go func() {
		s.cache.PutShortCode(ctx, short_code, long_url)
		fmt.Println("stored", short_code, "->", long_url, " in cache")
	}()
	return long_url, nil
}
