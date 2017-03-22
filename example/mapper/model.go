package mapper

import (
	"database/sql"

	"github.com/arstd/light/example/enum"
	"github.com/arstd/light/example/model"
)

//go:generate go run ../../main.go

// ModelMapper 示例接口
type ModelMapper interface {

	// insert into models(name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Array}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// returning id
	Insert(trans *sql.Tx, m *model.Model) error

	// insert into models(uint32, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice)
	// values [{ i, m := range ms | , }
	//  (${i}+888, ${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Array},
	//  ${m.Slice}, ${m.Status}, ${m.Pointer}, ${m.StructSlice})
	// ]
	BatchInsert(tx *sql.Tx, ms []*model.Model) (int64, error)

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where id=${id}
	Get(tx *sql.Tx, id int) (*model.Model, error)

	// update models
	// set name=${m.Name}, flag=${m.Flag}, score=${m.Score},
	//   map=${m.Map}, time=${m.Time}, slice=${m.Slice},
	//   status=${m.Status}, pointer=${m.Pointer}, struct_slice=${m.StructSlice},
	//   uint32=${m.Uint32}
	// where id=${m.Id}
	Update(tx *sql.Tx, m *model.Model) (int64, error)

	// delete from models
	// where id=${id}
	Delete(tx *sql.Tx, id int) (int64, error)

	// select count(*)
	// from models
	// where name like ${m.Name}
	// [{ m.Flag != false } and flag=${m.Flag} ]
	// [{ len(ss) != 0 } and status in ( [{range ss}] ) ]
	// [{ len(m.Slice) != 0 } and slice && ${m.Slice} ]
	Count(tx *sql.Tx, m *model.Model, ss []enum.Status) (int64, error)

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where name like ${m.Name}
	// [{ m.Flag != false }
	//   [{ len(ss) != 0 } and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag} ]
	// [{ len(m.Array) != 0 } and xarray && array[ [{range m.Array}] ] ]
	// [{ len(m.Slice) != 0 } and slice && ${m.Slice} ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *model.Model, ss []enum.Status, offset, limit int) ([]*model.Model, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice
	// from models
	// where name like ${m.Name}
	// [{ m.Flag != false }
	//   [{ len(ss) != 0 } and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag} ]
	// [{ len(m.Slice) != 0 } and slice && ${m.Slice} ]
	// order by id
	// offset ${offset} limit ${limit}
	Page(tx *sql.Tx, m *model.Model, ss []enum.Status, offset, limit int) (int64, []*model.Model, error)
}
