package usecase

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/lib/pq"
	tracer "github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
)

const (
	shortenedURLLength = 3
)

func (u *usecase) SearchOriginalURL(ctx context.Context, shortURL string) (string, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "d.SearchURLFromShortURL")
	defer span.Finish()

	url, err := u.database.SearchURLFromShortURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("failed to search url from database: %w", err)
	}

	return url, nil
}

func (u *usecase) GenerateURL(ctx context.Context, originalURL string) (string, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "d.GenerateURL")
	defer span.Finish()

	maxRetries := 3
	retries := 0

	for {
		if retries >= maxRetries {
			return "", fmt.Errorf("failed to insert short url due to duplicate key error: (retry count: %v)", retries)
		}

		shortURL, err := u.fetchOrGenerateShortURL(ctx, originalURL)
		if err != nil {
			var pqErr *pq.Error

			// Error code 23505 means 'unique_violation' error in PostgreSQL.
			// This error occurs when attempting to insert a short URL that already exists in the database.
			// In this case, regenerate a new short URL and retry the insertion process.
			if errors.As(err, &pqErr) && pqErr.Code == "23505" {
				retries++

				continue
			}

			return "", fmt.Errorf("failed to insert short url to database: %w", err)
		}

		return shortURL, nil
	}
}

func (u *usecase) fetchOrGenerateShortURL(ctx context.Context, originalURL string) (string, error) {
	var shortURL string

	if err := u.txManager.ReadWriteTransaction(ctx, func(ctx context.Context, tx transaction.RWTx) error {
		var err error

		shortURL, err = u.urlRepo.SelectShortURL(ctx, tx, originalURL)
		if err != nil && !errors.Is(err, apperr.ErrShortURLNotFound) {
			return fmt.Errorf("failed to select short url from database: %w", err)
		}

		if shortURL != "" {
			return nil
		}

		shortURL, err = generateRandomString(shortenedURLLength)
		if err != nil {
			return fmt.Errorf("failed to generate random string: %w", err)
		}

		err = u.urlRepo.InsertURL(ctx, tx, originalURL, shortURL)
		if err != nil {
			return fmt.Errorf("failed to insert short url to database: %w", err)
		}

		return nil
	}); err != nil {
		return "", fmt.Errorf("failed to exec txManager.ReadWriteTransaction: %w", err)
	}

	return shortURL, nil
}

// [a-zA-Z0-9] からランダムに n 文字の文字列を生成する。
func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, n)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", fmt.Errorf("failed to rand.Int: %w", err)
		}

		result[i] = letters[num.Int64()]
	}

	return string(result), nil
}
