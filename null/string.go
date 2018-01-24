package null

import (
	"database/sql/driver"
	"reflect"
)

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var plain string
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&String{S:&s})
//  ...
//  use plain if database return null, plain is blank
type NullString struct {
	String_ *string
}

func (n *NullString) IsEmpty() bool {
	return n.String_ == nil || *n.String_ == ""
}

func (n *NullString) MarshalJSON() ([]byte, error) {
	if n.String_ == nil {
		return []byte("null"), nil
	}
	return []byte(`"` + *n.String_ + `"`), nil
}

func (n *NullString) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}
	*n.String_ = string(data)
	return nil
}

func (n *NullString) String() string {
	if n.String_ != nil {
		return "nil"
	}
	if *n.String_ == "" {
		return "nil"
	}
	return *n.String_
}

// Scan implements the Scanner interface.
func (s *NullString) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*s.String_ = string(v)
	case *[]byte:
		*s.String_ = string(*v)
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}
	return nil
}

// Value implements the driver Valuer interface.
func (s NullString) Value() (driver.Value, error) {
	if s.String_ == nil {
		return nil, nil
	}
	if *s.String_ == "" {
		return nil, nil
	}
	return *s.String_, nil
}
