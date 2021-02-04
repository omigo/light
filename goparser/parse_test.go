package goparser

import (
	"os"
	"strings"
	"testing"

	"github.com/omigo/log"
)

func TestParse(t *testing.T) {
	gopath := strings.TrimSuffix(os.Getenv("PWD"), "/")

	t.Log(os.Getenv(""))
	filename := gopath + "/../example/store/user.go"
	src := `package store
import (
	// "database/sql"
	"github.com/omigo/light/example/model"
)
var User IUser
type IUser interface {
	// insert ignore into users(username, phone, address, status, birth_day, created, updated)
	// values (?,?,?,?,?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	// Insert(tx *sql.Tx, u *model.User) (int64, error)

	// UPDATE users
	// SET [username=?,]
	//     [phone=?,]
	//     [address=?,]
	//     [status=?,]
	//     [birth_day=?,]
	//     updated=CURRENT_TIMESTAMP
	// WHERE id=?
	// Update(u *model.User) (int64, error)

	// select id, username, phone, address, status, birth_day, created, updated
	// FROM users WHERE id=?
	Get(id uint64) (*model.User, error)
}
`

	itf, err := Parse(filename, src)
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(itf)
}
