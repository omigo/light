package null

import (
	"database/sql/driver"
	"reflect"
	"time"
)

type NullTime struct {
	Time *time.Time
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
