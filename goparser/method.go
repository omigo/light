package goparser

import (
	"bytes"
	"strings"

	"github.com/arstd/light/sqlparser"
)

type MethodType string

const (
	MethodTypeDDL    = "ddl"
	MethodTypeInsert = "insert"
	MethodTypeBulky  = "bulky"
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

	Name      string // Insert
	Doc       string // insert into users ...
	Signature string // Insert(tx *sql.Tx, u *model.User) (int64, error)

	Statement *sqlparser.Statement
	Type      MethodType

	Params  *Params
	Results *Results
}

func NewMethod(itf *Interface, name, doc string) *Method {
	return &Method{
		Interface: itf,
		Name:      name,
		Doc:       doc,
	}
}

func (m *Method) SetSignature() {
	var buf bytes.Buffer
	buf.WriteString(m.Name)

	buf.WriteByte('(')
	for i, v := range m.Params.List {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(v.Define())
	}
	buf.WriteByte(')')

	buf.WriteByte('(')
	for i, v := range m.Results.List {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(v.Define())
	}
	buf.WriteByte(')')

	m.Signature = buf.String()
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
		if m.Params.List[len(m.Params.List)-1].Slice {
			m.Type = MethodTypeBulky
		} else {
			m.Type = MethodTypeInsert
		}

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
