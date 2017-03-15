package testdata

import (
	"database/sql"

	"github.com/arstd/light/example/model"
	"github.com/arstd/yan/example/enum"
)

//go:generate go run ../main.go

// Interface1 示例接口
type Interface1 interface {

	// insert into model(id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32)
	// values [{ i, m := range ms | , } (${i}, ${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32}) ]
	BatchInsert(tx *sql.Tx, ms []*model.Model) (int64, error)

	// insert into model(name, flag, score, map, time, slice, status, pointer, struct_slice, uint32)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// returning id
	Insert(tx *sql.Tx, m *model.Model) error

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32
	// from model
	// where name like ${m.Name}
	//   [{ m.Flag != false }
	//       [{ m.Status != 0 } status=${m.Status} ]
	//       and flag=${m.Flag} ]
	//   [{ len(ss) != 0 } and status in ([{range ss}]) ]
	//   [{ len(m.Slice) != 0 } and slice ?| array[[{range m.Slice}]] ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *model.Model, ss []enum.Status, offset, limit int) ([]*model.Model, error)
}
