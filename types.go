package main

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
	Store *Store

	Name string // Insert
	Doc  string // insert into users ...

	Params  *Tuple
	Results *Tuple
}

type Tuple struct {
	Store *Store
	*types.Tuple
}

func (t *Tuple) List() []*Var {
	list := make([]*Var, t.Len())
	for i := 0; i < t.Len(); i++ {
		list[i] = &Var{t.Store, t.At(i), ""}
	}
	return list
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
			v = &Var{t.Store, x, ""}
			break
		}
	}
	if v == nil {
		panic(name + " not exist")
	}

	if len(parts) == 2 {
		panic(name + " to long")
	}

	return v
}

func (t *Tuple) LightByName(name string) *Var {
	if name == "" {
		panic("name must not blank")
	}
	parts := strings.Split(name, ".")
	var v *Var
	for i := 0; i < t.Len(); i++ {
		x := t.At(i)
		if x.Name() == parts[0] {
			v = &Var{t.Store, x, ""}
			break
		}
	}
	if v == nil {
		panic(name + " not exist")
	}

	if len(parts) != 2 {
		panic(name + " not short")
	}

	s := underlying(v.Type())
	for i := 0; i < s.NumFields(); i++ {
		x := s.Field(i)
		if x.Name() == parts[1] {
			return &Var{t.Store, x, s.Tag(i)}
		}
	}

	panic(name + " not exist")
}

type Var struct {
	Store *Store
	*types.Var
	Tag string
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
	case *types.Pointer:
		return ""
	case *types.Basic:
		if v.Nullable() {
			switch u.Kind() {
			case types.Uint8:
				return "light.Uint8"
			case types.String:
				return "light.String"
			default:
				panic(u.Kind())
			}
		} else {
			return ""
		}
	default:
		panic(u.String())
	}
}

func underlying(t types.Type) *types.Struct {
	switch u := t.(type) {
	case *types.Named:
		return underlying(u.Underlying())

	//
	// case *types.Basic:
	//
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
	if name == "" {
		name = strings.ToLower(typ)[:1]
	}
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

	default:
		panic(" unimplement " + reflect.TypeOf(u).String())
	}
}

func shortPkg(path string) string {
	return path[strings.LastIndex(path, "/")+1:]
}
