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
{{- if .Method.Store.Log}}
	log.Debug(query)
	log.Debug(args...)
{{- end}}
	var count {{call .Method.ResultTypeName}}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := exec.QueryRowContext(ctx, query, args...).Scan({{call .Method.ResultTypeWrap}}(&count))
	if err != nil {
		if err == sql.ErrNoRows {
			return count, nil
		}
{{- if .Method.Store.Log}}
		log.Error(query)
		log.Error(args...)
		log.Error(err)
{{- end}}
		return count, err
	}
{{- if .Method.Store.Log}}
	log.Debug(count)
{{- end}}
	return count, nil
`

var tplCount = template.Must(template.New("textCount").Parse(textCount))

type Wrapper struct {
	Method    *goparser.Method
	Statement *sqlparser.Statement
}

func writeCount(buf *bytes.Buffer, m *goparser.Method, stmt *sqlparser.Statement) {
	log.Errorn(tplCount.Execute(buf, &Wrapper{Method: m, Statement: stmt}))
}
