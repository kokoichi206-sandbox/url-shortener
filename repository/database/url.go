package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	tracer "github.com/opentracing/opentracing-go"

	"github.com/kokoichi206-sandbox/url-shortener/domain/repository"
	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
)

const searchURLFromShortURLStmt = `
SELECT
	url
FROM shorturl
WHERE short = $1;
`

func (d *database) SearchURLFromShortURL(ctx context.Context, shortURL string) (string, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "d.SearchURLFromShortURL")
	defer span.Finish()

	row := d.db.QueryRowContext(ctx, searchURLFromShortURLStmt, shortURL)

	var url string
	if err := row.Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperr.ErrShortURLNotFound
		}

		return "", fmt.Errorf("failed to scan: %w", err)
	}

	return url, nil
}

type urlRepo struct {
	extractRWTx func(transaction.RWTx) (*RwTx, error)
}

func NewURLRepo(
	extractRWTx func(transaction.RWTx) (*RwTx, error),
) repository.URLRepository {
	return &urlRepo{
		extractRWTx: extractRWTx,
	}
}

const selectShortURLStmt = `
SELECT
	short
FROM shorturl
WHERE url = $1;
`

func (u *urlRepo) SelectShortURL(ctx context.Context, ttx transaction.RWTx, originalURL string) (string, error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "u.SelectShortURL")
	defer span.Finish()

	tx, err := u.extractRWTx(ttx)
	if err != nil {
		return "", fmt.Errorf("failed to extract tx: %w", err)
	}

	row := tx.QueryRowContext(ctx, selectShortURLStmt, originalURL)

	var shortURL string
	if err := row.Scan(&shortURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", apperr.ErrShortURLNotFound
		}

		return "", fmt.Errorf("failed to scan: %w", err)
	}

	return shortURL, nil
}

const insertURLStmt = `
INSERT INTO shorturl (
	url,
	short
) VALUES (
	$1,
	$2
);
`

func (u *urlRepo) InsertURL(ctx context.Context, ttx transaction.RWTx, originalURL string, shortURL string) error {
	span, ctx := tracer.StartSpanFromContext(ctx, "t.InsertURL")
	defer span.Finish()

	tx, err := u.extractRWTx(ttx)
	if err != nil {
		return fmt.Errorf("failed to extract tx: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertURLStmt, originalURL, shortURL); err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}
