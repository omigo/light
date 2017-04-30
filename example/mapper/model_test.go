package mapper

import (
	"encoding/json"
	"flag"
	"testing"
	"time"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
	"github.com/arstd/log"
)

func TestInit(t *testing.T) {
	var args = flag.Args()
	if len(args) > 0 && args[0] == "pg" {
		initPG()
	} else {
		initSQLMock(t)
	}
}

var mapper ModelMapper = &ModelMapperImpl{}
var id int = 1

func TestModelMapperInsert(t *testing.T) {
	m := &domain.Model{
		Name:  "name",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": 1},
		Time:  time.Now(),
		Array: []int64{1, 2, 3},
		Slice: []string{"Slice Elem 1", "Slice Elem 2"},

		Status:  enum.StatusNormal,
		Pointer: &domain.Model{Name: "Pointer 2"},
		StructSlice: []*domain.Model{
			{Name: "StructSlice"},
		},

		Uint32: uint32(time.Now().Unix()),
	}
	err := mapper.Insert(m)
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}

	id = m.Id
	log.Infof("id=%d", m.Id)
}

func TestModelMapperInsertTx(t *testing.T) {
	m := &domain.Model{
		Name:  "name",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": 1},
		Time:  time.Now().Add(time.Hour),
		Array: []int64{1, 2, 3},
		Slice: []string{"Slice Elem 1", "Slice Elem 2"},

		Status:  enum.StatusNormal,
		Pointer: &domain.Model{Name: "Pointer"},
		StructSlice: []*domain.Model{
			{Name: "StructSlice"},
		},

		Uint32: uint32(time.Now().Unix()),
	}

	tx, err := db.Begin()
	if err != nil {
		t.Error(err)
	}
	defer tx.Commit()

	err = mapper.Insert(m)
	if err != nil {
		tx.Rollback()
		t.Fatalf("insert error: %s", err)
	}

	id = m.Id
	log.Infof("id=%d", m.Id)
}

func TestModelMapperBatchInsert(t *testing.T) {
	m := &domain.Model{
		Name:  "name",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": 1},
		Time:  time.Now(),
		Array: []int64{1, 2, 3},
		Slice: []string{"Slice Elem 1", "Slice Elem 2"},

		Status:  enum.StatusNormal,
		Pointer: &domain.Model{Name: "Pointer"},
		StructSlice: []*domain.Model{
			{Name: "StructSlice"},
		},

		Uint32: uint32(time.Now().Unix()),
	}
	a, err := mapper.BatchInsert([]*domain.Model{m, m, m})
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}

	log.Infof("affect %d rows", a)
}

func TestModelMapperGet(t *testing.T) {
	m, err := mapper.Get(id)
	if err != nil {
		t.Fatalf("get error: %s", err)
	}

	log.Info(json.Marshal(m))
}

func TestModelMapperUpdate(t *testing.T) {
	m := &domain.Model{
		Id:    id,
		Name:  "name update",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": "1  update"},
		Time:  time.Now().Add(-3 * time.Hour),
		Slice: []string{"Slice Elem 1 update", "Slice Elem 2 update"},

		Status:  enum.StatusNormal,
		Pointer: &domain.Model{Name: "Pointer update"},
		StructSlice: []*domain.Model{
			{Name: "StructSlice update"},
		},
		Uint32: 32,
	}

	a, err := mapper.Update(m)
	if err != nil {
		t.Fatalf("update error: %s", err)
	}

	log.Infof("affected=%d", a)
}

func TestModelMapperDelete(t *testing.T) {
	a, err := mapper.Delete(id)
	if err != nil {
		t.Fatalf("delete error: %s", err)
	}
	log.JSON(a)
}

func TestModelMapperCount(t *testing.T) {
	m := &domain.Model{
		Name:   "name%", // like 'name%'
		Flag:   true,
		Status: enum.StatusNormal,
		Slice:  []string{"Slice Elem 1", "Slice Elem Not Exist"},
	}
	count, err := mapper.Count(m, []enum.Status{enum.StatusNormal, enum.StatusDeleted})
	if err != nil {
		t.Fatalf("count(%+v) error: %s", m, err)
	}
	log.Info(count)
}

func TestModelMapperList(t *testing.T) {
	m := &domain.Model{
		Name:  "name%", // like 'name%'
		Flag:  true,
		Array: []int64{11, 22, 3},
		Slice: []string{"Slice Elem 1", "Slice Elem Not Exist"},
	}
	ss := []enum.Status{enum.StatusNormal, enum.StatusDeleted}
	ms, err := mapper.List(m, ss, time.Now().AddDate(0, 0, -1), time.Now(), 0, 20)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	log.Info(json.Marshal(ms))
}

func TestModelMapperPage(t *testing.T) {
	m := &domain.Model{
		Name:  "name%", // like 'name%'
		Flag:  true,
		Array: []int64{11, 22, 3},
		Slice: []string{"Slice Elem 1", "Slice Elem Not Exist"},
	}
	ss := []enum.Status{enum.StatusNormal, enum.StatusDeleted}
	cnt, ms, err := mapper.Page(m, ss, time.Now().AddDate(0, 0, -1), time.Now(), "id desc", 0, 20)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	log.Info(cnt)
	log.Info(json.Marshal(ms))
}

func TestClean(t *testing.T) {
	db.Close()
}
