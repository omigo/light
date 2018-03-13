package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeUpdate(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	wln("query := buf.String()")
	if m.Store.Log {
		wln("log.Debug(query)")
		wln("log.Debug(args...)")
	}

	wln("res, err := exec.Exec(query, args...)")
	wln("if err != nil {")
	if m.Store.Log {
		wln("log.Error(query)")
		wln("log.Error(args...)")
		wln("log.Error(err)")
	}
	wln("return 0, err")
	wln("}")
	wln("return res.RowsAffected()")
}
