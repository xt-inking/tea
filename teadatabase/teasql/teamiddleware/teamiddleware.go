package teamiddleware

import (
	"context"
	"database/sql"
)

type Middleware interface {
	Execer(Execer) Execer
	Queryer(Queryer) Queryer
	Beginner(Beginner) Beginner

	TxExecer(TxExecer) TxExecer
	TxQueryer(TxQueryer) TxQueryer

	TxCommit(TxCommit) TxCommit
	TxRollback(TxRollback) TxRollback
}

type Execer interface {
	Exec(db *sql.DB, ctx context.Context, query string, args ...any) (sql.Result, error)
}

type ExecerFunc func(db *sql.DB, ctx context.Context, query string, args ...any) (sql.Result, error)

func (f ExecerFunc) Exec(db *sql.DB, ctx context.Context, query string, args ...any) (sql.Result, error) {
	return f(db, ctx, query, args...)
}

type Queryer interface {
	Query(db *sql.DB, ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type QueryerFunc func(db *sql.DB, ctx context.Context, query string, args ...any) (*sql.Rows, error)

func (f QueryerFunc) Query(db *sql.DB, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return f(db, ctx, query, args...)
}

type Beginner interface {
	Begin(db *sql.DB, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type BeginnerFunc func(db *sql.DB, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)

func (f BeginnerFunc) Begin(db *sql.DB, ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return f(db, ctx, opts)
}

type TxExecer interface {
	Exec(tx *sql.Tx, ctx context.Context, query string, args ...any) (sql.Result, error)
}

type TxExecerFunc func(tx *sql.Tx, ctx context.Context, query string, args ...any) (sql.Result, error)

func (f TxExecerFunc) Exec(tx *sql.Tx, ctx context.Context, query string, args ...any) (sql.Result, error) {
	return f(tx, ctx, query, args...)
}

type TxQueryer interface {
	Query(tx *sql.Tx, ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type TxQueryerFunc func(tx *sql.Tx, ctx context.Context, query string, args ...any) (*sql.Rows, error)

func (f TxQueryerFunc) Query(tx *sql.Tx, ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return f(tx, ctx, query, args...)
}

type TxCommit interface {
	Commit(tx *sql.Tx, ctx context.Context) error
}

type TxCommitFunc func(tx *sql.Tx, ctx context.Context) error

func (f TxCommitFunc) Commit(tx *sql.Tx, ctx context.Context) error {
	return f(tx, ctx)
}

type TxRollback interface {
	Rollback(tx *sql.Tx, ctx context.Context) error
}

type TxRollbackFunc func(tx *sql.Tx, ctx context.Context) error

func (f TxRollbackFunc) Rollback(tx *sql.Tx, ctx context.Context) error {
	return f(tx, ctx)
}
