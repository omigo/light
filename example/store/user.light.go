package store

import (
	"bytes"

	"github.com/arstd/light/example/model"
	"github.com/arstd/light/light"
	"github.com/arstd/log"
)

type UserStore struct{}

func (*UserStore) Insert(u *model.User) (i int64, e error) {
	var buf bytes.Buffer
	var args []interface{}
	buf.WriteString("insert into users(username, phone, address, status, birthday, created, updated) values (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)")
	args = append(args, light.String(&u.Username), light.String(&u.Phone), u.Address, light.Uint8(&u.Status), u.Birthday)
	query := buf.String()
	log.Debug(query)
	log.Debug(args...)
	res, err := db.Exec(query, args...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, err
	}
	return res.LastInsertId()
}
