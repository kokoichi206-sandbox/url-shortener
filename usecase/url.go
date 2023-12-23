package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	tracer "github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
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

		shortURL = uuid.New().String()

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
