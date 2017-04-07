package light

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

func JSON(a interface{}) interface {
	driver.Valuer
	sql.Scanner
} {
	return &JSONWapper{a}
}

type JSONWapper struct {
	a interface{}
}

func (b JSONWapper) Value() (driver.Value, error) {
	if b.a == nil {
		return []byte(""), nil
	}
	return json.Marshal(b.a)
}

func (b *JSONWapper) Scan(src interface{}) error {
	var js []byte
	switch s := src.(type) {
	case string:
		js = []byte(s)
	case []byte:
		js = s
	case nil:
		b = nil
		return nil
	}

	return json.Unmarshal(js, b.a)
}
