package generator

import (
	"bytes"
	"html/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const CountResult = `
	query := buf.String()
{{- if .Store.Log}}
	log.Debug(query)
	log.Debug(args...)
{{end -}}
	var count {{call .ResultTypeName}}
	err := db.QueryRow(query, args...).Scan({{call .ResultTypeWrap}}(&count))
	if err != nil {
{{- if .Store.Log}}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
{{end -}}
		return count, err
	}
{{- if .Store.Log}}
	log.Debug(count)
{{end -}}
	return count, nil
`

var countResultTpl = template.Must(template.New("tplCountResult").Parse(CountResult))

func writeCount(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(countResultTpl.Execute(buf, m))
}
