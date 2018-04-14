package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textUpdate = textDelete

var tplUpdate = template.Must(template.New("textUpdate").Parse(textUpdate))

func writeUpdate(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	log.Errorn(tplUpdate.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
