package database_test

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/repository/database"
)

func Test_Database_ExtractRWTx(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL string
	}

	testCases := map[string]struct {
		args     args
		makeMock func(m sqlmock.Sqlmock)
		wantErr  string
	}{
		"success": {
			args: args{
				originalURL: "https://example.com",
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
			},
		},
		"failure": {
			args: args{
				originalURL: "https://example.com",
			},
			wantErr: "failed to extract rwTx of sql",
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

			var rwt transaction.RWTx
			// var err error

			if tc.makeMock != nil {
				tc.makeMock(mock)
				tx, err := db.BeginTx(context.Background(), nil)
				assert.Nil(t, err, "error of BeginTx should be nil")
				rwt = &database.RwTx{tx}
			}

			// Act
			_, err = database.ExtractRWTx(rwt)

			// Assert
			if tc.wantErr == "" {
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}
