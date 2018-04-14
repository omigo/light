package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textDelete = `
query := buf.String()
{{- if .Method.Store.Log }}
	log.Debug(query)
	log.Debug(args...)
{{- end}}

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
res, err := exec.ExecContext(ctx, query, args...)
if err != nil {
	{{- if .Method.Store.Log }}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
	{{- end}}
	return 0, err
}
return res.RowsAffected()
`

var tplDelete = template.Must(template.New("textDelete").Parse(textDelete))

func writeDelete(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	log.Errorn(tplDelete.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
