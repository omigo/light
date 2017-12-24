package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeDelete(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	writeExec(wln)
}
