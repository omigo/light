package goparser

import (
	"go/types"
	"reflect"
	"strings"

	"github.com/arstd/log"
)

type Store struct {
	Package string            // store
	Imports map[string]string // fmt database/sql
	Name    string            // User

	Methods []*Method
}

func (s *Store) MethodByName(name string) *Method {
	for _, a := range s.Methods {
		if a.Name == name {
			return a
		}
	}
	return nil
}

type Method struct {
	Store *Store `json:"-"`

	Name string // Insert
	Doc  string // insert into users ...

	Params  *Tuple
	Results *Tuple
}

type Tuple struct {
	Store *Store `json:"-"`
	*types.Tuple
}

func (t *Tuple) String() string {
	var ss []string
	for i := 0; i < t.Len(); i++ {
		ss = append(ss, t.At(i).String())
	}
	return strings.Join(ss, ", ")
}

func (t *Tuple) At(i int) *Var {
	return &Var{t.Store, t.Tuple.At(i), ""}
}

func (t *Tuple) VarByName(name string) *Var {
	if name == "" {
		panic("name must not blank")
	}
	parts := strings.Split(name, ".")
	var v *Var
	for i := 0; i < t.Len(); i++ {
		x := t.At(i)
		if x.Name() == parts[0] {
			v = x
			break
		}
	}
	if v == nil {
		panic(name + " not exist")
	}

	switch len(parts) {
	case 1:
		return v

	case 2:
		s := underlying(v.Type())
		for i := 0; i < s.NumFields(); i++ {
			x := s.Field(i)
			if x.Name() == parts[1] {
				return &Var{t.Store, x, s.Tag(i)}
			}
		}

	default:
	}
	panic(name + " to long")
}

type Var struct {
	Store *Store `json:"-"`
	*types.Var
	Tag string
}

func (v *Var) NotDefault(name string) string {
	switch u := v.Type().(type) {
	case *types.Named:
		if u.String() == "time.Time" {
			return "!" + name + ".IsZero()"
		}
		return `name != ""`

	case *types.Basic:
		switch u.Kind() {
		case types.String:
			return `name != ""`
		case types.Int, types.Int8, types.Int16, types.Int32, types.Int64,
			types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64,
			types.Float32, types.Float64:
			return "name != 0"
		case types.Bool:
			return "name"
		case types.Uintptr, types.UnsafePointer:
			return "name != nil"
		default:
			panic(reflect.TypeOf(u))
		}

	case *types.Pointer:
		return "name != nil"

	case *types.Struct:
		return "name != nil"

	default:
		panic(" unimplement " + reflect.TypeOf(u).String() + u.String())
	}
}

func (v *Var) Value(name string) string {
	if v.Wrap() == "" {
		return name
	}
	return v.Wrap() + "(&" + name + ")"
}

func (v *Var) Scan(name string) string {
	panic("unimplemented")
}

func (v *Var) Nullable() bool {
	log.Error("umimplemented")
	return true
}

func (v *Var) Pointer() bool {
	_, ok := v.Type().(*types.Pointer)
	return ok
}

func (v *Var) Wrap() string {
	switch u := v.Type().(type) {
	case *types.Pointer, *types.Named:
		return ""

	case *types.Basic:
		if v.Nullable() {
			switch u.Kind() {
			case types.Uint8:
				return "light.Uint8"
			case types.String:
				return "light.String"
			default:
				log.Warn(u)
				return ""
			}
		} else {
			return ""
		}

	default:
		panic(reflect.TypeOf(u))
	}
}

func underlying(t types.Type) *types.Struct {
	switch u := t.(type) {
	case *types.Named:
		return underlying(u.Underlying())

	case *types.Pointer:
		return underlying(u.Elem())

	case *types.Struct:
		return u

	default:
		panic(" unimplement " + reflect.TypeOf(u).String())
	}
}

func (v *Var) String() string {
	typ := typeString(v.Store, v.Type())
	name := v.Name()
	return name + " " + typ
}

func typeString(store *Store, t types.Type) string {
	switch u := t.(type) {
	case *types.Named:
		if obj := u.Obj(); obj != nil {
			if pkg := obj.Pkg(); pkg != nil {
				store.Imports[pkg.Path()] = ""
				return shortPkg(pkg.Path()) + "." + obj.Name()
			}
			return obj.Name()
		}
		return typeString(store, u.Underlying())

	case *types.Basic:
		return u.String()

	case *types.Pointer:
		return "*" + typeString(store, u.Elem())

	case *types.Struct:
		return u.String()

	case *types.Slice:
		return "[]" + typeString(store, u.Elem())

	default:
		panic(" unimplement " + reflect.TypeOf(u).String())
	}
}

func shortPkg(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}
