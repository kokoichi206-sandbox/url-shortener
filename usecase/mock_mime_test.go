package usecase_test

import (
	"context"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
)

// FIXME: gomock で引数のメソッドを実行する方法がわからないため自作。
type myMockTxManager struct {
	ReadWriteTransactionFunc func(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error
}

func (m *myMockTxManager) ReadWriteTransaction(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) error {
	return m.ReadWriteTransactionFunc(ctx, f)
}
