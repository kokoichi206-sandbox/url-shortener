package usecase_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
	"github.com/kokoichi206-sandbox/url-shortener/model/apperr"
	"github.com/kokoichi206-sandbox/url-shortener/usecase"
	"github.com/kokoichi206-sandbox/url-shortener/util/logger"
)

func Test_Usecase_SearchOriginalURL(t *testing.T) {
	t.Parallel()

	type args struct {
		shortURL string
	}

	testCases := map[string]struct {
		args             args
		makeMockDatabase func(m *MockDatabase)
		want             string
		wantErr          string
	}{
		"success": {
			args: args{
				shortURL: "R0D",
			},
			makeMockDatabase: func(m *MockDatabase) {
				m.
					EXPECT().
					SearchURLFromShortURL(gomock.Any(), "R0D").
					Return("https://example.com", nil)
			},
			want: "https://example.com",
		},
		"failure: no url in repository": {
			args: args{
				shortURL: "NUL",
			},
			makeMockDatabase: func(m *MockDatabase) {
				m.
					EXPECT().
					SearchURLFromShortURL(gomock.Any(), "NUL").
					Return("", apperr.ErrShortURLNotFound)
			},
			wantErr: "failed to search url from database: short url not found",
		},
		"failure: db error": {
			args: args{
				shortURL: "R0D",
			},
			makeMockDatabase: func(m *MockDatabase) {
				m.
					EXPECT().
					SearchURLFromShortURL(gomock.Any(), "R0D").
					Return("", errors.New("db error"))
			},
			wantErr: "failed to search url from database: db error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			m := NewMockDatabase(ctrl)
			tc.makeMockDatabase(m)

			b := bytes.NewBuffer([]byte{})
			logger.NewBasicLogger(b, "test", "searchOriginalURL")

			u := usecase.New(m, nil, nil, nil)

			// Act
			got, err := u.SearchOriginalURL(context.Background(), tc.args.shortURL)

			// Assert
			assert.Equal(t, tc.want, got, "result does not match")
			if tc.wantErr == "" {
				assert.NoError(t, err, "error should be nil")
			} else {
				assert.Regexp(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}

func Test_Usecase_GenerateURL(t *testing.T) {
	t.Parallel()

	type args struct {
		originalURL string
	}

	testCases := map[string]struct {
		args             args
		makeMockDatabase func(m *MockDatabase)
		makeURLsRepo     func(m *MockURLRepository)
		myMockTxManager  *myMockTxManager // FIXME: gomock で引数のメソッドを実行する方法がわからないため自作。
		want             string
		wantErr          string
	}{
		"success": {
			args: args{
				originalURL: "https://example.com",
			},
			makeURLsRepo: func(m *MockURLRepository) {
				m.
					EXPECT().
					SelectShortURL(gomock.Any(), gomock.Any(), "https://example.com").
					Return("", apperr.ErrShortURLNotFound)
				m.
					EXPECT().
					InsertURL(gomock.Any(), gomock.Any(), "https://example.com", gomock.Any()).
					Return(nil)
			},
			myMockTxManager: &myMockTxManager{
				// ReadWriteTransaction の中を実行させるために、ここで実行する関数を定義する。
				ReadWriteTransactionFunc: func(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error {
					return f(ctx, nil)
				},
			},
			want: "[a-zA-Z0-9]{3}",
		},
		"success: existing url": {
			args: args{
				originalURL: "https://example.com",
			},
			makeURLsRepo: func(m *MockURLRepository) {
				m.
					EXPECT().
					SelectShortURL(gomock.Any(), gomock.Any(), "https://example.com").
					Return("R0D", apperr.ErrShortURLNotFound)
			},
			myMockTxManager: &myMockTxManager{
				ReadWriteTransactionFunc: func(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error {
					return f(ctx, nil)
				},
			},
			want: "R0D",
		},
		"failure: select url": {
			args: args{
				originalURL: "https://example.com",
			},
			makeURLsRepo: func(m *MockURLRepository) {
				m.
					EXPECT().
					SelectShortURL(gomock.Any(), gomock.Any(), "https://example.com").
					Return("", errors.New("db error"))
			},
			myMockTxManager: &myMockTxManager{
				ReadWriteTransactionFunc: func(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error {
					return f(ctx, nil)
				},
			},
			wantErr: "failed to exec txManager.ReadWriteTransaction: failed to select short url from database: db error",
		},
		"failure: insert url": {
			args: args{
				originalURL: "https://example.com",
			},
			makeURLsRepo: func(m *MockURLRepository) {
				m.
					EXPECT().
					SelectShortURL(gomock.Any(), gomock.Any(), "https://example.com").
					Return("", apperr.ErrShortURLNotFound)
				m.
					EXPECT().
					InsertURL(gomock.Any(), gomock.Any(), "https://example.com", gomock.Any()).
					Return(errors.New("db error"))
			},
			myMockTxManager: &myMockTxManager{
				ReadWriteTransactionFunc: func(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error {
					return f(ctx, nil)
				},
			},
			wantErr: "failed to exec txManager.ReadWriteTransaction: failed to insert short url to database: db error",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ur := NewMockURLRepository(ctrl)
			tc.makeURLsRepo(ur)

			b := bytes.NewBuffer([]byte{})
			logger.NewBasicLogger(b, "test", "generateURL")

			u := usecase.New(nil, tc.myMockTxManager, ur, nil)

			// Act
			got, err := u.GenerateURL(context.Background(), tc.args.originalURL)

			// Assert
			assert.Regexp(t, tc.want, got, "result does not match")
			if tc.wantErr == "" {
				assert.NoError(t, err, "error should be nil")
			} else {
				assert.Regexp(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}

func Test_Usecase_GetRoomUsers(t *testing.T) {
	t.Parallel()

	type args struct {
		num int
	}

	testCases := map[string]struct {
		args    args
		wantReg string
		wantErr string
	}{
		"success": {
			args: args{
				num: 3,
			},
			wantReg: "^[a-zA-Z0-9]{3}$",
		},
		"success: length 5": {
			args: args{
				num: 5,
			},
			wantReg: "^[a-zA-Z0-9]{5}$",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange

			// Act
			got, err := usecase.GenerateRandomString(tc.args.num)

			// Assert
			assert.Regexp(t, tc.wantReg, got, "result does not match")
			if tc.wantErr == "" {
				assert.NoError(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}
