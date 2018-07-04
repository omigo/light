package store

import (
	"database/sql/driver"
	"strings"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/arstd/light/example/model"
	"github.com/arstd/light/null"
	"github.com/arstd/log"
)

var mock sqlmock.Sqlmock

func init() {
	var err error
	db, mock, err = sqlmock.New()
	log.Fataln(err)
	// defer db.Close()
}

func TestUserCreate(t *testing.T) {
	mock.ExpectExec("CREATE TABLE").WillReturnResult(sqlmock.NewResult(0, 0))

	err := User.Create("users")
	if err != nil {
		t.Error(err)
	}
}

func TestUserInsert(t *testing.T) {
	mock.ExpectBegin()
	mock.ExpectExec("INSERT ").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// mock.ExpectRollback()

	username := "admin" + time.Now().Format("150405")
	u := &model.User{
		Username: username,
		Phone:    username,
	}
	tx, err := db.Begin()
	if err != nil {
		t.Fatal(err)
	}
	// defer tx.Rollback()
	id0, err := User.Insert(tx, u)
	if err != nil {
		t.Error(err)
	}
	tx.Commit()
	if id0 == 0 {
		t.Errorf("expect id > 1, but %d", id0)
	}
}

func TestUserBulky(t *testing.T) {
	mock.ExpectBegin()
	stmt := mock.ExpectPrepare("INSERT ")
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	stmt.ExpectExec().WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	// mock.ExpectRollback()

	us := []*model.User{
		{
			Username: "admin1" + time.Now().Format("150405"),
			Phone:    "admin2" + time.Now().Format("150405"),
		},
		{
			Username: "admin1" + time.Now().Format("150405"),
			Phone:    "admin2" + time.Now().Format("150405"),
		},
	}

	affect, _, err := User.Bulky(us)
	if err != nil {
		t.Error(err)
	}
	if affect <= 1 {
		t.Errorf("expect affect > 1, but %d", affect)
	}
}

func TestUserUpsert(t *testing.T) {
	mock.ExpectBegin()
	args := []driver.Value{sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()}
	mock.ExpectExec("INSERT INTO").WithArgs(args...).WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()
	// mock.ExpectRollback()

	username := "admin" + time.Now().Format("150405")
	u := &model.User{
		Username: username,
		Phone:    username,
	}
	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Rollback()
	id0, err := User.Upsert(u, tx)
	if err != nil {
		t.Error(err)
	}
	tx.Commit()
	if id0 != 0 {
		t.Errorf("expect id = 0, but %d", id0)
	}
}

func TestUserReplace(t *testing.T) {
	mock.ExpectExec("REPLACE INTO").WillReturnResult(sqlmock.NewResult(1, 2))

	u := &model.User{
		Username: "admin" + time.Now().Format("150405"),
	}
	id0, err := User.Replace(u)
	if err != nil {
		t.Error(err)
	}
	if id0 == 0 {
		t.Errorf("expect id > 1, but %d", id0)
	}
}

func TestUserUpdate(t *testing.T) {
	mock.ExpectExec("UPDATE").WillReturnResult(sqlmock.NewResult(0, 1))

	addr := "address3"
	birth := time.Now()
	u := &model.User{
		Id:       1,
		Username: "admin3" + time.Now().Format("150405"),
		Phone:    "phone3",
		Address:  &addr,
		Status:   3,
		BirthDay: &birth,
	}
	a, err := User.Update(u)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect affect 1 rows, but %d", a)
	}
}

func TestUserGet(t *testing.T) {
	columns := strings.Split("id, username, phone, address, status, birth_day, created, updated", ", ")
	returns := []driver.Value{int64(1), []byte("admin"), []byte("13812341234"),
		[]byte("Pudong"), int64(1), time.Now(), time.Now(), time.Now()}
	rows := sqlmock.NewRows(columns).AddRow(returns...)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	if u, err := User.Get(1); err != nil {
		t.Error(err)
	} else if u == nil {
		t.Error("expect get one record, but not")
	} else if u.Username != "admin" {
		t.Errorf("expect username=admin, but got %s", u.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUserList(t *testing.T) {
	columns := strings.Split("id, username, phone, address, status, birth_day, created, updated", ", ")
	returns := []driver.Value{int64(1), []byte("admin"), []byte("13812341234"),
		[]byte("Pudong"), int64(1), time.Now(), time.Now(), time.Now()}
	rows := sqlmock.NewRows(columns).AddRow(returns...)
	mock.ExpectQuery("SELECT").WithArgs(1).WillReturnRows(rows)

	if u, err := User.Get(1); err != nil {
		t.Error(err)
	} else if u == nil {
		t.Error("expect get one record, but not")
	} else if u.Username != "admin" {
		t.Errorf("expect username=admin, but got %s", u.Username)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error(err)
	}
}

func TestUserPage(t *testing.T) {
	count := sqlmock.NewRows([]string{"count"}).AddRow(int64(10))
	mock.ExpectQuery("SELECT").WillReturnRows(count)

	columns := strings.Split("id, username, phone, address, status, birth_day, created, updated", ", ")
	returns := []driver.Value{int64(1), []byte("admin"), []byte("13812341234"),
		[]byte("Pudong"), int64(1), time.Now(), time.Now(), time.Now()}
	rows := sqlmock.NewRows(columns).AddRow(returns...).AddRow(returns...)
	mock.ExpectQuery("SELECT").WillReturnRows(rows)

	update := time.Now().Add(-time.Hour)
	u := &model.User{
		Username: "ad%",
		Updated:  null.Timestamp{Time: &update},
		Status:   9,
	}
	total, data, err := User.Page(u, []model.Status{1, 2, 3}, 1, 2)
	if err != nil {
		log.Error(err)
	}
	if total == 0 || len(data) == 0 {
		t.Error("expect get one or more records, but not")
	}
}

func TestUserDelete(t *testing.T) {
	mock.ExpectExec("DELETE").WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 1))

	a, err := User.Delete(1)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect affect 1 rows, but %d", a)
	}
}
