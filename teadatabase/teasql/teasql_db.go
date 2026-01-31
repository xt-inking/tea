package teasql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/tea-frame-go/tea/teaconfig"
	"github.com/tea-frame-go/tea/teadatabase/teasql/internal"
	"github.com/tea-frame-go/tea/teadatabase/teasql/teamiddleware"
	"github.com/tea-frame-go/tea/teaerrors"
	"github.com/tea-frame-go/tea/tealog"
)

type DB struct {
	Raw    *sql.DB
	logger *tealog.Logger

	sqlBuilder sqlBuilder

	execer   teamiddleware.Execer
	queryer  teamiddleware.Queryer
	beginner teamiddleware.Beginner

	txExecer  teamiddleware.TxExecer
	txQueryer teamiddleware.TxQueryer

	txCommit   teamiddleware.TxCommit
	txRollback teamiddleware.TxRollback
}

func NewDB(config *teaconfig.SqlConfig, options ...dbOption) *DB {
	raw, err := sql.Open(config.DriverName, config.DataSourceName)
	if err != nil {
		panic(err)
	}
	err = raw.Ping()
	if err != nil {
		panic(err)
	}
	raw.SetMaxOpenConns(config.MaxOpenConns)
	raw.SetMaxIdleConns(config.MaxIdleConns)
	raw.SetConnMaxLifetime(config.ConnMaxLifetime)
	raw.SetConnMaxIdleTime(config.ConnMaxIdleTime)
	sqlBuildersMu.RLock()
	sqlBuilderi, ok := sqlBuilders[config.DriverName]
	sqlBuildersMu.RUnlock()
	if !ok {
		panic(fmt.Errorf("teasql: unknown driver sqlBuilder %q (forgotten import?)", config.DriverName))
	}
	db := &DB{
		Raw: raw,
		logger: tealog.New(
			tealog.NewRecordHandlerText(),
			tealog.NewWriterCloserFile(tealog.FileDir, "sql"),
		),

		sqlBuilder: sqlBuilderi,

		execer:     teamiddleware.ExecerFunc(internal.ExecerFunc),
		queryer:    teamiddleware.QueryerFunc(internal.QueryerFunc),
		beginner:   teamiddleware.BeginnerFunc(internal.BeginnerFunc),
		txExecer:   teamiddleware.TxExecerFunc(internal.TxExecerFunc),
		txQueryer:  teamiddleware.TxQueryerFunc(internal.TxQueryerFunc),
		txCommit:   teamiddleware.TxCommitFunc(internal.TxCommitFunc),
		txRollback: teamiddleware.TxRollbackFunc(internal.TxRollbackFunc),
	}
	for _, o := range options {
		o(db)
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
