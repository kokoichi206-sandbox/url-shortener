package repository

import (
	"context"
)

type Database interface {
	Health(ctx context.Context) error

	SearchURLFromShortURL(ctx context.Context, shortURL string) (string, error)
}

type roTx interface {
	ROTxImpl()
}

type rwTx interface {
	ROTxRWTxImpl()
}

type txManager interface {
	ReadWriteTransaction(ctx context.Context, f func(ctx context.Context, tx rwTx) error) error
}
