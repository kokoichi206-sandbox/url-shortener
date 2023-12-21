package usecase

import (
	"github.com/kokoichi206-sandbox/url-shortener/repository"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

type Usecase interface{}

type usecase struct {
	database repository.Database

	logger logger.Logger
}

func New(database repository.Database, logger logger.Logger) Usecase {
	usecase := &usecase{
		database: database,
		logger:   logger,
	}

	return usecase
}
