package mapper

import (
	"testing"
	"time"

	"github.com/arstd/log"

	m "github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
)

func TestCreateTable(t *testing.T) {
	_, err := db.Exec("drop table if exists model")
	if err != nil {
		log.Error(err)
	}
	_, err = db.Exec(`
		create table model (
			id serial primary key,
			name text not null default '',
			flag bool not null default false,
			score numeric(20,2) not null default 0.0,

			map jsonb not null default '{}',
			time timestamptz default current_timestamp,
			slice jsonb not null default '[]',

			status smallint not null default 0,
			pointer jsonb not null default '{}',
			struct_slice jsonb not null default '[]',

			uint32 timestamptz default current_timestamp
		)
	`)
	if err != nil {
		log.Error(err)
	}
}

var x ModelMapper = &ModelMapperImpl{}
var id int = 5

func TestModelMapperInsert(t *testing.T) {
	m := &m.Model{
		Name:  "name",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": 1},
		Time:  time.Now().Add(-3 * time.Hour),
		Slice: []string{"Slice Elem 1", "Slice Elem 2"},

		Status:  enum.StatusNormal,
		Pointer: &m.Model{Name: "Pointer"},
		StructSlice: []*m.Model{
			{Name: "StructSlice"},
		},

		Uint32: uint32(time.Now().Add(-5 * time.Minute).Unix()),
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	err = x.Insert(tx, m)
	if err != nil {
		t.Fatalf("insert error: %s", err)
	}

	CommitTx(tx)
	id = m.Id
	log.Infof("id=%d", m.Id)
}

func TestModelMapperUpdate(t *testing.T) {
	m := &m.Model{
		Id:    id,
		Name:  "name update",
		Flag:  true,
		Score: 1.23,

		Map:   map[string]interface{}{"a": "1  update"},
		Time:  time.Now().Add(-3 * time.Hour),
		Slice: []string{"Slice Elem 1 update", "Slice Elem 2 update"},

		Status:  enum.StatusNormal,
		Pointer: &m.Model{Name: "Pointer update"},
		StructSlice: []*m.Model{
			{Name: "StructSlice update"},
		},
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	a, err := x.Update(tx, m)
	if err != nil {
		t.Fatalf("update error: %s", err)
	}

	CommitTx(tx)
	log.Infof("affected=%d", a)
}

func TestModelMapperGet(t *testing.T) {
	tx, err := BeginTx()
	defer RollbackTx(tx)
	m, err := x.Get(tx, id)
	if err != nil {
		t.Fatalf("get error: %s", err)
	}

	CommitTx(tx)
	log.JSON(m)
}

func TestModelMapperCount(t *testing.T) {
	m := &m.Model{
		Name:   "name%", // like 'name%'
		Flag:   true,
		Status: enum.StatusNormal,
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	count, err := x.Count(tx, m, []enum.Status{enum.StatusNormal, enum.StatusDeleted})
	if err != nil {
		t.Fatalf("count(%+v) error: %s", m, err)
	}

	CommitTx(tx)
	log.JSON(count)
}

func TestModelMapperList(t *testing.T) {
	m := &m.Model{
		Name:   "name%", // like 'name%'
		Flag:   true,
		Slice:  []string{"Slice Elem 1 update error", "Slice Elem 2 update"},
		Status: enum.StatusNormal,
	}
	tx, err := BeginTx()
	defer RollbackTx(tx)
	ms, err := x.List(tx, m, []enum.Status{enum.StatusNormal, enum.StatusDeleted}, 0, 20)
	if err != nil {
		t.Fatalf("list(%+v) error: %s", m, err)
	}

	CommitTx(tx)
	log.JSON(ms)
}

func TestModelMapperDelete(t *testing.T) {
	tx, err := BeginTx()
	defer RollbackTx(tx)
	a, err := x.Delete(tx, id)
	CommitTx(tx)

	if err != nil {
		t.Fatalf("delete error: %s", err)
	}

	log.JSON(a)
}
