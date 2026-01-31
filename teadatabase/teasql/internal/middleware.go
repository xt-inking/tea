package internal

import (
	"context"
	"database/sql"
)

func ExecerFunc(db *sql.DB, ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.ExecContext(ctx, query, args...)
}

func QueryerFunc(db *sql.DB, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.QueryContext(ctx, query, args...)
}

func BeginnerFunc(db *sql.DB, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.BeginTx(ctx, opts)
}

func TxExecerFunc(tx *sql.Tx, ctx context.Context, query string, args ...any) (sql.Result, error) {
	return tx.ExecContext(ctx, query, args...)
}

func TxQueryerFunc(tx *sql.Tx, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return tx.QueryContext(ctx, query, args...)
}

func TxCommitFunc(tx *sql.Tx, ctx context.Context) error {
	return tx.Commit()
}

func TxRollbackFunc(tx *sql.Tx, ctx context.Context) error {
	return tx.Rollback()
}
