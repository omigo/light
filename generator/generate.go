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
		if strings.HasPrefix(store.Source, v) {
			store.Source = store.Source[len(v)+5:]
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

		genCondition(stmt, m)
		// log.JSONIndent(stmt)

		m.SetType()
	}

	var t *template.Template
	t = template.New("tpl")
	t.Funcs(template.FuncMap{
		"sub": func(a, b int) int { return a - b },
		"aggregate": func(m *goparser.Method, v *sqlparser.Fragment) *Aggregate {
			return &Aggregate{Method: m, Fragment: v}
		},
		"paramsVarByNameValue": func(m *goparser.Method, name string) string {
			x := m.Params.VarByName(name)
			return x.Value(x.VName)
		},
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

func genCondition(stmt *sqlparser.Statement, m *goparser.Method) {
	for _, v := range stmt.Fragments {
		deepGenCondition(v, m)
	}
}

func deepGenCondition(f *sqlparser.Fragment, m *goparser.Method) {
	if len(f.Fragments) == 0 {
		if f.Condition == "-" {
			var cs []string
			for _, name := range f.Variables {
				v := m.Params.VarByName(name)
				d := v.NotDefault(v.VName)
				cs = append(cs, "("+d+")")
			}
			f.Condition = strings.Join(cs, " && ")
		}
		return
	}

	for _, v := range f.Fragments {
		deepGenCondition(v, m)
	}

	if f.Condition != "-" {
		return
	}

	var cs []string
	for _, v := range f.Fragments {
		if v.Condition == "" {
			continue
		}
		cs = append(cs, "("+v.Condition+")")
	}
	f.Condition = strings.Join(cs, " || ")
}
