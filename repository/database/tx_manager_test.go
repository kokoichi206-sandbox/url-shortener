package database_test

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/repository/database"
	"github.com/stretchr/testify/assert"
)

func Test_Database_NewTxManager(t *testing.T) {
	t.Parallel()

	type args struct {
		f func(ctx context.Context, tx transaction.RWTx) error
	}

	testCases := map[string]struct {
		args     args
		makeMock func(m sqlmock.Sqlmock)
		want     string
		wantErr  string
	}{
		"success": {
			args: args{
				f: func(ctx context.Context, tx transaction.RWTx) error {
					// transaction success !
					return nil
				},
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				// f が正常終了した時は、変更内容が commit されること。
				m.ExpectCommit()
			},
		},
		"success: rollback due to function error": {
			args: args{
				f: func(ctx context.Context, tx transaction.RWTx) error {
					// transaction failure ...
					return errors.New("f error")
				},
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				// f が異常終了した時は、変更内容が rollback されること。
				m.ExpectRollback()
			},
			wantErr: "failed to execute f: f error",
		},
		"failure: begin tx": {
			args: args{
				f: func(ctx context.Context, tx transaction.RWTx) error {
					// transaction success !
					return nil
				},
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin().WillReturnError(errors.New("begin error"))
			},
			wantErr: "failed to begin tx: begin error",
		},
		"failure: rollback": {
			args: args{
				f: func(ctx context.Context, tx transaction.RWTx) error {
					// transaction failure ...
					return errors.New("f error")
				},
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectRollback().WillReturnError(errors.New("rollback error"))
			},
			wantErr: "failed to rollback tx: rollback error",
		},
		"failure: commit": {
			args: args{
				f: func(ctx context.Context, tx transaction.RWTx) error {
					// transaction success !
					return nil
				},
			},
			makeMock: func(m sqlmock.Sqlmock) {
				m.ExpectBegin()
				m.ExpectCommit().WillReturnError(errors.New("commit error"))
			},
			wantErr: "failed to commit tx: commit error",
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

			txManager := database.NewTxManager(db)

			// Act
			err = txManager.ReadWriteTransaction(context.Background(), tc.args.f)

			// Assert
			if tc.wantErr == "" {
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}
