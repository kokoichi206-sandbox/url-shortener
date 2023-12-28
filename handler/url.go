package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tracer "github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
	"github.com/kokoichi206-sandbox/url-shortener/model/request"
)

func (h *handler) GetOriginalURL(c *gin.Context) error {
	ctx := c.Request.Context()

	span, ctx := tracer.StartSpanFromContext(ctx, "h.GetOriginalURL")
	defer span.Finish()

	shortURL := c.Param("shortURL")

	url, err := h.usecase.SearchOriginalURL(ctx, shortURL)
	if err != nil {
		return fmt.Errorf("failed to exec usecase.SearchOriginalURL: %w", err)
	}

	c.Redirect(http.StatusMovedPermanently, url)

	return nil
}

func (h *handler) GenerateURL(c *gin.Context) error {
	ctx := c.Request.Context()

	span, ctx := tracer.StartSpanFromContext(ctx, "h.GenerateURL")
	defer span.Finish()

	var body request.CreateURL
	if err := json.NewDecoder(c.Request.Body).Decode(&body); err != nil {
		// body is empty or invalid json format.
		return apperr.ErrRequestBodyInvalid
	}

	url, err := h.usecase.GenerateURL(ctx, body.OriginalURL)
	if err != nil {
		return fmt.Errorf("failed to exec usecase.SearchOriginalURL: %w", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"short_url": url,
	})

	return nil
}
