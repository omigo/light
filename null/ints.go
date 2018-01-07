package null

import (
	"database/sql/driver"
	"reflect"
	"strconv"
)

type NullInt struct{ Int *int }
type NullInt8 struct{ Int8 *int8 }
type NullUint8 struct{ Uint8 *uint8 }
type NullInt16 struct{ Int16 *int16 }
type NullUint16 struct{ Uint16 *uint16 }
type NullInt32 struct{ Int32 *int32 }
type NullUint32 struct{ Uint32 *uint32 }
type NullInt64 struct{ Int64 *int64 }
type NullUint64 struct{ Uint64 *uint64 }

func (n *NullInt) String() string    { return toString(n.Int) }
func (n *NullInt8) String() string   { return toString(n.Int8) }
func (n *NullUint8) String() string  { return toString(n.Uint8) }
func (n *NullInt16) String() string  { return toString(n.Int16) }
func (n *NullUint16) String() string { return toString(n.Uint16) }
func (n *NullInt32) String() string  { return toString(n.Int32) }
func (n *NullUint32) String() string { return toString(n.Uint32) }
func (n *NullInt64) String() string  { return toString(n.Int64) }
func (n *NullUint64) String() string { return toString(n.Uint64) }

func (n *NullInt) Scan(value interface{}) error    { return scan(n.Int, value) }
func (n *NullInt8) Scan(value interface{}) error   { return scan(n.Int8, value) }
func (n *NullUint8) Scan(value interface{}) error  { return scan(n.Uint8, value) }
func (n *NullInt16) Scan(value interface{}) error  { return scan(n.Int16, value) }
func (n *NullUint16) Scan(value interface{}) error { return scan(n.Uint16, value) }
func (n *NullInt32) Scan(value interface{}) error  { return scan(n.Int32, value) }
func (n *NullUint32) Scan(value interface{}) error { return scan(n.Uint32, value) }
func (n *NullInt64) Scan(value interface{}) error  { return scan(n.Int64, value) }
func (n *NullUint64) Scan(value interface{}) error { return scan(n.Uint64, value) }

func (n NullInt) Value() (driver.Value, error)    { return value(n.Int) }
func (n NullInt8) Value() (driver.Value, error)   { return value(n.Int8) }
func (n NullUint8) Value() (driver.Value, error)  { return value(n.Uint8) }
func (n NullInt16) Value() (driver.Value, error)  { return value(n.Int16) }
func (n NullUint16) Value() (driver.Value, error) { return value(n.Uint16) }
func (n NullInt32) Value() (driver.Value, error)  { return value(n.Int32) }
func (n NullUint32) Value() (driver.Value, error) { return value(n.Uint32) }
func (n NullInt64) Value() (driver.Value, error)  { return value(n.Int64) }
func (n NullUint64) Value() (driver.Value, error) { return value(n.Uint64) }

func toString(ptr interface{}) string {
	if ptr == nil {
		return "nil"
	}

	i64 := toInt64(ptr)
	if i64 == 0 {
		return "nil"
	}

	return strconv.FormatInt(i64, 10)
}

func value(ptr interface{}) (driver.Value, error) {
	if ptr == nil {
		return nil, nil
	}

	i64 := toInt64(ptr)
	if i64 == 0 {
		return nil, nil
	}
	return i64, nil
}

func toInt64(ptr interface{}) (i64 int64) {
	switch v := ptr.(type) {
	case *int8:
		i64 = int64(*v)
	case *uint8:
		i64 = int64(*v)
	// case *byte:
	// 	i64 = int64(*v)
	case *int16:
		i64 = int64(*v)
	case *uint16:
		i64 = int64(*v)
	case *int32:
		i64 = int64(*v)
	case *uint32:
		i64 = int64(*v)
	case *int:
		i64 = int64(*v)
	// case *rune:
	// 	i64 = int64(*v)
	case *int64:
		i64 = *v
	case *uint64:
		i64 = int64(*v)

	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}
	return
}

func scan(ptr, value interface{}) error {
	if value == nil {
		return nil
	}

	var i64 int64
	switch v := value.(type) {
	case int64:
		i64 = v
	case *int64:
		i64 = *v
	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	switch v := ptr.(type) {
	case *int8:
		*v = int8(i64)
	case *uint8:
		*v = uint8(i64)
	// case *byte:
	// 	*v = byte(i64)
	case *int16:
		*v = int16(i64)
	case *uint16:
		*v = uint16(i64)
	case *int32:
		*v = int32(i64)
	case *uint32:
		*v = uint32(i64)
	case *int:
		*v = int(i64)
	// case *rune:
	// 	*v = rune(i64)
	case *int64:
		*v = i64
	case *uint64:
		*v = uint64(i64)

	default:
		panic("unsupported type " + reflect.TypeOf(v).String())
	}

	return nil
}
