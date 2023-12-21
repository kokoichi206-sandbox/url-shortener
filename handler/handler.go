package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
	"github.com/kokoichi206-sandbox/url-shortener/usecase"
	"github.com/kokoichi206-sandbox/url-shortener/util"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

type handler struct {
	logger  logger.Logger
	usecase usecase.Usecase

	Engine *gin.Engine
}

//nolint:revive
func New(logger logger.Logger, usecase usecase.Usecase) *handler {
	r := gin.Default()

	h := &handler{
		logger:  logger,
		usecase: usecase,
		Engine:  r,
	}
	h.setupRoutes()

	return h
}

func (h *handler) setupRoutes() {
	base := h.Engine.Group("/api/v1")
	base.Use(h.requestIDMW())

	base.Handle(http.MethodGet, "/health", handlerWrapper(h.Health, h.logger))
}

func handlerWrapper(fun func(c *gin.Context) error, logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := fun(c); err != nil {
			handleError(c, logger, err)
		}
	}
}

// handleError is a helper function to handle error.
// This function writes status code and error message to response body.
func handleError(c *gin.Context, logger logger.Logger, err error) {
	var e apperr.AppError
	if ok := errors.As(err, &e); !ok {
		e = apperr.AppError{
			StatusCode: http.StatusInternalServerError,
			Message:    "internal server error",
			Log:        err.Error(),
		}
	}

	if e.Log != "" {
		logger.Error(util.DetachedCtx{
			Parent: c.Request.Context(),
		}, e.Log)
	}

	c.JSON(e.StatusCode, gin.H{
		"error": e.Message,
	})
}

func (h *handler) Health(c *gin.Context) error {
	if err := h.usecase.Health(c.Request.Context()); err != nil {
		return fmt.Errorf("failed to health check: %w", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"health": "ok",
	})

	return nil
}
