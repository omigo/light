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

func MethodSignature(m *Method) string {
	name := m.Store.Name
	if name[0] == 'I' {
		name = name[1:]
	}
	return m.Name + "(" + m.Params.String() + ")(" + m.Results.String() + ")"
}

func MethodTx(m *Method) string {
	for i := 0; i < m.Params.Len(); i++ {
		v := m.Params.List[i]
		typ := typeString(v.Store, v.Var.Type())
		if typ == "*sql.Tx" {
			return v.Var.Name()
		}
	}
	return ""
}

func HasVariable(m *Method) bool {
	for _, f := range m.Statement.Fragments {
		if len(f.Variables) > 0 || f.Range != "" {
			return true
		}
	}
	return false
}

type Method struct {
	Store *Store `json:"-"`

	Name string // Insert
	Doc  string // insert into users ...

	Statement *sqlparser.Statement
	Type      MethodType

	Params  *Params
	Results *Results
}

func NewMethod(store *Store, name, doc string) *Method {
	m := &Method{Store: store, Name: name, Doc: doc}
	return m
}

func (m *Method) SetType() {
	switch m.Statement.Type {
	case sqlparser.SELECT:
		if m.Results.Len() == 3 {
			m.Type = MethodTypePage
		} else if m.Results.List[0].IsSlice() {
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
