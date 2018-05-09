package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/arstd/light/goparser"
	"github.com/arstd/light/sqlparser"
	"github.com/arstd/log"
)

func Generate(store *goparser.Store) []byte {
	gopath := os.Getenv("GOPATH")
	for _, v := range strings.Split(gopath, string(filepath.ListSeparator)) {
		v = strings.TrimRight(v, "/")
		if strings.HasPrefix(store.Source, v) {
			store.Source = store.Source[len(v)+5:]
			break
		}
	}

	for k, v := range store.Imports {
		if i := strings.Index(k, "/vendor/"); i > 0 {
			delete(store.Imports, k)
			store.Imports[k[i+8:]] = v
		}
	}

	for _, m := range store.Methods {
		p := sqlparser.NewParser(bytes.NewBufferString(m.Doc))
		stmt, err := p.Parse()
		if err != nil {
			panic(err)
		}
		m.Statement = stmt
		// log.JSONIndent(stmt)

		m.GenCondition()
		m.SetType()
	}

	var t *template.Template
	t = template.New("tpl")
	t.Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"aggregate": func(m *goparser.Method, v *sqlparser.Fragment) *Aggregate {
			return &Aggregate{Method: m, Fragment: v}
		},
		"MethodTx":             goparser.MethodTx,
		"MethodSignature":      goparser.MethodSignature,
		"ParamsVarByNameValue": goparser.ParamsVarByNameValue,
		"VariableTypeName":     goparser.VariableTypeName,
		"VariableWrap":         goparser.VariableWrap,
		"VariableElemTypeName": goparser.VariableElemTypeName,
		"VariableVarByTagScan": goparser.VariableVarByTagScan,
	})
	log.Fataln(t.Parse(tpl))
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))
	log.Fataln(t.Execute(buf, store))
	return buf.Bytes()
}

type Aggregate struct {
	Method   *goparser.Method
	Fragment *sqlparser.Fragment
}
