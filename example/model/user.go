package model

import "time"

type Status = uint8

type User struct {
	Id       uint64
	Username string
	Phone    string `light:",nullable"`
	Address  *string
	Status   Status `light:",nullable"`
	BirthDay *time.Time
	Created  time.Time
	Updated  time.Time
}
