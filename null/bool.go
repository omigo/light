package null

import (
	"database/sql/driver"
	"fmt"
	"reflect"
)

type NullBool struct {
	Bool *bool
}

func (n *NullBool) IsEmpty() bool {
	return n.Bool == nil || *n.Bool
}

func (n *NullBool) MarshalJSON() ([]byte, error) {
	if n.Bool == nil {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("%t", *n.Bool)), nil
}

func (n *NullBool) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	if string(data) == "true" {
		var b bool = true
		*n.Bool = b
	}
	return nil
}

func (n *NullBool) String() string {
	if n.Bool != nil {
		return "nil"
	}
	if *n.Bool {
		return "true"
	}
	return "false"
}

// Scan implements the Scanner interface.
func (s *NullBool) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s.Bool = v == 1
	case *int64:
		*s.Bool = *v == 1
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}
	return nil
}

// Value implements the driver Valuer interface.
func (s NullBool) Value() (driver.Value, error) {
	if s.Bool == nil {
		return nil, nil
	}
	if *s.Bool {
		return 1, nil
	}
	return 0, nil
}
