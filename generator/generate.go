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

func Generate(itf *goparser.Interface) []byte {
	gopath := os.Getenv("GOPATH")
	for _, v := range strings.Split(gopath, string(filepath.ListSeparator)) {
		v = strings.TrimRight(v, "/")
		if strings.HasPrefix(itf.Source, v) {
			itf.Source = itf.Source[len(v)+5:]
			break
		}
	}

	for k, v := range itf.Imports {
		if i := strings.Index(k, "/vendor/"); i > 0 {
			delete(itf.Imports, k)
			itf.Imports[k[i+8:]] = v
		}
	}

	for _, m := range itf.Methods {
		p := sqlparser.NewParser(bytes.NewBufferString(m.Doc))
		stmt, err := p.Parse()
		if err != nil {
			panic(err)
		}
		m.Statement = stmt
		// log.JSONIndent(stmt)

		m.GenCondition()
		m.SetType()

		m.SetSignature()
	}

	var t *template.Template
	t = template.New("tpl")
	t.Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"aggregate": func(m *goparser.Method, v *sqlparser.Fragment, buf, args string) *Aggregate {
			return &Aggregate{Method: m, Fragment: v, Buf: buf, Args: args}
		},
		"MethodTx":            goparser.MethodTx,
		"HasVariable":         goparser.HasVariable,
		"ResultWrap":          goparser.ResultWrap,
		"ResultTypeName":      goparser.ResultTypeName,
		"ResultElemTypeName":  goparser.ResultElemTypeName,
		"LookupScanOfResults": goparser.LookupScanOfResults,
		"LookupValueOfParams": goparser.LookupValueOfParams,
	})
	log.Fataln(t.Parse(tpl))
	buf := bytes.NewBuffer(make([]byte, 0, 1024*16))
	log.Fataln(t.Execute(buf, itf))
	return buf.Bytes()
}

type Aggregate struct {
	Method   *goparser.Method
	Fragment *sqlparser.Fragment
	Buf      string
	Args     string
}
