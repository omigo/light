package goparser

import (
	"go/types"
	"strings"
)

type Results struct {
	Store   *Store `json:"-"`
	Tuple   *types.Tuple
	Results []*Variable
	Result  *Variable
}

func NewResults(store *Store, tuple *types.Tuple) *Results {
	rs := &Results{
		Store:   store,
		Tuple:   tuple,
		Results: make([]*Variable, tuple.Len()),
	}

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		rs.Results[i] = &Variable{
			Store: store,
			VName: v.Name(),
			// Tag   :
			Var: v,
		}
	}

	switch tuple.Len() {
	case 1:
		// ddl
	case 2:
		rs.Result = rs.Results[0]
	case 3:
		rs.Result = rs.Results[1]
	default:
		panic(rs.Len())
	}
	return rs
}

func (rs *Results) String() string {
	var ss []string
	for _, r := range rs.Results {
		ss = append(ss, r.String())
	}
	return strings.Join(ss, ", ")
}

func (rs *Results) Len() int {
	return rs.Tuple.Len()
}

func (rs *Results) VarByName(name string) *Variable {
	name = strings.Trim(name, "`")
	if name == "" {
		panic("name must not blank")
	}

	var v *Variable
	parts := strings.Split(name, ".")

	parts0 := lowerCamelCase(parts[0])
	// 从参数列表中查找
	for i := 0; i < rs.Len(); i++ {
		x := rs.Results[i]
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
					return &Variable{VName: name, Store: rs.Store, Var: x, Tag: z}
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
	for _, r := range rs.Results {
		s := underlying(r.Var.Type())
		if s != nil {
			for j := 0; j < s.NumFields(); j++ {
				x := s.Field(j)
				if x.Name() == out {
					z := getTag(s.Tag(j), "light")
					return &Variable{
						VName: r.Var.Name() + "." + x.Name(),
						Store: rs.Store,
						Var:   x,
						Tag:   z,
					}
				}
			}
		}
	}
	panic("variable " + name + " not exist")
}
