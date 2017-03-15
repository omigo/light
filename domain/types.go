package domain

import (
	"fmt"
	"strings"
)

type MethodKind string

const (
	Insert MethodKind = "insert"
	Batch             = "batch"
	Update            = "update"
	Delete            = "delete"
	Get               = "get"
	Count             = "count"
	List              = "list"
	Page              = "page"
)

type Package struct {
	Source string
	Path   string
	Name   string

	Imports map[string]string

	Interfaces []*Interface
}

type Interface struct {
	Name string

	Methods []*Method
}

func (itf *Interface) ImplName() string {
	return itf.Name + "Impl"
}

type Method struct {
	Name string
	Doc  string

	Kind      MethodKind
	Fragments []*Fragment

	Params  []*VarType
	Results []*VarType

	Returnings []*VarType
}

func (m *Method) ParamsExpr() string {
	return varTypesExpr(m.Params)
}

func (m *Method) ResultsExpr() string {
	return varTypesExpr(m.Results)
}

func varTypesExpr(vts []*VarType) string {
	var elems []string
	for _, vt := range vts {
		elems = append(elems, vt.Expr())
	}
	return strings.Join(elems, ", ")
}

type VarType struct {
	// ms []*domain.Model
	Var string `json:"Var,omitempty"` //  ms
	Tag string `json:"Tag,omitempty"`

	Path    string `json:"Path,omitempty"`    //  github.com/arstd/light/example/domain
	Array   string `json:"Array,omitempty"`   //  []
	Slice   string `json:"Slice,omitempty"`   //  []
	Pointer string `json:"Pointer,omitempty"` //  *
	Pkg     string `json:"Pkg,omitempty"`     //  domain
	Name    string `json:"Name,omitempty"`    //  Model
	Alias   string `json:"Alias,omitempty"`   //  e.g. domain.State => string

	Key  string `json:"Key,omitempty"`
	Elem string `json:"Elem,omitempty"`

	Deep   bool       `json:"Deep,omitempty"` //  深入解析这个类型
	Fields []*VarType `json:"Fields,omitempty"`
}

func (vt *VarType) Expr() string {
	var pkg string
	if vt.Pkg != "" {
		pkg = vt.Pkg + "."
	}
	return fmt.Sprintf("%s %s%s%s%s%s", vt.Var, vt.Array, vt.Slice, vt.Pointer, pkg, vt.Name)
}

type Fragment struct {
	Cond    string     `json:"Cond,omitempty"`
	Stmt    string     `json:"Stmt,omitempty"`
	Prepare string     `json:"Prepare,omitempty"`
	Args    []*VarType `json:"Args,omitempty"`

	Range     *VarType `json:"Range,omitempty"`
	Index     *VarType `json:"Index,omitempty"`
	Iterator  *VarType `json:"Iterator,omitempty"`
	Seperator string   `json:"Seperator,omitempty"`

	Fragments []*Fragment `json:"Fragments,omitempty"`
}
