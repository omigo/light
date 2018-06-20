package model

import (
	"time"

	"github.com/arstd/light/null"
)

type Status uint8

type User struct {
	Id       uint64
	Username string
	Phone    string `json:"mobile" light:"mobile,nullable"`
	Address  *string
	Status   Status
	BirthDay *time.Time
	Created  time.Time
	Updated  null.Timestamp `light:",nullable"`
}
