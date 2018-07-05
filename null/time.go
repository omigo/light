package null

import (
	"bytes"
	"database/sql/driver"
	"reflect"
	"time"
)

const (
	formatDate     = `"2006-01-02"`
	formatDatetime = `"2006-01-02 15:04:05"`
)

type NullTime struct {
	Time *time.Time
}

func (n *NullTime) IsEmpty() bool {
	return n.Time == nil || n.Time.IsZero()
}

func (n *NullTime) MarshalJSON() ([]byte, error) {
	if n.Time == nil || n.Time.IsZero() {
		return []byte("null"), nil
	}
	if n.Time.Hour() == 0 && n.Time.Minute() == 0 && n.Time.Second() == 0 {
		return []byte(n.Time.Format(formatDate)), nil
	}
	return []byte(n.Time.Format(formatDatetime)), nil
}

func (n *NullTime) UnmarshalJSON(data []byte) error {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return err
	}
	if n.Time == nil {
		n.Time = new(time.Time)
	}
	if bytes.Equal(data, []byte("0000-00-00")) || bytes.Equal(data, []byte("0000-00-00 00:00:00")) {
		var tmp time.Time
		*n.Time = tmp
		return nil
	}
	if len(data) == len(formatDate) {
		*n.Time, err = time.ParseInLocation(formatDate, string(data), loc)
	} else {
		*n.Time, err = time.ParseInLocation(formatDatetime, string(data), loc)
	}
	return err
}

func (n *NullTime) String() string {
	if n.Time == nil {
		return "nil"
	}
	if n.Time.IsZero() {
		return "nil"
	}

	return n.Time.Format("2006-01-02 15:04:05.999")
}

func (n *NullTime) Scan(value interface{}) error {
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
func (n NullTime) Value() (driver.Value, error) {
	if n.Time == nil {
		return nil, nil
	}

	if n.Time.IsZero() {
		return nil, nil
	}

	return *n.Time, nil
}
