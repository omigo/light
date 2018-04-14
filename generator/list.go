package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

var textList = `
query := buf.String()
{{- if .Method.Store.Log }}
	log.Debug(query)
	log.Debug(args...)
{{- end}}

ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
rows, err := exec.QueryContext(ctx, query, args...)
if err != nil {
	{{- if .Method.Store.Log }}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
	{{- end}}
	return nil, err
}
defer rows.Close()

var data {{call .Method.ResultTypeName}}

for rows.Next() {
	xu := new({{ call .Method.ResultElemTypeName }})
	data = append(data, xu)
	xdst := []interface{}{ {{- $method := .Method -}}
		{{- range $i, $field := .Statement.Fields -}}
			{{- if $i -}} , {{- end -}}
			{{- call $method.ResultVarByTagScan $field -}}
		{{- end -}}
	}

	err = rows.Scan(xdst...)
	if err != nil {
		{{- if .Method.Store.Log }}
			log.Error(query)
			log.Error(args...)
			log.Error(err)
		{{- end}}
		return nil, err
	}
	{{- if .Method.Store.Log }}
		log.JSON(xdst)
	{{- end}}
}
if err = rows.Err(); err != nil {
	{{- if .Method.Store.Log }}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
	{{- end}}
	return nil, err
}

return data, nil
`
var tplList = template.Must(template.New("textList").Parse(textList))

func writeList(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(tplList.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
