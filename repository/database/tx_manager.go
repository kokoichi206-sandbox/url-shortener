package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
)

type txManager struct {
	db *sql.DB
}

func newTxManager(db *sql.DB) transaction.TxManager {
	return &txManager{
		db: db,
	}
}

func (t *txManager) ReadWriteTransaction(ctx context.Context, f func(ctx context.Context, tx transaction.RWTx) error) (err error) {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if err != nil {
			if e := tx.Rollback(); e != nil {
				err = fmt.Errorf("failed to rollback tx: %w", e)
			}

			return
		}

		if e := tx.Commit(); e != nil {
			err = fmt.Errorf("failed to commit tx: %w", e)
		}
	}()

	if err = f(ctx, &rwTx{tx}); err != nil {
		return fmt.Errorf("failed to execute f: %w", err)
	}

	return nil
}
