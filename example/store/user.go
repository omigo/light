package store

import "github.com/arstd/light/example/model"

type User interface {

	// CREATE TABLE IF NOT EXISTS users (
	// 	id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
	// 	username VARCHAR(32) NOT NULL UNIQUE,
	// 	Phone VARCHAR(32),
	// 	address VARCHAR(256),
	// 	status TINYINT UNSIGNED,
	// 	birthday DATE,
	// 	created TIMESTAMP default CURRENT_TIMESTAMP,
	// 	updated TIMESTAMP default CURRENT_TIMESTAMP
	// ) ENGINE=InnoDB DEFAULT CHARSET=utf8
	Create() error

	// insert into users(username, phone, address, status, birthday, created, updated)
	// values (${u.Username}, ${u.Phone}, ${u.Address}, ${u.Status}, ${u.Birthday},
	//   CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	Insert(u *model.User) (int64, error)

	// UPDATE users
	// SET [username=${u.Username},]
	//     [phone=${u.Phone},]
	//     [address=${u.Address},]
	//     [status=${u.Status},]
	//     [birthday=${u.Birthday},]
	//     updated=CURRENT_TIMESTAMP
	// WHERE id=${u.Id}
	Update(u *model.User) (int64, error)

	// DELETE FROM users WHERE id=${id}
	Delete(id uint64) (int64, error)

	// select id, username, phone, address, status, birthday, created, updated
	// FROM users WHERE id=${id}
	Get(id uint64) (*model.User, error)

	// select id, username, phone, address, status, birthday, created, updated
	// from users
	// where username like ${u.Username}
	// [
	// 	and address = ${u.Address}
	// 	[and phone like ${u.Phone}]
	// 	and created > ${u.Created}
	//  [{(u.Birthday != nil && !u.Birthday.IsZero()) || u.Id > 1 }
	//   [and birthday > ${u.Birthday}]
	//   [and id > ${u.Id}]
	//  ]
	// ]
	// and status != ${u.Status}
	// [and updated > ${u.Updated}]
	// and birthday is not null
	// order by updated desc
	// limit ${offset}, ${size}
	List(u *model.User, offset, size int) ([]*model.User, error)

	// select id, username, phone, address, status, birthday, created, updated
	// from users
	// where username like ${u.Username}
	// [
	// 	and address = ${u.Address}
	// 	[and phone like ${u.Phone}]
	// 	and created > ${u.Created}
	// ]
	// and birthday is not null
	// and status != ${u.Status}
	// [and updated > ${u.Updated}]
	// order by updated desc
	// limit ${page}, ${size}
	Page(u *model.User, page int, size int) (int64, []*model.User, error)
}
