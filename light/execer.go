package light

import (
	"context"
	"database/sql"
)

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

func GetExec(tx *sql.Tx, db *sql.DB) Execer {
	if tx != nil {
		return tx
	}
	return db
}
