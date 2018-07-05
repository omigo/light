package null

import (
	"bytes"
	"database/sql/driver"
	"reflect"
	"time"
)

var zero = time.Unix(0, 0)

type ClickHouseTime struct {
	Time *time.Time
}

func (n *ClickHouseTime) IsEmpty() bool {
	return n.Time == nil || n.Time.IsZero()
}

func (n *ClickHouseTime) MarshalJSON() ([]byte, error) {
	if n.Time == nil || n.Time.IsZero() {
		return []byte("null"), nil
	}
	if n.Time.Hour() == 0 && n.Time.Minute() == 0 && n.Time.Second() == 0 {
		return []byte(n.Time.Format(formatDate)), nil
	}
	return []byte(n.Time.Format(formatDatetime)), nil
}

func (n *ClickHouseTime) UnmarshalJSON(data []byte) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	if n.Time == nil {
		n.Time = new(time.Time)
	}
	if bytes.Equal(data, []byte("0000-00-00")) || bytes.Equal(data, []byte("0000-00-00 00:00:00")) {
		*n.Time = zero
		return nil
	}
	if len(data) == len(formatDate) {
		*n.Time, err = time.ParseInLocation(formatDate, string(data), loc)
	} else {
		*n.Time, err = time.ParseInLocation(formatDatetime, string(data), loc)
	}
	return err
}

func (n *ClickHouseTime) String() string {
	if n.Time == nil {
		return "1970-01-01 08:00:00"
	}
	if n.Time.IsZero() {
		return "1970-01-01 08:00:00"
	}

	return n.Time.Format("2006-01-02 15:04:05")
}

func (n *ClickHouseTime) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	if n.Time == nil {
		n.Time = new(time.Time)
	}

	switch v := value.(type) {
	case time.Time:
		*n.Time = v

	case *time.Time:
		*n.Time = *v

	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}
	return nil
}
func (n ClickHouseTime) Value() (driver.Value, error) {
	if n.Time == nil {
		return nil, nil
	}

	if n.Time.IsZero() {
		return nil, nil
	}

	return *n.Time, nil
}
