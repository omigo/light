package mapper

import (
	"database/sql"

	// import sql driver
	_ "github.com/lib/pq"

	"github.com/arstd/log"
)

const url = "postgres://postgres:@127.0.0.1:5432/test?sslmode=disable"

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	// connect()
}

func connect() (err error) {
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return err
	}

	log.Infof("successfully connect to %s", url)
	return nil
}

func BeginTx() (*sql.Tx, error) {
	return db.Begin()
}

func CommitTx(tx *sql.Tx) error {
	return tx.Commit()
}

func RollbackTx(tx *sql.Tx) error {
	return tx.Rollback()
}
