package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeFragment(buf *bytes.Buffer, m *goparser.Method, v *sqlparser.Fragment) {
	w := buf.WriteString
	wln := func(s string) { buf.WriteString(s + "\n") }

	if v.Condition != "" {
		w("if ")
		w(v.Condition)
		wln(" {")
	}

	if v.Statement != "" {
		if len(v.Replacers) > 0 {
			w("fmt.Fprintf(&buf, `")
			w(v.Statement)
			w(" `")
			for _, name := range v.Replacers {
				w(",")
				w(name)
			}
			wln(")")
		} else {
			w("buf.WriteString(`")
			w(v.Statement)
			wln(" `)")
		}
		if len(v.Variables) > 0 {
			w("args = append(args")
			for _, name := range v.Variables {
				w(", ")
				w(m.Params.VarByName(name).Value(name))
			}
			wln(")")
		}
	} else {
		for _, x := range v.Fragments {
			writeFragment(buf, m, x)
		}
	}

	if v.Condition != "" {
		wln("}")
	}
}
