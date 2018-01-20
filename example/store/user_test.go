package store

import (
	"testing"
	"time"

	"github.com/arstd/light/example/model"
	"github.com/arstd/log"
)

var store User = &UserStore{}

var id uint64

func TestUserCreate(t *testing.T) {
	err := store.Create("users")
	if err != nil {
		t.Error(err)
	}
}

var username string

func TestUserInsert(t *testing.T) {
	username = "admin" + time.Now().Format("150405")
	u := &model.User{
		Username: username,
	}
	id0, err := store.Insert(u)
	if err != nil {
		t.Error(err)
	}
	if id0 == 0 {
		t.Errorf("expect id > 1, but %d", id0)
	}
	id = uint64(id0)
}

func TestUserUpsert(t *testing.T) {
	u := &model.User{
		Username: username,
	}
	id0, err := store.Upsert(u)
	if err != nil {
		t.Error(err)
	}
	if id0 != 0 {
		t.Errorf("expect id = 0, but %d", id0)
	}
}

func TestUserDelete1(t *testing.T) {
	a, err := store.Delete(id)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect affect 1 rows, but %d", a)
	}
}

func TestUserReplace(t *testing.T) {
	u := &model.User{
		Username: "admin" + time.Now().Format("150405"),
	}
	id0, err := store.Replace(u)
	if err != nil {
		t.Error(err)
	}
	if id0 == 0 {
		t.Errorf("expect id > 1, but %d", id0)
	}
	id = uint64(id0)
}

func TestUserUpdate(t *testing.T) {
	addr := "address3"
	birth := time.Now()
	u := &model.User{
		Id:       id,
		Username: "admin3" + time.Now().Format("150405"),
		Phone:    "phone3",
		Address:  &addr,
		Status:   3,
		BirthDay: &birth,
	}
	a, err := store.Update(u)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect affect 1 rows, but %d", a)
	}
}

func TestUserGet(t *testing.T) {
	u, err := store.Get(id)
	if err != nil {
		t.Error(err)
	}
	if u == nil {
		t.Error("expect get one record, but not")
	}
}

func TestUserList(t *testing.T) {
	u := &model.User{
		Username: "ad%",
		Updated:  time.Now().Add(-time.Hour),
		Status:   9,
	}
	data, err := store.List(u, 0, 2)
	if err != nil {
		log.Error(err)
	}
	if len(data) == 0 {
		t.Error("expect get one or more records, but not")
	}
}

func TestUserPage(t *testing.T) {
	u := &model.User{
		Username: "ad%",
		Updated:  time.Now().Add(-time.Hour),
		Status:   9,
	}
	total, data, err := store.Page(u, 0, 1)
	if err != nil {
		log.Error(err)
	}
	if total == 0 || len(data) == 0 {
		t.Error("expect get one or more records, but not")
	}
}

func TestUserDelete(t *testing.T) {
	a, err := store.Delete(id)
	if err != nil {
		t.Error(err)
	}
	if a != 1 {
		t.Errorf("expect affect 1 rows, but %d", a)
	}
}
