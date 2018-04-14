package generator

import (
	"bytes"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

const textFragment = `
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
`

var tplFragment *template.Template

func init() {
	tplFragment = template.New("textFragment")
	tplFragment.Funcs(template.FuncMap{
		"aggregate": func(m *goparser.Method, v *sqlparser.Fragment) *Aggregate {
			return &Aggregate{Method: m, Fragment: v}
		},
		"paramsVarByNameValue": func(m *goparser.Method, name string) string {
			x := m.Params.VarByName(name)
			return x.Value(x.VName)
		},
	})
	log.Fataln(tplFragment.Parse(textFragment))
}

type Aggregate struct {
	Method   *goparser.Method
	Fragment *sqlparser.Fragment
}

func writeFragment(buf *bytes.Buffer, m *goparser.Method, v *sqlparser.Fragment) {
	data := &Aggregate{Method: m, Fragment: v}
	tplFragment.Execute(buf, data)
}
