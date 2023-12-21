package usecase

import (
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

type Usecase interface{}

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
