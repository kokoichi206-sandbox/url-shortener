package usecase

import (
	"context"

	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

type Usecase interface {
	Health(ctx context.Context) error
}

type usecase struct {
	database database

	logger logger.Logger
}

func New(database database, logger logger.Logger) Usecase {
	usecase := &usecase{
		database: database,
		logger:   logger,
	}

	return usecase
}

func (u *usecase) Health(ctx context.Context) error {
	// db の接続確認。
	return u.database.Health(ctx)
}
