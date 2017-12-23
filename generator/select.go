package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeSelect(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}
}

func writeFragment(buf *bytes.Buffer, m *goparser.Method, v *sqlparser.Fragment) {
	w := buf.WriteString
	wln := func(s string) { buf.WriteString(s + "\n") }

	if v.Condition != "" {
		w("if ")
		w(v.Condition)
		wln(" {")
	}

	if v.Statement != "" {
		w("buf.WriteString(`")
		w(v.Statement)
		wln("`)")
		w("args = append(args")
		for _, name := range v.Variables {
			w(", ")
			w(m.Params.VarByName(name).Value(name))
		}
		wln(")")
	} else {
		for _, x := range v.Fragments {
			writeFragment(buf, m, x)
		}
	}

	if v.Condition != "" {
		wln("}")
	}
}
