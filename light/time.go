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
		return &TimeWapper{Time: time.Unix(int64(v), 0), Uint32: v, Valid: v != 0}

	case *uint32:
		return &TimeWapper{Time: time.Unix(int64(*v), 0), Uint32: *v, Valid: *v != 0}

	default:
		panic("type not implemented")
	}
}

type TimeWapper struct {
	Uint32 uint32
	Time   time.Time
	Valid  bool
}

func (b TimeWapper) Value() (driver.Value, error) {
	if !b.Valid {
		return nil, nil
	}
	return b.Time, nil
}

func (b *TimeWapper) Scan(src interface{}) error {
	b.Time, b.Valid = src.(time.Time)
	if b.Valid {
		b.Uint32 = uint32(b.Time.Unix())
	}
	return nil
}
