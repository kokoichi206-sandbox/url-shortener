package database

import (
	"context"
	"database/sql"
	"fmt"

	tracer "github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
)

const searchURLFromShortURLQuery = `
SELECT
	url
FROM shorturl
WHERE short = $1;
`

func (d *database) SearchURLFromShortURL(ctx context.Context, shortURL string) (string, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "d.SearchURLFromShortURL")
	defer span.Finish()

	row := d.db.QueryRowContext(ctx, searchURLFromShortURLQuery, shortURL)

	var url string
	if err := row.Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			return "", apperr.ShortURLNotFound
		}

		return "", fmt.Errorf("failed to scan: %w", err)
	}

	return url, nil
}
