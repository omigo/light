package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textGet = `
query := buf.String()
{{- if .Method.Store.Log }}
	log.Debug(query)
	log.Debug(args...)
{{- end}}

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
row := exec.QueryRowContext(ctx, query, args...)

xu := new({{call .Method.ResultTypeName}})
xdst := []interface{}{ {{- $method := .Method -}}
	{{- range $i, $field := .Statement.Fields -}}
		{{- if $i -}} , {{- end -}}
		{{- call $method.ResultVarByTagScan $field -}}
	{{- end -}}
}
err := row.Scan(xdst...)
	if err != nil {
	if err == sql.ErrNoRows {
		return nil, nil
	}
{{- if .Method.Store.Log}}
	log.Error(query)
	log.Error(args...)
	log.Error(err)
{{- end }}
	return nil, err
{{- if .Method.Store.Log}}
	log.JSON(xdst)
{{- end }}
}

return xu, err
`

var tplGet = template.Must(template.New("textGet").Parse(textGet))

func writeGet(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(tplGet.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
