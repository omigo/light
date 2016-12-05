package main

import (
	"go/ast"
	"go/token"
)

type MethodKind string

const (
	KindInsert = "insert"
	KindUpdate = "update"
	KindDelete = "delete"
	KindGet    = "get"
	KindCount  = "count"
	KindList   = "list"
	KindPage   = "page"
)

type Analyzer struct {
	Filename string

	Package string
	Imports []string

	Interfaces []*Interface

	// used internal
	fset *token.FileSet
	f    *ast.File
}

type Interface struct {
	Analyzer *Analyzer `json:"-"`

	Name    string
	Comment string

	Methods []*Method
}

type Method struct {
	Analyzer  *Analyzer  `json:"-"`
	Interface *Interface `json:"-"`

	Name    string
	Comment string

	Params  []*VarAndType
	Results []*VarAndType

	Kind MethodKind
}

type VarAndType struct {
	Analyzer  *Analyzer  `json:"-"`
	Interface *Interface `json:"-"`
	Method    *Method    `json:"-"`

	Var  string
	Type string
}

func (vt *VarAndType) IsPrimitive() bool {
	// TODO
	return true
}

func (vt *VarAndType) IsStruct() bool {
	// TODO
	return true
}

func (vt *VarAndType) IsArray() bool {
	// TODO
	return true
}
