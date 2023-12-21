package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	tracer "github.com/opentracing/opentracing-go"
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

	fmt.Printf("url: %v\n", url)
	c.Redirect(http.StatusMovedPermanently, url)

	return nil
}
