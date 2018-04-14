package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textPage = `
{{- define "textFragment" -}}
	{{- if .Fragment.Condition}}
		if {{.Fragment.Condition}} {
	{{- end }}
	{{- if .Fragment.Statement }}
		{{- if .Fragment.Range }}
			if len({{.Fragment.Range}}) > 0 {
				fmt.Fprintf(&buf, "{{.Fragment.Statement}} ", strings.Repeat(",?", len({{.Fragment.Range}}))[1:])
				for _, v := range {{.Fragment.Range}} {
					args = append(args, v)
				}
			}
		{{- else if .Fragment.Replacers }}
			fmt.Fprintf(&buf, "{{.Fragment.Statement}} "{{range $elem := .Fragment.Replacers}}, {{$elem}}{{end}})
		{{- else }}
			buf.WriteString("{{.Fragment.Statement}} ")
		{{- end }}
		{{- if .Fragment.Variables }}{{$method := .Method}}
			args = append(args{{range $elem := .Fragment.Variables}}, {{paramsVarByNameValue $method $elem}}{{end}})
		{{- end }}
	{{- else }}{{$method := .Method}}
		{{- range $fragment := .Fragment.Fragments }}
			{{template "textFragment" (aggregate $method $fragment)}}
		{{- end }}
	{{- end }}
	{{- if .Fragment.Condition}}
	 }
	{{- end }}
{{- end -}}

var total int64
totalQuery := "SELECT count(1) "+ buf.String()
{{- if .Method.Store.Log }}
	log.Debug(totalQuery)
	log.Debug(args...)
{{- end}}
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
err := exec.QueryRowContext(ctx, totalQuery, args...).Scan(&total)
if err != nil {
	{{- if .Method.Store.Log }}
		log.Error(totalQuery)
		log.Error(args...)
		log.Error(err)
	{{- end}}
	return 0, nil, err
}
{{- if .Method.Store.Log }}
log.Debug(total)
{{- end}}

{{$i := sub (len .Statement.Fragments) 1}}
{{ $fragment := index .Statement.Fragments $i }}
{{template "textFragment" (aggregate .Method $fragment)}}

{{ $fragement0 := index .Statement.Fragments 0 }}
query := "{{$fragement0.Statement}} " + buf.String()
{{- if .Method.Store.Log }}
	log.Debug(query)
	log.Debug(args...)
{{- end}}

ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()
rows, err := exec.QueryContext(ctx, query, args...)
if err != nil {
	{{- if .Method.Store.Log }}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
	{{- end}}
	return 0, nil, err
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
		return 0, nil, err
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
	return 0, nil, err
}

return total, data, nil
`

var tplPage *template.Template

func init() {
	tplPage = template.New("textPage")
	tplPage.Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"aggregate": func(m *goparser.Method, v *sqlparser.Fragment) *Aggregate {
			return &Aggregate{Method: m, Fragment: v}
		},
		"paramsVarByNameValue": func(m *goparser.Method, name string) string {
			x := m.Params.VarByName(name)
			return x.Value(x.VName)
		},
	})
	log.Fataln(tplPage.Parse(textPage))
}

func writePage(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(tplPage.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
