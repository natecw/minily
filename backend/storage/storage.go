package storage

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/natecw/minily/encoding"
	"github.com/natecw/minily/models"
)

type Storage struct {
	pool *pgxpool.Pool
}

func NewStorage(databaseUrl string) (*Storage, error) {
	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		return nil, fmt.Errorf("could not open db: %w", err)
	}
	return &Storage{
		pool: pool,
	}, nil
}

// basic md5 and just insert. not handling duplicates; resulting short_code could be longer
func (s *Storage) CreateMinily(ctx context.Context, m models.CreateRequest) (*models.Minily, error) {
	var code string
	s.pool.QueryRow(ctx, "select short_code from urls where long_url = $1 and coalesce(expiration, now() + '1 years'::interval) > $2", m.URL, time.Now()).
		Scan(&code)
	if code == "" {
		return &models.Minily{
			ShortCode: code,
		}, nil
	}
	// todo: if alias not given then encode long_url to shorter; else use alias
	// todo: when generating path, refactor to redis incr
	short_code := encoding.EncodeMd5(m.URL)

	// todo: if alias used check it's not already in db, if it is and long_url is different then error; else then guchi

	s.pool.QueryRow(ctx, "insert into urls(short_code, long_url, alias, expiration, created_by) values($1, $2, $3, $4, $5)",
		short_code, m.URL, m.Alias, m.ExpiresAt, m.CreatedBy)
	return &models.Minily{
		ShortCode: short_code,
	}, nil
}

func (s *Storage) GetUrl(ctx context.Context, short_code string) (string, error) {
	// todo: not deleting expirations so cleanup required at some point
	var long_url string
	s.pool.QueryRow(ctx, "select long_url from urls where short_code=$1 and coalesce(expiration, now() + '1 years'::interval) > $2", short_code, time.Now()).
		Scan(&long_url)
	if long_url == "" {
		return "", fmt.Errorf("unknown url %v", short_code)
	}
	return long_url, nil
}
