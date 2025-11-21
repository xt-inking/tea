package teasql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tea-frame-go/tea/teadatabase/teasql/internal"
	"github.com/tea-frame-go/tea/teadatabase/teasql/teamiddleware"
	"github.com/tea-frame-go/tea/teaerrors"
	"github.com/tea-frame-go/tea/tealog"
)

type DB struct {
	Raw    *sql.DB
	logger *tealog.Logger

	identifierQuoter identifierQuoter

	execer   teamiddleware.Execer
	queryer  teamiddleware.Queryer
	beginner teamiddleware.Beginner

	txExecer  teamiddleware.TxExecer
	txQueryer teamiddleware.TxQueryer

	txCommit   teamiddleware.TxCommit
	txRollback teamiddleware.TxRollback
}

func NewDB(driverName string, dataSourceName string, middleware ...teamiddleware.Middleware) *DB {
	raw, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		panic(err)
	}
	err = raw.Ping()
	if err != nil {
		panic(err)
	}
	identifierQuotersMu.RLock()
	identifierQuoteri, ok := identifierQuoters[driverName]
	identifierQuotersMu.RUnlock()
	if !ok {
		panic(fmt.Errorf("teasql: unknown driver identifierQuoter %q (forgotten import?)", driverName))
	}
	var (
		execer     teamiddleware.Execer     = teamiddleware.ExecerFunc(internal.ExecerFunc)
		queryer    teamiddleware.Queryer    = teamiddleware.QueryerFunc(internal.QueryerFunc)
		beginner   teamiddleware.Beginner   = teamiddleware.BeginnerFunc(internal.BeginnerFunc)
		txExecer   teamiddleware.TxExecer   = teamiddleware.TxExecerFunc(internal.TxExecerFunc)
		txQueryer  teamiddleware.TxQueryer  = teamiddleware.TxQueryerFunc(internal.TxQueryerFunc)
		txCommit   teamiddleware.TxCommit   = teamiddleware.TxCommitFunc(internal.TxCommitFunc)
		txRollback teamiddleware.TxRollback = teamiddleware.TxRollbackFunc(internal.TxRollbackFunc)
	)
	for i := len(middleware) - 1; i >= 0; i-- {
		execer = middleware[i].Execer(execer)
		queryer = middleware[i].Queryer(queryer)
		beginner = middleware[i].Beginner(beginner)
		txExecer = middleware[i].TxExecer(txExecer)
		txQueryer = middleware[i].TxQueryer(txQueryer)
		txCommit = middleware[i].TxCommit(txCommit)
		txRollback = middleware[i].TxRollback(txRollback)
	}
	db := &DB{
		Raw: raw,
		logger: tealog.New(
			tealog.NewRecordHandlerText(),
			tealog.NewWriterCloserFile(tealog.FileDir, "sql"),
		),

		identifierQuoter: identifierQuoteri,

		execer:     execer,
		queryer:    queryer,
		beginner:   beginner,
		txExecer:   txExecer,
		txQueryer:  txQueryer,
		txCommit:   txCommit,
		txRollback: txRollback,
	}
	return db
}

func (db *DB) Close() error {
	defer func() {
		db.logger.Close()
	}()
	return db.Raw.Close()
}

func (db *DB) Exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.exec(ctx, query, args...)
	}
	return db.exec(ctx, query, args...)
}

func (db *DB) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := txFromContext(ctx); ok {
		return tx.query(ctx, query, args...)
	}
	return db.query(ctx, query, args...)
}

func (db *DB) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.execer.Exec(db.Raw, ctx, query, args...)
}

func (db *DB) query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.queryer.Query(db.Raw, ctx, query, args...)
}

func (db *DB) begin(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.beginner.Begin(db.Raw, ctx, opts)
}

func (db *DB) Tx(ctx context.Context, opts *sql.TxOptions, f func(ctx context.Context) error) error {
	if tx, ok := txFromContext(ctx); ok {
		return tx.tx(ctx, f)
	}
	raw, err := db.begin(ctx, opts)
	if err != nil {
		return err
	}
	tx := newTx(raw, db)
	ctx = txIntoContext(ctx, tx)
	defer func() {
		errRollback := tx.rollback(ctx)
		if errRollback != nil {
			db.logger.Error(ctx, teaerrors.New(errRollback, 0).ErrorStack())
		}
	}()
	err = f(ctx)
	if err != nil {
		return err
	}
	return tx.commit(ctx)
}
