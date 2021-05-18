package database

import (
	"context"
	"database/sql"
	"fmt"
)

// This type helps us execute queries against a database.  It will internally manage
// connections, transations, etc. depending on our current needs.
type dbExecutor struct {
	db   *sql.DB
	tx   *sql.Tx
	ctx  context.Context
	txOK bool
}

// Create an instance of dbExecutor.  If transaction options are provided a transaction will be created, otherwise
// queries will be run
func CreateExecutor(connectFunc func() *sql.DB, txOptions *sql.TxOptions) *dbExecutor {
	return CreateExecutorWithContext(connectFunc, context.Background(), txOptions)
}

func CreateExecutorWithContext(connectFunc func() *sql.DB, ctx context.Context, txOptions *sql.TxOptions) *dbExecutor {
	res := &dbExecutor{
		db:   connectFunc(),
		ctx:  ctx,
		txOK: false,
	}
	if txOptions != nil {
		newTx, txErr := res.db.BeginTx(ctx, txOptions)
		if txErr != nil {
			panic(txErr)
		}
		res.tx = newTx
	}
	return res
}

func (dbe *dbExecutor) QueryRow(query string, args ...interface{}) *sql.Row {
	var res *sql.Row

	if dbe.tx != nil {
		res = dbe.tx.QueryRowContext(dbe.ctx, query, args...)
	} else {
		res = dbe.db.QueryRowContext(dbe.ctx, query, args...)
	}

	return res
}

func (dbe *dbExecutor) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var res *sql.Rows
	var err error

	if dbe.tx != nil {
		res, err = dbe.tx.QueryContext(dbe.ctx, query, args...)
	} else {
		res, err = dbe.db.QueryContext(dbe.ctx, query, args...)
	}

	return res, err
}

func (dbe *dbExecutor) Exec(query string, args ...interface{}) (sql.Result, error) {
	var res sql.Result
	var err error

	if dbe.tx != nil {
		res, err = dbe.tx.ExecContext(dbe.ctx, query, args...)
	} else {
		res, err = dbe.db.ExecContext(dbe.ctx, query, args...)
	}

	return res, err
}

func (dbe *dbExecutor) Complete() {
	if dbe.tx != nil {
		if dbe.txOK {
			dbe.tx.Commit()
		} else {
			dbe.tx.Rollback()
		}
	}

	// Always close our connection!
	dbe.db.Close()
}

func (dbe *dbExecutor) SetTransationFlag(status bool) {
	dbe.txOK = status
}
