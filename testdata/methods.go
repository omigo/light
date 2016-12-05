package testdata

import (
	"database/sql"

	"github.com/arstd/light/example/domain"
	"github.com/arstd/light/example/enum"
)

//go:generate yan -force

// Interface1 示例接口
type Interface1 interface {

	// insert into model(name, flag, score, map, time, slice, status, pointer, struct_slice, uint32)
	// values (${m.Name}, ${m.Flag}, ${m.Score}, ${m.Map}, ${m.Time}, ${m.Slice},
	//   ${m.Status}, ${m.Pointer}, ${m.StructSlice}, ${m.Uint32})
	// returning id
	Insert(tx *sql.Tx, m *domain.Model) error

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32
	// from model
	// where name like ${m.Name}
	//   [?{ m.Flag != false } and flag=${m.Flag} ]
	//   [?{ len(ss) != 0 } and status in (${ss}) ]
	//   [?{ len(m.Slice) != 0 } and slice ?| array[${m.Slice}] ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, m *domain.Model, ss []enum.Status, offset, limit int) ([]*domain.Model, error)
}
