package teasql

import (
	"github.com/tea-frame-go/tea/teadatabase/teasql/teamiddleware"
	"github.com/tea-frame-go/tea/tealog"
)

var DBOptions = dbOptions{}

type dbOptions struct{}

func (dbOptions) Logger(logger *tealog.Logger) dbOption {
	return func(db *DB) {
		db.logger.Close()
		db.logger = logger
	}
}

func (dbOptions) Middleware(middleware ...teamiddleware.Middleware) dbOption {
	return func(db *DB) {
		for i := len(middleware) - 1; i >= 0; i-- {
			db.execer = middleware[i].Execer(db.execer)
			db.queryer = middleware[i].Queryer(db.queryer)
			db.beginner = middleware[i].Beginner(db.beginner)
			db.txExecer = middleware[i].TxExecer(db.txExecer)
			db.txQueryer = middleware[i].TxQueryer(db.txQueryer)
			db.txCommit = middleware[i].TxCommit(db.txCommit)
			db.txRollback = middleware[i].TxRollback(db.txRollback)
		}
	}
}

type dbOption func(db *DB)
