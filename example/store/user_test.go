package store

import (
	"testing"
	"time"

	"github.com/arstd/light/example/model"
	"github.com/arstd/log"
)

var store = &UserStore{}

var id uint64

func TestUserInsert(t *testing.T) {
	u := &model.User{
		Username: "admin" + time.Now().Format("150405"),
	}
	a, err := store.Insert(u)
	if err != nil {
		t.Error(err)
	}
	log.Json(a)
	id = uint64(a)
}

func TestUserDelete(t *testing.T) {
	a, err := store.Delete(id)
	if err != nil {
		t.Error(err)
	}
	log.Json(a)
}

func TestUserInsert2(t *testing.T) {
	addr := "address"
	birth := time.Now()
	u := &model.User{
		Username: "admin2" + time.Now().Format("150405"),
		Phone:    "phone",
		Address:  &addr,
		Status:   2,
		Birthday: &birth,
	}
	a, err := store.Insert(u)
	if err != nil {
		t.Error(err)
	}
	log.Json(a)
	id = uint64(a)
}

func TestUserDelete2(t *testing.T) {
	a, err := store.Delete(id)
	if err != nil {
		t.Error(err)
	}
	log.Json(a)
}

func TestUserUpdate(t *testing.T) {
	addr := "address3"
	birth := time.Now()
	u := &model.User{
		Id:       1,
		Username: "admin3" + time.Now().Format("150405"),
		Phone:    "phone3",
		Address:  &addr,
		Status:   3,
		Birthday: &birth,
	}
	a, err := store.Update(u)
	if err != nil {
		t.Error(err)
	}
	log.Json(a)
}

func TestUserGet(t *testing.T) {
	var id uint64 = 1
	u, err := store.Get(id)
	if err != nil {
		t.Error(err)
	}
	log.Json(u)
}

func TestUserGet2(t *testing.T) {
	var id uint64 = 2
	u, err := store.Get(id)
	if err != nil {
		log.Error(err)
	}
	log.Json(u)
}

func TestUserList(t *testing.T) {
	u := &model.User{
		Username: "ad%",
		Updated:  time.Now().Add(-time.Hour),
		Status:   0,
	}
	data, err := store.List(u, 0, 2)
	if err != nil {
		log.Error(err)
	}
	log.Json(data)
}

func TestUserPage(t *testing.T) {
	u := &model.User{
		Username: "ad%",
		Updated:  time.Now().Add(-time.Hour),
	}
	total, data, err := store.Page(u, 0, 1)
	if err != nil {
		log.Error(err)
	}
	log.Json(total, data)
}
