package store

import (
	"database/sql"

	"github.com/arstd/light/example/model"
)

//go:generate light -log

var User IUser

type IUser interface {

	// CREATE TABLE IF NOT EXISTS #{name} (
	// 	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	// 	username VARCHAR(32) NOT NULL UNIQUE,
	// 	Phone VARCHAR(32),
	// 	address VARCHAR(256),
	// 	status TINYINT UNSIGNED,
	// 	birth_day DATE,
	// 	created TIMESTAMP default CURRENT_TIMESTAMP,
	// 	updated TIMESTAMP default CURRENT_TIMESTAMP
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8
	Create(name string) error

	// insert ignore into users(`username`, phone, address, status, birth_day, created, updated)
	// values (${u.Username},?,?,?,?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	Insert(tx *sql.Tx, u *model.User) (int64, error)

	// insert into users(username, phone, address, status, birth_day, created, updated)
	// values (?,?,?,?,?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	// on duplicate key update
	//   username=values(username), phone=values(phone), address=values(address),
	//   status=values(status), birth_day=values(birth_day), updated=CURRENT_TIMESTAMP
	Upsert(u *model.User, tx *sql.Tx) (int64, error)

	// replace into users(username, phone, address, status, birth_day, created, updated)
	// values (?,?,?,?,?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	Replace(u *model.User) (int64, error)

	// UPDATE users
	// SET [username=?,]
	//     [phone=?,]
	//     [address=?,]
	//     [status=?,]
	//     [birth_day=?,]
	//     updated=CURRENT_TIMESTAMP
	// WHERE id=?
	Update(u *model.User) (int64, error)

	// DELETE FROM users WHERE id=?
	Delete(id uint64) (int64, error)

	// select id, username, mobile, address, status, birth_day, created, updated
	// FROM users WHERE id=?
	Get(id uint64) (*model.User, error)

	// select count(1)
	// from users
	Count() (int64, error)

	// select (select id from users where id=a.id) as id,
	// `username`, phone as phone, address, status, birth_day, created, updated
	// from users a
	// where id != -1 and  username <> 'admin' and username like ?
	// [
	// 	and address = ?
	// 	[and phone like ?]
	// 	and created > ?
	//  [{(u.BirthDay != nil && !u.BirthDay.IsZero()) || u.Id > 1 }
	//   [and birth_day > ?]
	//   [and id > ?]
	//  ]
	// ]
	// and status != ?
	// [and updated > ?]
	// and birth_day is not null
	// order by updated desc
	// limit ${offset}, ${size}
	List(u *model.User, offset, size int) ([]*model.User, error)

	// select id, username, phone, address, status, birth_day, created, updated
	// from users
	// where username like ?
	// [
	// 	and address = ?
	// 	[and phone like ?]
	// 	and created > ?
	// ]
	// and birth_day is not null
	// and status != ?
	// [{ range } and status in (#{ss})]
	// [and updated > ?]
	// order by updated desc
	// limit ${offset}, ${size}
	Page(u *model.User, ss []model.Status, offset int, size int) (int64, []*model.User, error)
}
