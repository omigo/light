package light

import (
	"database/sql"
)

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

func GetExecer(db *sql.DB, txs []*sql.Tx) Execer {
	if len(txs) > 0 {
		return txs[0]
	}
	return db
}
