package goparser

import (
	"go/types"
	"strings"
)

//
// func ParamsVarByNameValue(ps *Params, name string) string {
// 	x := ps.VarByName(name)
// 	return x.Value(x.Name)
// }

type Params struct {
	Tuple *types.Tuple
	List  []*Variable

	Names map[string]*Variable
}

func NewParams(tuple *types.Tuple) *Params {
	ps := &Params{
		Tuple: tuple,
		List:  make([]*Variable, tuple.Len()),
		Names: make(map[string]*Variable),
	}

	for i := 0; i < tuple.Len(); i++ {
		v := tuple.At(i)
		ps.List[i] = NewVariable(v)
	}

	return ps
}

func (p *Params) Lookup(name string) *Variable {
	name = strings.Trim(name, "`")
	name = strings.TrimSpace(name)
	return p.Names[name]
}
