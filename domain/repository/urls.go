package repository

import (
	"context"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
)

type URLRepository interface {
	SelectShortURL(ctx context.Context, tx transaction.RWTx, originalURL string) (string, error)
	InsertURL(ctx context.Context, tx transaction.RWTx, originalURL string, shortURL string) error
}
