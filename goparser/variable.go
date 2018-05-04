package goparser

import (
	"go/types"
	"log"
	"reflect"
	"strings"
)

func VariableTypeName(v *Variable) string     { return v.TypeName() }
func VariableWrap(v *Variable) string         { return v.Wrap(true) }
func VariableElemTypeName(v *Variable) string { return v.ElemTypeName() }
func VariableVarByTagScan(s *Variable, name string) string {
	v := s.VarByTag(name)
	return v.Scan("xu." + v.VName)
}

type Variable struct {
	VName string
	Store *Store `json:"-"`
	Var   *types.Var
	Tag   string
}

func (v *Variable) VarByTag(field string) *Variable {
	field = strings.Trim(field, "`")
	s := underlying(v.Var.Type())
	for i := 0; i < s.NumFields(); i++ {
		tag := s.Tag(i)
		t := getTag(tag, "light")
		if t != "" {
			tt := strings.Split(t, ",")
			if tt[0] != "" {
				if strings.HasPrefix(t, tt[0]) {
					return &Variable{
						VName: s.Field(i).Name(),
						Store: v.Store,
						Var:   s.Field(i),
						Tag:   t,
					}
				}
			}
		}
	}

	out := upperCamelCase(field)
	for i := 0; i < s.NumFields(); i++ {
		x := s.Field(i)
		if strings.EqualFold(out, x.Name()) {
			t := getTag(s.Tag(i), "light")
			return &Variable{
				VName: s.Field(i).Name(),
				Store: v.Store,
				Var:   s.Field(i),
				Tag:   t,
			}
		}
	}
	panic(field + " not found")
}

func lowerCamelCase(field string) (out string) {
	return camelCase(field, false)
}

func upperCamelCase(field string) (out string) {
	return camelCase(field, true)
}

func camelCase(name string, first bool) (out string) {
	for _, v := range name {
		if first {
			out += strings.ToUpper(string(v))
			first = false
		} else if v == '_' {
			first = true
		} else {
			out += string(v)
			first = false
		}
	}
	return out
}

func getTag(tag, key string) string {
	idx := strings.Index(tag, key+`:"`)
	if idx == -1 {
		return ""
	}
	tag = tag[idx+len(key)+2:]
	idx = strings.Index(tag, `"`)
	if idx == -1 {
		panic(tag)
	}
	return tag[:idx]
}

func (v *Variable) IsBasic() bool {
	_, ok := v.Var.Type().(*types.Basic)
	return ok
}

func (v *Variable) NotDefault(name string) string {
	switch u := v.Var.Type().(type) {
	case *types.Named:
		if u.String() == "time.Time" {
			return "!" + name + ".IsZero()"
		}

		t, ok := u.Underlying().(*types.Basic)
		if !ok {
			log.Fatalf("%#v", u)
		}
		bi := t.Info()
		switch {
		case bi&types.IsString == types.IsString:
			return name + ` != ""`
		case bi&types.IsInteger == types.IsInteger:
			return name + ` != 0`
		case bi&types.IsFloat == types.IsFloat:
			return name + ` != 0`
		case bi&types.IsBoolean == types.IsBoolean:
			return name
		default:
			panic(t)
		}

	case *types.Basic:
		bi := u.Info()
		switch {
		case bi&types.IsString == types.IsString:
			return name + ` != ""`
		case bi&types.IsInteger == types.IsInteger:
			return name + ` != 0`
		case bi&types.IsFloat == types.IsFloat:
			return name + ` != 0`
		case bi&types.IsBoolean == types.IsBoolean:
			return name
		default:
			panic(u)
		}

	case *types.Pointer:
		return name + " != nil"

	case *types.Struct:
		return name + " != nil"

	case *types.Slice:
		return "len(" + name + ") != 0"

	default:
		panic(" unimplement " + reflect.TypeOf(u).String() + u.String())
	}
}

func (v *Variable) Value(name string) string {
	if v.Wrap() == "" {
		return name
	}
	return v.Wrap() + "(&" + name + ")"
}

func (v *Variable) Scan(name string) string {
	s := v.Value(name)
	if strings.HasPrefix(s, "null.") {
		return s
	}
	return "&" + s
}

func (v *Variable) Nullable() bool {
	for i, v := range strings.Split(v.Tag, ",") {
		if i == 0 {
			continue
		}
		if v == "nullable" {
			return true
		}
	}
	return false
}

func (v *Variable) IsSlice() bool {
	_, ok := v.Var.Type().(*types.Slice)
	return ok
}

func (v *Variable) Wrap(force ...bool) string {
	switch u := v.Var.Type().(type) {
	case *types.Pointer, *types.Named, *types.Slice, *types.Array:
		return ""

	case *types.Basic:
		if v.Nullable() || (len(force) > 0 && force[0]) {
			name := u.Name()
			return "null." + strings.ToUpper(name[:1]) + name[1:]
		}
		return ""

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

	case *types.Slice:
		return underlying(u.Elem())

	case *types.Struct:
		return u

	default:
		return nil
	}
}

func (v *Variable) String() string {
	typ := typeString(v.Store, v.Var.Type())
	name := v.Var.Name()
	return name + " " + typ
}

func (v *Variable) TypeName() string {
	return strings.TrimLeft(v.String(), " *")
}

func (v *Variable) ElemTypeName() string {
	return strings.TrimLeft(v.String(), " []*")
}

func typeString(store *Store, t types.Type) string {
	switch u := t.(type) {
	case *types.Named:
		if obj := u.Obj(); obj != nil {
			if pkg := obj.Pkg(); pkg != nil {
				path := pkg.Path()
				if path != "" && path[0] != '/' {
					store.Imports[pkg.Path()] = ""
					return shortPkg(pkg.Path()) + "." + obj.Name()
				}
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
