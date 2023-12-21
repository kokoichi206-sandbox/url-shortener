package usecase

import (
	"context"
	"fmt"

	tracer "github.com/opentracing/opentracing-go"
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
