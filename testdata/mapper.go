package testdata

import (
	"database/sql"

	"github.com/arstd/light/example/enum"
	"github.com/arstd/light/example/model"
)

//go:generate light -force

// ModelMapper 示例接口
type ModelMapper interface {

	// insert into model(name, flag, score, map, time, slice, status, pointer, struct_slice, uint32)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// returning id
	Insert(tx *sql.Tx, m *model.Model) error

	// update model
	// set name=${m.Name}, flag=${m.Flag}, score=${m.Score},
	//   map=${m.Map}, time=${m.Time}, slice=${m.Slice},
	//   status=${m.Status}, pointer=${m.Pointer}, struct_slice=${m.StructSlice},
	//   uint32=${m.Uint32}
	// where id=${m.Id}
	Update(tx *sql.Tx, m *model.Model) (int64, error)

	// delete from model
	// where id=${id}
	Delete(tx *sql.Tx, id int) (int64, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32
	// from model
	// where id=${id}
	Get(tx *sql.Tx, id int) (*model.Model, error)

	// select count(*)
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	Count(tx *sql.Tx, m *model.Model, ss []enum.Status) (int64, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	//   [?{ len(m.Slice) != 0 } and slice ?| array[${m.Slice}] ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *model.Model, ss []enum.Status, offset, limit int) ([]*model.Model, error)

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	//   [?{ len(m.Slice) != 0 } and slice ?| array[${m.Slice}] ]
	// order by id
	// offset ${offset} limit ${limit}
	Paging(tx *sql.Tx, m *model.Model, ss []enum.Status) (int64, []*model.Model, error)
}
