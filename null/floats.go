package null

import (
	"database/sql/driver"
	"reflect"
	"strconv"
)

type NullFloat32 struct{ Float32 *float32 }
type NullFloat64 struct{ Float64 *float64 }

func (n *NullFloat32) String() string { return floatToString(n.Float32) }
func (n *NullFloat64) String() string { return floatToString(n.Float64) }

func (n *NullFloat32) Scan(value interface{}) error { return scanFloat(n.Float32, value) }
func (n *NullFloat64) Scan(value interface{}) error { return scanFloat(n.Float64, value) }

func (n NullFloat32) Value() (driver.Value, error) { return valueFloat(n.Float32) }
func (n NullFloat64) Value() (driver.Value, error) { return valueFloat(n.Float64) }

func floatToString(ptr interface{}) string {
	if ptr == nil {
		return "nil"
	}

	f64 := toFloat64(ptr)

	if f64 == 0 {
		return "nil"
	}

	return strconv.FormatFloat(f64, 'e', 2, 32)
}

func valueFloat(ptr interface{}) (driver.Value, error) {
	if ptr == nil {
		return nil, nil
	}

	f64 := toFloat64(ptr)
	if f64 == 0 {
		return nil, nil
	}
	return f64, nil
}

func toFloat64(ptr interface{}) (f64 float64) {
	switch v := ptr.(type) {
	case *float32:
		f64 = float64(*v)
	case *float64:
		f64 = float64(*v)
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}
	return
}

func scanFloat(ptr, value interface{}) error {
	if value == nil {
		return nil
	}

	var f64 int64
	switch v := value.(type) {
	case int64:
		f64 = v
	case *int64:
		f64 = *v
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	switch v := ptr.(type) {
	case *float32:
		*v = float32(f64)
	case *float64:
		*v = float64(f64)
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	return nil
}
