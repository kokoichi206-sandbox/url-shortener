package database

import (
	"database/sql"
	"errors"

	"github.com/kokoichi206-sandbox/url-shortener/domain/transaction"
)

type RwTx struct {
	*sql.Tx
}

func (t *RwTx) ROTxImpl() {}
func (t *RwTx) RWTxImpl() {}

func ExtractRWTx(_tx transaction.RWTx) (*RwTx, error) {
	tx, ok := _tx.(*RwTx)
	if !ok {
		return nil, errors.New("failed to extract rwTx of sql")
	}

	return tx, nil
}
