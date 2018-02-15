package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writeCount(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	w := buf.WriteString
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("query := buf.String()")
	if m.Store.Log {
		wln("log.Debug(query)")
		wln("log.Debug(args...)")
	}

	w("var total ")
	wln(m.Results.Result().TypeName())
	wln(`err := db.QueryRow(query, args...).Scan(` + m.Results.Result().Wrap(true) + `(&total))
		if err != nil {`)
	if m.Store.Log {
		wln(`log.Error(query)
			log.Error(args...)
			log.Error(err)`)
	}
	wln(`return total, err
		}`)
	if m.Store.Log {
		wln(`log.Debug(total)`)
	}

	wln("return total, nil")
}
