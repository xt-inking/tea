package teasql

import (
	"context"
	"database/sql"
	"strconv"
	"sync/atomic"

	"github.com/tea-frame-go/tea/teaerrors"
)

type tx struct {
	Raw        *sql.Tx
	db         *DB
	identifier atomic.Int64
}

func newTx(raw *sql.Tx, db *DB) *tx {
	tx := &tx{
		Raw:        raw,
		db:         db,
		identifier: atomic.Int64{},
	}
	return tx
}

func (tx *tx) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.db.txExecer.Exec(tx.Raw, ctx, query, args...)
}

func (tx *tx) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.db.txQueryer.Query(tx.Raw, ctx, query, args...)
}

func (tx *tx) commit(ctx context.Context) error {
	return tx.db.txCommit.Commit(tx.Raw, ctx)
}

func (tx *tx) rollback(ctx context.Context) error {
	return tx.db.txRollback.Rollback(tx.Raw, ctx)
}

func (tx *tx) tx(ctx context.Context, f func(ctx context.Context) error) error {
	identifierCount := tx.identifier.Add(1)
	identifier := tx.db.sqlBuilder.QuoteIdentifier(strconv.FormatInt(identifierCount, 10))
	err := tx.beginSavepoint(ctx, identifier)
	if err != nil {
		return err
	}
	err = f(ctx)
	if err != nil {
		errRollbackSavepoint := tx.rollbackSavepoint(ctx, identifier)
		if errRollbackSavepoint != nil {
			return errRollbackSavepoint
		}
		return err
	}
	err = tx.commitSavepoint(ctx, identifier)
	if err != nil {
		tx.db.logger.Error(ctx, teaerrors.New(err, 0).ErrorStack())
	}
	return nil
}

func (tx *tx) beginSavepoint(ctx context.Context, identifier string) error {
	_, err := tx.exec(ctx, "SAVEPOINT "+identifier+";")
	return err
}

func (tx *tx) commitSavepoint(ctx context.Context, identifier string) error {
	_, err := tx.exec(ctx, "RELEASE SAVEPOINT "+identifier+";")
	return err
}

func (tx *tx) rollbackSavepoint(ctx context.Context, identifier string) error {
	_, err := tx.exec(ctx, "ROLLBACK TO "+identifier+";")
	return err
}

func txIntoContext(ctx context.Context, tx *tx) context.Context {
	return context.WithValue(ctx, keyContextTx, tx)
}

func txFromContext(ctx context.Context) (*tx, bool) {
	tx, ok := ctx.Value(keyContextTx).(*tx)
	return tx, ok
}

var keyContextTx = contextKeyTx{}

type contextKeyTx struct{}
