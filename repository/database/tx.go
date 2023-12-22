package database

import (
	"database/sql"
	"errors"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
)

type rwTx struct {
	*sql.Tx
}

func (t *rwTx) ROTxImpl() {}
func (t *rwTx) RWTxImpl() {}

func ExtractRWTx(_tx transaction.ROTx) (*rwTx, error) {
	tx, ok := _tx.(*rwTx)
	if !ok {
		return nil, errors.New("failed to extract rwTx of sql")
	}

	return tx, nil
}
