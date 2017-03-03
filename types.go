package main

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/arstd/log"
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

func (a *Analyzer) GetPath(pkg string) string {
	for _, imp := range a.Imports {
		if imp[0] != '"' && strings.HasPrefix(imp, pkg+" ") {
			return imp[len(pkg)+2 : len(imp)-1]
		} else if strings.HasSuffix(imp, "/"+pkg+`"`) {
			return imp[1 : len(imp)-1]
		}
	}
	log.Panicf("import path not found for %s", pkg)
	return "" // unreachable code
}

type Interface struct {
	a *Analyzer

	Name    string
	Comment string

	Methods []*Method
}

type Method struct {
	a *Analyzer
	i *Interface

	Name    string
	Comment string

	Params  []*VarAndType
	Results []*VarAndType

	Kind MethodKind
}

type VarAndType struct {
	a *Analyzer
	i *Interface
	m *Method

	Var  string
	Type string

	Slice  string
	Star   string
	Pkg    string
	Alias  string
	Path   string
	Fields []*VarAndType
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
