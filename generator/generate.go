package generator

import (
	"bytes"
	"os"
	"strings"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

func Generate(store *goparser.Store) {
	buf := bytes.NewBuffer(make([]byte, 0, 65536))
	for _, m := range store.Methods {
		p := sqlparser.NewParser(bytes.NewBufferString(m.Doc))
		stmt, err := p.Parse()
		if err != nil {
			panic(err)
		}

		writeSignature(buf, m)

		genCondition(stmt, m)
		log.JSONIndent(stmt)
		switch stmt.Type {
		case sqlparser.SELECT:
			writeSelect(buf, m, stmt)

		default:
			panic("unimplemented " + m.Doc)
		}
		buf.WriteByte('}')
	}

	buf.WriteTo(os.Stdout)
}

func writeSignature(buf *bytes.Buffer, m *goparser.Method) {
	w := buf.WriteString
	w("func (*")
	w(m.Store.Name)
	w("Store)")
	w(m.Name)
	w("(")
	w(m.Params.String())
	w(")(")
	w(m.Results.String())
	w("){\n")
}

func genCondition(stmt *sqlparser.Statement, m *goparser.Method) {
	for _, v := range stmt.Fragments {
		deepGenCondition(v, m)
	}
}

func deepGenCondition(f *sqlparser.Fragment, m *goparser.Method) {
	if f.Condition == "" {
		return
	}

	if len(f.Fragments) == 0 {
		var cs []string
		for _, name := range f.Variables {
			v := m.Params.VarByName(name)
			d := v.NotDefault(name)
			cs = append(cs, d)
		}
		f.Condition = strings.Join(cs, " && ")
		return
	}

	for _, v := range f.Fragments {
		deepGenCondition(v, m)
	}

	var cs []string
	for _, v := range f.Fragments {
		if v.Condition == "" {
			continue
		}
		cs = append(cs, v.Condition)
	}
	f.Condition = strings.Join(cs, " && ")
}
