package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeInsert(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	writeExec(wln)
}

func writeExec(wln func(string)) {
	wln("query := buf.String()")
	wln("log.Debug(query)")
	wln("log.Debug(args...)")

	wln("res, err := db.Exec(query, args...)")
	wln("if err != nil {")
	wln("log.Error(query)")
	wln("log.Error(args...)")
	wln("log.Error(err)")
	wln("return 0, err")
	wln("}")
	wln("return res.LastInsertId()")
}
