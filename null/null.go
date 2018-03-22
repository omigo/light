package null

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

func String(v *string) ValueScanner   { return &NullString{String_: v} }
func Uint8(v *uint8) ValueScanner     { return &NullUint8{Uint8: v} }
func Byte(v *byte) ValueScanner       { return &NullUint8{Uint8: v} }
func Int8(v *int8) ValueScanner       { return &NullInt8{Int8: v} }
func Uint16(v *uint16) ValueScanner   { return &NullUint16{Uint16: v} }
func Int16(v *int16) ValueScanner     { return &NullInt16{Int16: v} }
func Uint32(v *uint32) ValueScanner   { return &NullUint32{Uint32: v} }
func Int32(v *int32) ValueScanner     { return &NullInt32{Int32: v} }
func Rune(v *rune) ValueScanner       { return &NullInt32{Int32: v} }
func Int(v *int) ValueScanner         { return &NullInt{Int: v} }
func Uint64(v *uint64) ValueScanner   { return &NullUint64{Uint64: v} }
func Int64(v *int64) ValueScanner     { return &NullInt64{Int64: v} }
func Float32(v *float32) ValueScanner { return &NullFloat32{Float32: v} }
func Float64(v *float64) ValueScanner { return &NullFloat64{Float64: v} }
func Time(v *time.Time) ValueScanner  { return &NullTime{Time: v} }
func Bool(v *bool) ValueScanner       { return &NullBool{Bool: v} }

type ValueScanner interface {
	driver.Valuer
	sql.Scanner
}
