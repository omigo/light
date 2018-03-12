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
{{- if .Log}}
	log.Debug(query)
	log.Debug(args...)
{{end -}}
	var count {{.ResultTypeName}}
	err := db.QueryRow(query, args...).Scan({{.ResultTypeWrap}}(&count))
	if err != nil {
{{- if .Log}}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
{{end -}}
		return count, err
	}
{{- if .Log}}
	log.Debug(count)
{{end -}}
	return count, nil
`

var countResultTpl = template.Must(template.New("tplCountResult").Parse(CountResult))

func writeCount(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	data := map[string]interface{}{
		"Log":            m.Store.Log,
		"ResultTypeName": m.Results.Result().TypeName(),
		"ResultTypeWrap": m.Results.Result().Wrap(true),
	}
	log.Errorn(countResultTpl.Execute(buf, data))
}
