package null

import (
	"database/sql/driver"
	"reflect"
	"strconv"
	"time"
)

type Timestamp struct {
	Time *time.Time
}

func (n *Timestamp) IsEmpty() bool {
	return n.Time == nil || n.Time.IsZero()
}

func (n *Timestamp) MarshalJSON() ([]byte, error) {
	if n.Time == nil || n.Time.IsZero() {
		return []byte("0"), nil
	}
	return []byte(strconv.FormatInt(n.Time.Unix(), 10)), nil
}

func (n *Timestamp) UnmarshalJSON(data []byte) error {
	if n.Time == nil {
		n.Time = new(time.Time)
	}

	ts, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}
	*n.Time = time.Unix(ts, 0)

	return err
}

func (n *Timestamp) String() string {
	if n.Time == nil {
		return "0"
	}
	if n.Time.IsZero() {
		return "0"
	}
	return strconv.FormatInt(n.Time.Unix(), 10)
}

func (n *Timestamp) Scan(value interface{}) error {
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
func (n Timestamp) Value() (driver.Value, error) {
	if n.Time == nil {
		return nil, nil
	}

	if n.Time.IsZero() {
		return nil, nil
	}

	return *n.Time, nil
}
