package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
)

func writePage(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	w := buf.WriteString
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("query := buf.String()")
	wln("log.Debug(query)")
	wln("log.Debug(args...)")

	wln("rows, err := db.Query(query, args...)")
	wln("if err != nil {")
	wln("log.Error(query)")
	wln("log.Error(args...)")
	wln("log.Error(err)")
	wln("return 0, nil, err")
	wln("}")
	wln("defer rows.Close()")

	wln("var total int64")

	v := m.Results.Result()

	w("var data ")
	wln(v.TypeName())

	wln("for rows.Next() {")
	w("xu := new(")
	w(v.ElemTypeName())
	wln(")")
	wln("data = append(data, xu)")
	w("xdst := []interface{}{")
	for _, f := range stmt.Fields {
		s := m.Results.Result()
		v := s.VarByTag(f)
		name := "xu." + v.Name()
		w(v.Scan(name))
		w(",")
	}
	buf.Truncate(buf.Len() - 1)
	wln("}")

	wln("err = rows.Scan(xdst...)")
	wln("if err != nil {")
	wln("log.Error(query)")
	wln("log.Error(args...)")
	wln("log.Error(err)")
	wln("return 0, nil, err")
	wln("}")
	wln("}")
	wln("if err = rows.Err(); err != nil {")
	wln("log.Error(query)")
	wln("log.Error(args...)")
	wln("log.Error(err)")
	wln("return 0, nil, err")
	wln("}")

	wln("return total, data, nil")
}
