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
func (*UserStore) List(u *model.User, offset int, size int) ([]*model.User, error) {
	var buf bytes.Buffer
	var args []interface{}
	buf.WriteString(`SELECT id,username,phone,address,status,birthday,created,updated FROM users WHERE username LIKE ? `)
	args = append(args, u.Username)
	if u.Phone != "" {
		buf.WriteString(`AND address = ?`)
		args = append(args, u.Address)
		if u.Phone != "" {
			buf.WriteString(`AND phone LIKE ? `)
			args = append(args, light.String(&u.Phone))
		}
		buf.WriteString(`AND created > ? `)
		args = append(args, u.Created)
	}
	buf.WriteString(`AND status != ? `)
	args = append(args, light.Uint8(&u.Status))
	if !u.Updated.IsZero() {
		buf.WriteString(`AND updated > ? `)
		args = append(args, u.Updated)
	}
	buf.WriteString(`AND birthday IS NOT NULL `)
	buf.WriteString(`ORDER BY updated DESC LIMIT ?, ?`)
	args = append(args, offset, size)
	query := buf.String()
	log.Debug(query)
	log.Debug(args...)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return nil, err
	}
	defer rows.Close()
	var data []*model.User
	for rows.Next() {
		xu := new(model.User)
		data = append(data, xu)
		xdst := []interface{}{&xu.Id, &xu.Username, light.String(&xu.Phone), &xu.Address, light.Uint8(&xu.Status), &xu.Birthday, &xu.Created, &xu.Updated}
		err = rows.Scan(xdst...)
		if err != nil {
			log.Error(query)
			log.Error(args)
			log.Error(err)
			return nil, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return nil, err
	}
	return data, nil
}
