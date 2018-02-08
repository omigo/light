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
		if v.Range != "" {
			wln(`if len(` + v.Range + `) > 0 {`)
			w(`fmt.Fprintf(&buf, "`)
			w(v.Statement)
			wln(` ", strings.Repeat(",?", len(` + v.Range + `))[1:])`)
			wln(`for _, v := range ` + v.Range + ` {
						args = append(args, v)
				}
			}`)
		} else if len(v.Replacers) > 0 {
			w(`fmt.Fprintf(&buf, "`)
			w(v.Statement)
			w(` "`)
			for _, name := range v.Replacers {
				w(",")
				w(name)
			}
			wln(")")
		} else {
			w(`buf.WriteString("`)
			w(v.Statement)
			wln(` ")`)
		}

		if len(v.Variables) > 0 {
			w("args = append(args")
			for _, name := range v.Variables {
				w(", ")
				x := m.Params.VarByName(name)
				w(x.Value(x.VName))
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
