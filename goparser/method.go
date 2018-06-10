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

func MethodTx(m *Method) string {
	for _, v := range m.Params.List {
		if v.Tx {
			return v.Name
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
	Interface *Interface `json:"-"`

	Name string // Insert
	Doc  string // insert into users ...
	Expr string // Insert(tx *sql.Tx, u *model.User) (int64, error)

	Statement *sqlparser.Statement
	Type      MethodType

	Params  *Params
	Results *Results
}

func NewMethod(itf *Interface, name, doc string, expr string) *Method {
	return &Method{
		Interface: itf,
		Name:      name,
		Doc:       doc,
		Expr:      expr,
	}
}

func (m *Method) SetType() {
	switch m.Statement.Type {
	case sqlparser.SELECT:
		switch {
		case len(m.Results.List) == 3:
			m.Type = MethodTypePage
		case m.Results.Result.Slice:
			m.Type = MethodTypeList
		case !m.Results.Result.Array && !m.Results.Result.Slice && !m.Results.Result.Struct:
			m.Type = MethodTypeAgg
		default:
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
				v := m.Params.Names[name]
				if v == nil {
					panic("method `" + m.Name + "` variable `" + name + "` not found")
				}
				d := v.NotDefault()
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
