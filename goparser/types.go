package goparser

import (
	"go/types"
	"reflect"
	"strings"
)

type Store struct {
	Source string
	Log    bool

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

func (m *Method) Signature() string {
	return "func (*" + m.Store.Name + "Store)" + m.Name +
		"(" + m.Params.String() + ")(" + m.Results.String() + "){\n"
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
	x := t.Tuple.At(i)
	return &Var{VName: x.Name(), Store: t.Store, Var: t.Tuple.At(i)}
}

func (t *Tuple) VarByName(name string) *Var {
	name = strings.Trim(name, "`")
	if name == "" {
		panic("name must not blank")
	}
	parts := strings.Split(name, ".")
	var v *Var

	// 从参数列表中查找
	for i := 0; i < t.Len(); i++ {
		x := t.At(i)
		if x.Name() == parts[0] {
			v = x
			break
		}
	}
	// 如果找到了
	if v != nil {
		switch len(parts) {
		case 1:
			return v

		case 2:
			s := underlying(v.Type())
			for i := 0; i < s.NumFields(); i++ {
				x := s.Field(i)
				if x.Name() == parts[1] {
					z := getTag(s.Tag(i), "light")
					return &Var{VName: name, Store: t.Store, Var: x, Tag: z}
				}
			}
			panic("variable " + name + " not exist")

		default:
			panic("variable " + name + " to long")
		}
	}

	// 从结构体参数中查找
	if len(parts) > 1 {
		panic("variable " + parts[0] + " not exist")
	}

	out := capitalize(name)
	for i := 0; i < t.Len(); i++ {
		s := underlying(t.At(i).Type())
		if s != nil {
			for j := 0; j < s.NumFields(); j++ {
				x := s.Field(j)
				if x.Name() == out {
					z := getTag(s.Tag(j), "light")
					return &Var{
						VName: t.At(i).Name() + "." + x.Name(),
						Store: t.Store,
						Var:   x,
						Tag:   z,
					}
				}
			}
		}
	}
	panic("variable " + parts[0] + " not exist")
}

func (t *Tuple) Result() *Var {
	switch t.Len() {
	case 1:
		panic("unimplemented")
	case 2:
		return t.At(0)
	case 3:
		return t.At(1)
	default:
		panic(t.Len())
	}
}

type Var struct {
	VName string
	Store *Store `json:"-"`
	Tag   string
	*types.Var
}

func (v *Var) VarByTag(field string) *Var {
	field = strings.Trim(field, "`")
	s := underlying(v.Type())
	for i := 0; i < s.NumFields(); i++ {
		tag := s.Tag(i)
		t := getTag(tag, "light")
		if t != "" {
			tt := strings.Split(t, ",")
			if tt[0] != "" {
				if strings.HasPrefix(t, tt[0]) {
					return &Var{
						VName: s.Field(i).Name(),
						Store: v.Store,
						Var:   s.Field(i),
						Tag:   t,
					}
				}
			}
		}
	}

	out := capitalize(field)
	for i := 0; i < s.NumFields(); i++ {
		x := s.Field(i)
		if strings.EqualFold(out, x.Name()) {
			t := getTag(s.Tag(i), "light")
			return &Var{
				VName: s.Field(i).Name(),
				Store: v.Store,
				Var:   s.Field(i),
				Tag:   t,
			}
		}
	}
	panic(field + " not found")
}

func capitalize(field string) (out string) {
	var first bool = true
	for _, v := range field {
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

func (v *Var) NotDefault(name string) string {
	switch u := v.Type().(type) {
	case *types.Named:
		if u.String() == "time.Time" {
			return "!" + name + ".IsZero()"
		}
		return name + ` != ""`

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
			panic(u.Name())
		}

	case *types.Pointer:
		return name + " != nil"

	case *types.Struct:
		return name + " != nil"

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
	s := v.Value(name)
	if strings.HasPrefix(s, "null.") {
		return s
	}
	return "&" + s
}

func (v *Var) Nullable() bool {
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

func (v *Var) IsSlice() bool {
	_, ok := v.Type().(*types.Slice)
	return ok
}

func (v *Var) Wrap() string {
	switch u := v.Type().(type) {
	case *types.Pointer, *types.Named:
		return ""

	case *types.Basic:
		if v.Nullable() {
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

func (v *Var) String() string {
	typ := typeString(v.Store, v.Type())
	name := v.Name()
	return name + " " + typ
}

func (v *Var) TypeName() string {
	return strings.TrimLeft(v.String(), " *")
}

func (v *Var) ElemTypeName() string {
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
