package goparser

import (
	"strings"

	"github.com/arstd/light/sqlparser"
)

type MethodType string

const (
	MethodTypeDDL    = "ddl"
	MethodTypeInsert = "insert"
	MethodTypeUpdate = "update"
	MethodTypeDelete = "delete"
	MethodTypeGet    = "get"
	MethodTypeList   = "list"
	MethodTypePage   = "page"
	MethodTypeAgg    = "agg"
)

type Method struct {
	Store *Store `json:"-"`

	Name string // Insert
	Doc  string // insert into users ...

	Statement *sqlparser.Statement
	Type      MethodType

	Params  *Params
	Results *Results

	ResultTypeName     func() string
	ResultTypeWrap     func() string
	ResultElemTypeName func() string
	ResultVarByTagScan func(name string) string
	ParamsVarByName    func(string) *Variable
	Signature          func() string
	Tx                 func() string
}

func NewMethod(store *Store, name, doc string) *Method {
	m := &Method{Store: store, Name: name, Doc: doc}
	m.ResultTypeName = m.resultTypeName
	m.ResultTypeWrap = m.resultTypeWrap
	m.ResultElemTypeName = m.resultElemTypeName
	m.ResultVarByTagScan = m.resultVarByTagScan
	m.Signature = m.signature
	m.Tx = m.tx
	return m
}

func (m *Method) SetType() {
	switch m.Statement.Type {
	case sqlparser.SELECT:
		if m.Results.Len() == 3 {
			m.Type = MethodTypePage
		} else if m.Results.Results[0].IsSlice() {
			m.Type = MethodTypeList
		} else if m.Results.Result.IsBasic() {
			m.Type = MethodTypeAgg
		} else {
			m.Type = MethodTypeGet
		}

	case sqlparser.INSERT, sqlparser.REPLACE:
		m.Type = MethodTypeInsert

	case sqlparser.UPDATE:
		m.Type = MethodTypeUpdate

	case sqlparser.DELETE:
		m.Type = MethodTypeDelete

	default:
		m.Type = MethodTypeDDL
	}
}

func (m *Method) GenCondition() {
	for _, f := range m.Statement.Fragments {
		deepGenCondition(f, m)
	}
}

func deepGenCondition(f *sqlparser.Fragment, m *Method) {
	if len(f.Fragments) == 0 {
		if f.Condition == "-" {
			var cs []string
			for _, name := range f.Variables {
				v := m.Params.VarByName(name)
				d := v.NotDefault(v.VName)
				cs = append(cs, "("+d+")")
			}
			f.Condition = strings.Join(cs, " && ")
		}
		return
	}

	for _, v := range f.Fragments {
		deepGenCondition(v, m)
	}

	if f.Condition != "-" {
		return
	}

	var cs []string
	for _, v := range f.Fragments {
		if v.Condition == "" {
			continue
		}
		cs = append(cs, "("+v.Condition+")")
	}
	f.Condition = strings.Join(cs, " || ")
}

func (m *Method) signature() string {
	name := m.Store.Name
	if name[0] == 'I' {
		name = name[1:]
	}
	return m.Name + "(" + m.Params.String() + ")(" + m.Results.String() + ")"
}

func (m *Method) resultTypeName() string     { return m.Results.Result.TypeName() }
func (m *Method) resultTypeWrap() string     { return m.Results.Result.Wrap(true) }
func (m *Method) resultElemTypeName() string { return m.Results.Result.ElemTypeName() }
func (m *Method) resultVarByTagScan(name string) string {
	s := m.Results.Result
	v := s.VarByTag(name)
	return v.Scan("xu." + v.VName)
}

func (m *Method) tx() string {
	for i := 0; i < m.Params.Len(); i++ {
		v := m.Params.At(i)
		typ := typeString(v.Store, v.Var.Type())
		if typ == "*sql.Tx" {
			return v.Var.Name()
		}
	}
	return ""
}
