package mapper

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
	"github.com/arstd/log"
)

func TestCreateTable(t *testing.T) {
	_, err := db.Exec("drop table if exists models")
	if err != nil {
		log.Error(err)
	}
	_, err = db.Exec(`
		create table models (
			id serial primary key,
			name text not null,
			flag bool not null default false,
			score decimal(3,1) not null default 0.0,

			map jsonb not null default '{}',
			time timestamptz not null default now(),
			xarray text[] not null,
			slice text[] not null,

			status smallint not null default 0,
			state text not null default '',

			pointer jsonb not null default '{}',
			struct_slice jsonb not null default '[]',
			uint32 timestamptz not null default now()
		)
	`)
	if err != nil {
		log.Error(err)
	}
}

var mapper ModelMapper = &ModelMapperImpl{}
var id int = 1

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

	err = mapper.Insert(m, tx)
	if err != nil {
		tx.Rollback()
		t.Fatalf("insert error: %s", err)
	}

	id = m.Id
	log.Infof("id=%d", m.Id)
}

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
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	a, err := mapper.BatchInsert([]*domain.Model{m, m, m}, tx)
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}

	CommitTx(tx)
	log.Infof("affect %d rows", a)
}

func TestModelMapperGet(t *testing.T) {
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	m, err := mapper.Get(id, tx)
	if err != nil {
		t.Fatalf("get error: %s", err)
	}

	CommitTx(tx)
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
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	a, err := mapper.Update(m, tx)
	if err != nil {
		t.Fatalf("update error: %s", err)
	}

	CommitTx(tx)
	log.Infof("affected=%d", a)
}

func TestModelMapperDelete(t *testing.T) {
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	a, err := mapper.Delete(id, tx)
	CommitTx(tx)

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
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	count, err := mapper.Count(m, []enum.Status{enum.StatusNormal, enum.StatusDeleted}, tx)
	if err != nil {
		t.Fatalf("count(%+v) error: %s", m, err)
	}

	CommitTx(tx)
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
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	ms, err := mapper.List(m, ss, time.Now().AddDate(0, 0, -1), time.Now(), 0, 20, tx)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	CommitTx(tx)
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
	tx, err := BeginTx()
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}
	defer RollbackTx(tx)
	cnt, ms, err := mapper.Page(m, ss, time.Now().AddDate(0, 0, -1), time.Now(), 0, 20, tx)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	CommitTx(tx)
	log.Info(cnt)
	log.Info(json.Marshal(ms))
}
