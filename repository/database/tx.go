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

var _ transaction.ROTx = (*RwTx)(nil)

func ExtractRWTx(tx transaction.RWTx) (*RwTx, error) {
	rwTx, ok := tx.(*RwTx)
	if !ok {
		return nil, errors.New("failed to extract rwTx of sql")
	}

	return rwTx, nil
}
