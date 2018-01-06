package store

import (
	"bytes"

	"github.com/arstd/light/example/model"
	"github.com/arstd/light/null"
	"github.com/arstd/log"
)

// type User struct {
// 	Id       uint64     `db:"id BIGINT UNSIGNED AUTO_INCREMENT"`
// 	Username string     `db:"username VARCHAR(32) NOT NULL UNIQUE"`
// 	Phone    string     `db:"Phone VARCHAR(32)"`
// 	Address  *string    `db:"address VARCHAR(256)"`
// 	Status   Status     `db:"status TINYINT UNSIGNED"`
// 	Birthday *time.Time `db:"birthday DATE"`
// 	Created  time.Time  `db:"created TIMESTAMP default CURRENT_TIMESTAMP"`
// 	Updated  time.Time  `db:"updated TIMESTAMP default CURRENT_TIMESTAMP"`
// }

type UserStore0 struct{}

// INSERT INTO users(username, phone, address, status, birthday, created, updated)
// VALUES (${u.Username}, ${u.Phone}, ${u.Address}, ${u.Status}, ${u.Birthday},
//   CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
func (*UserStore0) Insert0(u *model.User) (int64, error) {
	query := `insert into users(username, phone, address, status, birthday, created,
		 updated) values (?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`
	args := []interface{}{u.Username, null.String(&u.Phone), u.Address,
		null.Uint8(&u.Status), u.Birthday}
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

// UPDATE users
// SET [username=${u.Username},]
//     [phone=${u.Phone},]
//     [address=${u.Address},]
//     [status=${u.Status},]
//     [birthday=${u.Birthday},]
//     updated=CURRENT_TIMESTAMP
// WHERE id=${u.Id}
func (*UserStore0) Update0(u *model.User) (int64, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	args := make([]interface{}, 0, 64)
	buf.WriteString(`UPDATE users SET `)
	if u.Phone != "" {
		buf.WriteString(`phone=?, `)
		args = append(args, null.String(&u.Phone))
	}
	if u.Address != nil {
		buf.WriteString(`address=?, `)
		args = append(args, u.Address)
	}
	if u.Status != 0 {
		buf.WriteString(`status=?, `)
		args = append(args, null.Uint8(&u.Status))
	}
	if u.Birthday != nil {
		buf.WriteString(`birthday=?, `)
		args = append(args, u.Birthday)
	}
	buf.WriteString(`updated=CURRENT_TIMESTAMP WHERE id=? `)
	args = append(args, u.Id)

	query := buf.String()
	log.Debug(query)
	log.Debug(args)

	res, err := db.Exec(query, args...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

// DELETE FROM users WHERE id=${id}
func (*UserStore0) Delete0(id uint64) (int64, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	args := make([]interface{}, 0, 64)
	buf.WriteString(`DELETE FROM users WHERE id=? `)
	args = append(args, id)

	query := buf.String()
	log.Debug(query)
	log.Debug(args)

	res, err := db.Exec(query, args...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, err
	}
	return res.RowsAffected()
}

// SELECT *
// FROM users WHERE id=${id}
func (*UserStore0) Get0(id uint64) (*model.User, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	args := make([]interface{}, 0, 64)
	buf.WriteString(`SELECT id, username, phone, address, status, birthday,
		created, updated FROM users WHERE id=? `)
	args = append(args, id)

	query := buf.String()
	log.Debug(query)
	log.Debug(args)

	row := db.QueryRow(query, args...)
	xu := new(model.User)
	xdst := []interface{}{&xu.Id, &xu.Username, null.String(&xu.Phone),
		&xu.Address, null.Uint8(&xu.Status), &xu.Birthday, &xu.Created, &xu.Updated}
	err := row.Scan(xdst...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return nil, err
	}
	return xu, nil
}

// select *
// from users
// where 1=1
// [and username like ${u.Username}]
// [and phone like ${u.Phone}]
// [and updated > ${u.Updated}]
// limit 10
func (*UserStore0) List0(u *model.User) ([]*model.User, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	args := make([]interface{}, 0, 64)
	buf.WriteString(`SELECT id, username, phone, address, status, birthday,
		created, updated FROM users WHERE 1=1 `)
	if u.Username != "" {
		buf.WriteString(`and username like ? `)
		args = append(args, u.Username)
	}
	if u.Phone != "" {
		buf.WriteString(`and phone like ? `)
		args = append(args, null.String(&u.Phone))
	}
	if !u.Updated.IsZero() {
		buf.WriteString(`and updated > ? `)
		args = append(args, u.Updated)
	}
	buf.WriteString(`limit 10`)

	query := buf.String()
	log.Debug(query)
	log.Debug(args)

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
		xdst := []interface{}{&xu.Id, &xu.Username, null.String(&xu.Phone),
			&xu.Address, null.Uint8(&xu.Status), &xu.Birthday, &xu.Created, &xu.Updated}
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

// select *
// from users
// where 1=1
// [and username like ${u.Username}]
// [and phone like ${u.Phone}]
// [and updated > ${u.Updated}]
// limit ${(page-1)*size}, ${size}
func (*UserStore0) Page0(u *model.User, page int, size int) (int64, []*model.User, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	args := make([]interface{}, 0, 64)
	listStmt := "SELECT id, username, phone, address, status, birthday, created, updated "
	totalStmt := "select count(1) "
	buf.WriteString(`FROM users WHERE 1=1 `)
	if u.Username != "" {
		buf.WriteString(`and username like ? `)
		args = append(args, u.Username)
	}
	if u.Phone != "" {
		buf.WriteString(`and phone like ? `)
		args = append(args, null.String(&u.Phone))
	}
	if !u.Updated.IsZero() {
		buf.WriteString(`and updated > ? `)
		args = append(args, u.Updated)
	}

	query := totalStmt + buf.String()
	log.Debug(query)
	log.Debug(args)

	var total int64
	err := db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, nil, err
	}

	buf.WriteString(`limit ?, ? `)
	args = append(args, (page-1)*size, size)

	query = listStmt + buf.String()
	log.Debug(query)
	log.Debug(args)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, nil, err
	}
	defer rows.Close()

	var data []*model.User
	for rows.Next() {
		xu := new(model.User)
		data = append(data, xu)
		xdst := []interface{}{&xu.Id, &xu.Username, null.String(&xu.Phone),
			&xu.Address, null.Uint8(&xu.Status), &xu.Birthday, &xu.Created, &xu.Updated}
		err = rows.Scan(xdst...)
		if err != nil {
			log.Error(query)
			log.Error(args)
			log.Error(err)
			return 0, nil, err
		}
	}
	if err = rows.Err(); err != nil {
		log.Error(query)
		log.Error(args)
		log.Error(err)
		return 0, nil, err
	}

	return total, data, nil
}
