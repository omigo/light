package generator

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/omigo/light/goparser"
	"github.com/omigo/light/sqlparser"
	"github.com/omigo/log"
)

func getGomodPath(path string) string {
	for {
		path = filepath.Dir(path)
		if path == "" || path == "/" || path == "." {
			return path
		}
		if fileInfo, err := os.Stat(path + "/go.mod"); err != nil {
			if os.IsExist(err) {
				return path
			}
		} else if !fileInfo.IsDir() {
			return path
		}
	}
}

func Generate(itf *goparser.Interface) []byte {
	path := getGomodPath(itf.Source)

	if i := strings.LastIndex(path, string(filepath.Separator)); i > 0 {
		path = path[:i]

		if strings.HasPrefix(itf.Source, path) {
			itf.Source = itf.Source[len(path)+1:]
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
		"ParamsLast":          goparser.ParamsLast,
		"ParamsLastElem":      goparser.ParamsLastElem,
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
