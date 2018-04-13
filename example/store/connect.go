package store

import (
	"database/sql"
	"fmt"

	"github.com/arstd/light/example/conf"
	"github.com/arstd/log"
)

var db *sql.DB

func init() {
	open()
	log.Fataln(Connect())
}

func open() {
	log.Json(conf.Conf)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		conf.Conf.DB.Username,
		conf.Conf.DB.Password,
		conf.Conf.DB.Host,
		conf.Conf.DB.Port,
		conf.Conf.DB.DBName,
		conf.Conf.DB.Params,
	)
	var err error
	db, err = sql.Open(conf.Conf.DB.Dialect, dsn)
	log.Fataln(err)

	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(1)
	db.SetConnMaxLifetime(0)
}

func Connect() error {
	return db.Ping()
}

func Close() {
	if db != nil {
		log.Errorn(db.Close())
	}
}

func Begin() (*sql.Tx, error) {
	return db.Begin()
}
