package generator

import (
	"bytes"
	"strings"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func Generate(store *goparser.Store) []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 65536))
	for _, m := range store.Methods {
		p := sqlparser.NewParser(bytes.NewBufferString(m.Doc))
		stmt, err := p.Parse()
		if err != nil {
			panic(err)
		}

		buf.WriteString(m.Signature())

		genCondition(stmt, m)
		// log.JSONIndent(stmt)

		switch stmt.Type {
		case sqlparser.SELECT:
			writeSelect(buf, m, stmt)

		case sqlparser.INSERT:
			writeInsert(buf, m, stmt)

		case sqlparser.UPDATE:
			writeUpdate(buf, m, stmt)

		case sqlparser.DELETE:
			writeDelete(buf, m, stmt)

		case sqlparser.CREATE:
			writeCreate(buf, m, stmt)

		default:
			panic("unimplemented " + m.Doc)
		}
		buf.WriteString("}\n\n")
	}

	header := writeHeader(store)
	buf.WriteTo(header)
	return header.Bytes()
}

func genCondition(stmt *sqlparser.Statement, m *goparser.Method) {
	for _, v := range stmt.Fragments {
		deepGenCondition(v, m)
	}
}

func deepGenCondition(f *sqlparser.Fragment, m *goparser.Method) {
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
