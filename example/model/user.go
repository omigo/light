package model

import "time"

type Status = uint8

// create table users (
// 	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
// 	username VARCHAR(32) NOT NULL UNIQUE,
// 	Phone VARCHAR(32),
// 	address VARCHAR(256),
// 	status TINYINT UNSIGNED,
// 	birthday DATE,
// 	created TIMESTAMP default CURRENT_TIMESTAMP,
// 	updated TIMESTAMP default CURRENT_TIMESTAMP
// )
type User struct {
	Id       uint64     `db:"id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY"`
	Username string     `db:"username VARCHAR(32) NOT NULL UNIQUE"`
	Phone    string     `db:"Phone VARCHAR(32)"`
	Address  *string    `db:"address VARCHAR(256)"`
	Status   Status     `db:"status TINYINT UNSIGNED"`
	Birthday *time.Time `db:"birthday DATE"`
	Created  time.Time  `db:"created TIMESTAMP default CURRENT_TIMESTAMP"`
	Updated  time.Time  `db:"updated TIMESTAMP default CURRENT_TIMESTAMP"`
}
