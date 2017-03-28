package light

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"

	"github.com/arstd/log"
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
	js, err := json.Marshal(b.a)
	if err != nil {
		log.Error(err)
	}
	return js, nil
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

	err := json.Unmarshal(js, b.a)
	if err != nil {
		log.Error(err)
	}
	return err
}
