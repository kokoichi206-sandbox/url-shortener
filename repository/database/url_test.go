package database_test

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
	"github.com/kokoichi206-sandbox/url-shortener/repository/database"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

func Test_Database_SearchURLFromShortURL(t *testing.T) {
	t.Parallel()

	type args struct {
		shortURL string
	}

	testCases := map[string]struct {
		args     args
		makeMock func(m sqlmock.Sqlmock)
		want     string
		wantErr  string
	}{
		"success": {
			args: args{
				shortURL: "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.
					ExpectQuery(regexp.QuoteMeta(database.SearchURLFromShortURLStmt)).
					WithArgs("R0D").
					WillReturnRows(
						sqlmock.NewRows([]string{"url"}).
							AddRow("https://example.com"),
					)
			},
			want: "https://example.com",
		},
		"failure: no row found": {
			args: args{
				shortURL: "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.
					ExpectQuery(regexp.QuoteMeta(database.SearchURLFromShortURLStmt)).
					WithArgs("R0D").
					WillReturnError(sql.ErrNoRows)
			},
			wantErr: apperr.ErrShortURLNotFound.Error(),
		},
		"failure: scan error": {
			args: args{
				shortURL: "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.
					ExpectQuery(regexp.QuoteMeta(database.SearchURLFromShortURLStmt)).
					WithArgs("R0D").
					WillReturnError(errors.New("scan error"))
			},
			wantErr: "failed to scan: scan error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.makeMock(mock)

			logger := logger.NewBasicLogger(nil, "test", "database")

			database := database.New(db, logger)

			// Act
			got, err := database.SearchURLFromShortURL(context.Background(), tc.args.shortURL)

			// Assert
			assert.Equal(t, tc.want, got, "result does not match")
			if tc.wantErr == "" {
				require.NoError(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}

func Test_Database_SelectShortURL(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL string
	}

	testCases := map[string]struct {
		args            args
		makeMock        func(m sqlmock.Sqlmock)
		makeExtractRWTx func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error)
		want            string
		wantErr         string
	}{
		"success": {
			args: args{
				originalURL: "https://example.com",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectQuery(regexp.QuoteMeta(database.SelectShortURLStmt)).
					WithArgs("https://example.com").
					WillReturnRows(
						sqlmock.NewRows([]string{"short"}).
							AddRow("R0D"),
					)
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return &database.RwTx{sqlTx}, nil
				}
			},
			want: "R0D",
		},
		"failure: extract rwtx": {
			args: args{
				originalURL: "https://example.com",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectQuery(regexp.QuoteMeta(database.SelectShortURLStmt)).
					WithArgs("https://example.com").
					WillReturnRows(
						sqlmock.NewRows([]string{"short"}).
							AddRow("R0D"),
					)
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return nil, errors.New("extract rwtx error")
				}
			},
			wantErr: "failed to extract tx: extract rwtx error",
		},
		"failure: no row found": {
			args: args{
				originalURL: "https://wtf.example.com",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectQuery(regexp.QuoteMeta(database.SelectShortURLStmt)).
					WithArgs("https://wtf.example.com").
					WillReturnError(sql.ErrNoRows)
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return &database.RwTx{sqlTx}, nil
				}
			},
			wantErr: apperr.ErrShortURLNotFound.Error(),
		},
		"failure: scan error": {
			args: args{
				originalURL: "https://example.com",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectQuery(regexp.QuoteMeta(database.SelectShortURLStmt)).
					WithArgs("https://example.com").
					WillReturnError(errors.New("scan error"))
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return &database.RwTx{sqlTx}, nil
				}
			},
			wantErr: "failed to scan: scan error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.makeMock(mock)

			// TODO: よりうまくテストができないか考える。
			// より外側の txManager で tx が作成される想定だが、テストではここで作成する。
			tx, err := db.BeginTx(context.Background(), nil)
			require.NoError(t, err, "error of BeginTx should be nil")
			rwt := &database.RwTx{tx}

			urlRepo := database.NewURLRepo(tc.makeExtractRWTx(tx))

			// Act
			got, err := urlRepo.SelectShortURL(context.Background(), rwt, tc.args.originalURL)

			// Assert
			assert.Equal(t, tc.want, got, "result does not match")
			if tc.wantErr == "" {
				require.NoError(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}

func Test_Database_InsertURL(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL string
		shortURL    string
	}

	testCases := map[string]struct {
		args            args
		makeMock        func(m sqlmock.Sqlmock)
		makeExtractRWTx func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error)
		want            string
		wantErr         string
	}{
		"success": {
			args: args{
				originalURL: "https://example.com",
				shortURL:    "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectExec(regexp.QuoteMeta(database.InsertURLStmt)).
					WithArgs("https://example.com", "R0D").
					WillReturnResult(driver.RowsAffected(1))
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return &database.RwTx{sqlTx}, nil
				}
			},
		},
		"failure: extract rwtx": {
			args: args{
				originalURL: "https://example.com",
				shortURL:    "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectExec(regexp.QuoteMeta(database.InsertURLStmt)).
					WithArgs("https://example.com", "R0D").
					WillReturnResult(driver.RowsAffected(1))
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return nil, errors.New("extract rwtx error")
				}
			},
			wantErr: "failed to extract tx: extract rwtx error",
		},
		"failure: exec error": {
			args: args{
				originalURL: "https://example.com",
				shortURL:    "R0D",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.
					ExpectExec(regexp.QuoteMeta(database.InsertURLStmt)).
					WithArgs("https://example.com", "R0D").
					WillReturnError(errors.New("exec error"))
			},
			makeExtractRWTx: func(sqlTx *sql.Tx) func(transaction.RWTx) (*database.RwTx, error) {
				return func(r transaction.RWTx) (*database.RwTx, error) {
					return &database.RwTx{sqlTx}, nil
				}
			},
			wantErr: "failed to insert: exec error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
			}
			defer db.Close()

			tc.makeMock(mock)

			// TODO: よりうまくテストができないか考える。
			// より外側の txManager で tx が作成される想定だが、テストではここで作成する。
			tx, err := db.BeginTx(context.Background(), nil)
			require.NoError(t, err, "error of BeginTx should be nil")
			rwt := &database.RwTx{tx}

			urlRepo := database.NewURLRepo(tc.makeExtractRWTx(tx))

			// Act
			err = urlRepo.InsertURL(context.Background(), rwt, tc.args.originalURL, tc.args.shortURL)

			// Assert
			if tc.wantErr == "" {
				require.NoError(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}
