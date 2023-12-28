package repository

import (
	"context"
)

type Database interface {
	Health(ctx context.Context) error

	SearchURLFromShortURL(ctx context.Context, shortURL string) (string, error)
}
