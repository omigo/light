package goparser

import (
	"go/types"
	"strings"
)

func ParamsLast(ps *Params) string { return ps.List[len(ps.List)-1].FullName() }

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
