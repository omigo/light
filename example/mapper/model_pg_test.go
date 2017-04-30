package mapper

import (
	"database/sql"

	// import sql driver
	_ "github.com/lib/pq"

	"github.com/arstd/log"
)

const url = "postgres://postgres:@127.0.0.1:5432/test?sslmode=disable"

func initPG() {
	var err error
	db, err = sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}

	connect()
	createTable()
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

func createTable() {
	_, err := db.Exec("drop table if exists models")
	if err != nil {
		log.Error(err)
	}
	_, err = db.Exec(`
		create table models (
			id serial primary key,
			name text not null,
			flag bool not null default false,
			score decimal(3,1) not null default 0.0,

			map jsonb not null default '{}',
			time timestamptz not null default now(),
			xarray text[] not null,
			slice text[] not null,

			status smallint not null default 0,
			state text not null default '',

			pointer jsonb not null default '{}',
			struct_slice jsonb not null default '[]',
			uint32 timestamptz not null default now()
		)
	`)
	if err != nil {
		log.Error(err)
	}
}
