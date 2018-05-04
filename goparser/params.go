package goparser

import (
	"go/types"
	"strings"
)

type Params struct {
	Store  *Store `json:"-"`
	Tuple  *types.Tuple
	Params []*Variable
}

func NewParams(store *Store, tuple *types.Tuple) *Params {
	return &Params{Store: store, Tuple: tuple}
}

func (t *Params) String() string {
	var ss []string
	for i := 0; i < t.Tuple.Len(); i++ {
		ss = append(ss, t.At(i).String())
	}
	return strings.Join(ss, ", ")
}

func (t *Params) Len() int {
	return t.Tuple.Len()
}

func (t *Params) At(i int) *Variable {
	x := t.Tuple.At(i)
	return &Variable{VName: x.Name(), Store: t.Store, Var: t.Tuple.At(i)}
}

func (t *Params) VarByName(name string) *Variable {
	name = strings.Trim(name, "`")
	if name == "" {
		panic("name must not blank")
	}

	var v *Variable
	parts := strings.Split(name, ".")

	parts0 := lowerCamelCase(parts[0])
	// 从参数列表中查找
	for i := 0; i < t.Tuple.Len(); i++ {
		x := t.At(i)
		if x.Var.Name() == parts0 {
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
			s := underlying(v.Var.Type())
			for i := 0; i < s.NumFields(); i++ {
				x := s.Field(i)
				if x.Name() == parts[1] {
					z := getTag(s.Tag(i), "light")
					return &Variable{VName: name, Store: t.Store, Var: x, Tag: z}
				}
			}
			panic("variable " + name + " not exist")

		default:
			panic("variable " + name + " to long")
		}
	}

	// 从结构体参数中查找
	if len(parts) > 1 {
		panic("variable " + name + " not exist")
	}

	out := upperCamelCase(name)
	for i := 0; i < t.Len(); i++ {
		s := underlying(t.At(i).Var.Type())
		if s != nil {
			for j := 0; j < s.NumFields(); j++ {
				x := s.Field(j)
				if x.Name() == out {
					z := getTag(s.Tag(j), "light")
					return &Variable{
						VName: t.At(i).Var.Name() + "." + x.Name(),
						Store: t.Store,
						Var:   x,
						Tag:   z,
					}
				}
			}
		}
	}
	panic("variable " + name + " not exist")
}
