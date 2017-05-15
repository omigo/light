package light

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

func Time(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	switch v := a.(type) {
	case uint32:
		return &TimeWapper{Uint32: &v, Time: time.Unix(int64(v), 0), Valid: true}
	case int32:
		return &TimeWapper{Int32: &v, Time: time.Unix(int64(v), 0), Valid: true}
	case int:
		return &TimeWapper{Int: &v, Time: time.Unix(int64(v), 0), Valid: true}
	case uint64:
		return &TimeWapper{Uint64: &v, Time: time.Unix(int64(v), 0), Valid: true}
	case int64:
		return &TimeWapper{Int64: &v, Time: time.Unix(v, 0), Valid: true}

	case *uint32:
		return &TimeWapper{Uint32: v}
	case *int32:
		return &TimeWapper{Int32: v}
	case *int:
		return &TimeWapper{Int: v}
	case *uint64:
		return &TimeWapper{Uint64: v}
	case *int64:
		return &TimeWapper{Int64: v}

	default:
		panic("type not implemented")
	}
}

type TimeWapper struct {
	Uint32 *uint32
	Int32  *int32
	Int    *int
	Uint64 *uint64
	Int64  *int64

	Time  time.Time
	Valid bool
}

func (b TimeWapper) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Time, nil
}

func (b *TimeWapper) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	b.Time, b.Valid = src.(time.Time)
	if b.Valid {
		switch {
		case b.Uint32 != nil:
			*b.Uint32 = uint32(b.Time.Unix())
		case b.Int32 != nil:
			*b.Int32 = int32(b.Time.Unix())
		case b.Int != nil:
			*b.Int = int(b.Time.Unix())
		case b.Uint64 != nil:
			*b.Uint64 = uint64(b.Time.Unix())
		case b.Int64 != nil:
			*b.Int64 = b.Time.Unix()
		}
	}
	return nil
}
