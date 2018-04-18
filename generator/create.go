package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textCreate = `
query := buf.String()
{{- if .Method.Store.Log }}
	log.Debug(query)
	log.Debug(args...)
{{- end}}

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
_, err := exec.ExecContext(ctx, query, args...)
{{- if .Method.Store.Log }}
	if err != nil {
		log.Error(query)
		log.Error(args...)
		log.Error(err)
	}
{{- end}}
return err
`

var tplCreate = template.Must(template.New("textCreate").Parse(textCreate))

func writeCreate(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	wln := func(s string) { buf.WriteString(s + "\n") }

	wln("var buf bytes.Buffer")
	wln("var args []interface{}")

	for _, f := range stmt.Fragments {
		writeFragment(buf, m, f)
	}

	log.Errorn(tplCreate.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
