package light

import (
	"database/sql"
	"database/sql/driver"
	"reflect"
)

type ValueScanner interface {
	driver.Valuer
	sql.Scanner
}

func String(v *string) ValueScanner {
	return &istring{S: v}
}

// NullString represents a string that may be null.
// NullString implements the Scanner interface so
// it can be used as a scan destination:
//
//  var plain string
//  s := &String{S:&s}
//  err := db.QueryRow("SELECT name FROM foo WHERE id=?", id).Scan(&s)
//  ...
//  use plain if database return null, plain is blank
type istring struct {
	S *string
}

// Scan implements the Scanner interface.
func (s *istring) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case []byte:
		*s.S = string(v)
	case *[]byte:
		*s.S = string(*v)
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	return nil
}

// Value implements the driver Valuer interface.
func (s istring) Value() (driver.Value, error) {
	if s.S == nil {
		return nil, nil
	}
	if *s.S == "" {
		return nil, nil
	}
	return *s.S, nil
}

func Uint8(v *uint8) ValueScanner {
	return &iuint8{S: v}
}

type iuint8 struct {
	S *uint8
}

// Scan implements the Scanner interface.
func (s *iuint8) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	switch v := value.(type) {
	case int64:
		*s.S = uint8(v)
	case *int64:
		*s.S = uint8(*v)
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	return nil
}

func (s iuint8) Value() (driver.Value, error) {
	if s.S == nil {
		return nil, nil
	}
	if *s.S == 0 {
		return nil, nil
	}
	return int64(*s.S), nil
}
