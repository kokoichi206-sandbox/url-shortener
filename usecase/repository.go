package usecase

import "context"

type database interface {
	Health(ctx context.Context) error

	SearchURLFromShortURL(ctx context.Context, shortURL string) (string, error)
}
