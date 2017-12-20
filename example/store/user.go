package store

import "github.com/arstd/light/example/model"

type User interface {

	// insert into users(username, phone, address, status, birthday, created, updated)
	// values (${u.Username}, ${u.Phone}, ${u.Address}, ${u.Status}, ${u.Birthday},
	//   CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	Insert(u *model.User) (int64, error)
	//
	// // UPDATE users
	// // SET [username=${u.Username},]
	// //     [phone=${u.Phone},]
	// //     [address=${u.Address},]
	// //     [status=${u.Status},]
	// //     [birthday=${u.Birthday},]
	// //     updated=CURRENT_TIMESTAMP
	// // WHERE id=${u.Id}
	// Update(u *model.User) (int64, error)
	//
	// // DELETE FROM users WHERE id=${id}
	// Delete(id uint64) (int64, error)
	//
	// // SELECT *
	// // FROM users WHERE id=${id}
	// Get(id uint64) (*model.User, error)
	//
	// // select *
	// // from users
	// // where 1=1
	// // [and username like ${u.Username}]
	// // [and phone like ${u.Phone}]
	// // [and updated > ${u.Updated}]
	// // limit 20,10
	// List(u *model.User) ([]*model.User, error)
	//
	// // select *
	// // from users
	// // where username like ${u.Username}
	// // [and phone like ${u.Phone}]
	// // and status != ${u.Status}
	// // [and updated > ${u.Updated}]
	// // limit ${(page-1)*size}, ${size}
	// Page(u *model.User, page int, size int) (int64, []*model.User, error)
}
