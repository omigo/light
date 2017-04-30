package mapper

import (
	"database/sql"
	"time"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
)

//go:generate light -skip=false

// ModelMapper 示例接口
type ModelMapper interface {
	// insert into models(name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Array}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// returning id
	Insert(m *domain.Model, tx ...*sql.Tx) error

	// insert into models(name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32)
	// values [{ i, m := range ms | , }
	//  (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Array},
	//  ${m.Slice}, ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// ]
	BatchInsert(ms []*domain.Model, tx ...*sql.Tx) (int64, error)

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where id=${id}
	Get(id int, tx ...*sql.Tx) (*domain.Model, error)

	// update models
	// set name=${m.Name}, flag=${m.Flag}, score=${m.Score},
	//   map=${m.Map}, time=${m.Time}, slice=${m.Slice},
	//   status=${m.Status}, pointer=${m.Pointer}, struct_slice=${m.StructSlice},
	//   uint32=${m.Uint32}
	// where id=${m.Id}
	Update(m *domain.Model, tx ...*sql.Tx) (int64, error)

	// delete from models
	// where id=${id}
	Delete(id int, tx ...*sql.Tx) (int64, error)

	// select count(*)
	// from models
	// where name like ${m.Name}
	// [{ m.Flag } and flag=${m.Flag} ]
	// [{ len(m.Array) != 0 } and xarray && array[ [{range m.Array}] ] ]
	// [{ len(ss) != 0 } and status in ( [{range ss}] ) ]
	// [{ len(m.Slice) != 0 } and slice && ${m.Slice} ]
	Count(m *domain.Model, ss []enum.Status, tx ...*sql.Tx) (int64, error)

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where name like ${m.Name}
	// [ [ and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag}
	//   [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// ]
	// [ and time between ${from} and ${to} ]
	// [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// [{ from.IsZero() && !to.IsZero() } and time <= ${to} ]
	// [ and xarray && array[ [{range m.Array}] ] ]
	// [ and slice && ${m.Slice} ]
	// order by id
	// offset ${offset} limit ${limit}
	List(m *domain.Model, ss []enum.Status, from, to time.Time, offset, limit int, tx ...*sql.Tx) ([]*domain.Model, error)

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where name like ${m.Name}
	// [{ m.Flag != false }
	//   [{ len(ss) != 0 } and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag} ]
	// [{ len(m.Slice) != 0 } and slice && ${m.Slice} ]
	// [ and time between ${from} and ${to} ]
	// [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// [{ from.IsZero() && !to.IsZero() } and time <= ${to} ]
	// order by #{orderBy}
	// offset ${offset} limit ${limit}
	Page(m *domain.Model, ss []enum.Status, from, to time.Time, orderBy string, offset, limit int, tx ...*sql.Tx) (int64, []*domain.Model, error)
}
