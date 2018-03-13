package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textCount = `
	query := buf.String()
{{- if .Store.Log}}
	log.Debug(query)
	log.Debug(args...)
{{end -}}
	var count {{call .ResultTypeName}}
	err := exec.QueryRow(query, args...).Scan({{call .ResultTypeWrap}}(&count))
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

var tplCount = template.Must(template.New("textCount").Parse(textCount))

func writeCount(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(tplCount.Execute(buf, m))
}
