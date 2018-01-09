package model

import "time"

type Status = uint8

/*
create table users (
	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	username VARCHAR(32) NOT NULL UNIQUE,
	Phone VARCHAR(32),
	address VARCHAR(256),
	status TINYINT UNSIGNED,
	birthday DATE,
	created TIMESTAMP default CURRENT_TIMESTAMP,
	updated TIMESTAMP default CURRENT_TIMESTAMP
)
*/
type User struct {
	Id       uint64
	Username string
	Phone    string     `light:",nullable"`
	Address  *string    `light:",nullable"`
	Status   Status     `light:",nullable"`
	Birthday *time.Time `light:",nullable"`
	Created  time.Time
	Updated  time.Time
}
