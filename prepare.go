package main

import (
	"strings"

	"github.com/arstd/log"
)

func prepare(pkg *Package) {
	for _, intf := range pkg.Interfaces {
		for _, m := range intf.Methods {
			addPathToImports(pkg, m)
			m.Kind = getMethodKind(m)

			checkLastParamTx(m)
			fillResultsVar(m)

			m.Fragments = splitToFragments(m.Doc)
			prepareArgs(m, m.Fragments)

			m.Returnings = getReturnings(m)
		}
	}
}

func checkLastParamTx(m *Method) {
	vt := m.Params[len(m.Params)-1]
	if vt.Name == "Tx" && vt.Path == "database/sql" && vt.Pointer == "*" && vt.Slice == "[]" {
		vt.Slice = "..."
		vt.Var = "xtx"
	} else {
		log.Panicf("last param expect `txs ...*sql.Tx`, but `%s` not", m.Name)
	}
}

func addPathToImports(pkg *Package, m *Method) {
	// TODO conflict
	for _, p := range m.Params {
		if p.Path != "" {
			pkg.Imports[p.Path[strings.LastIndex(p.Path, "/")+1:]] = p.Path
		}
		for _, f := range p.Fields {
			if f.Path != "" {
				pkg.Imports[f.Path[strings.LastIndex(f.Path, "/")+1:]] = f.Path
			}
		}
	}
	for _, p := range m.Results {
		if p.Path != "" {
			pkg.Imports[p.Path[strings.LastIndex(p.Path, "/")+1:]] = p.Path
		}
		for _, f := range p.Fields {
			if f.Path != "" {
				pkg.Imports[f.Path[strings.LastIndex(f.Path, "/")+1:]] = f.Path
			}
		}
	}
}

func getMethodKind(m *Method) MethodKind {
	if len(m.Results) < 0 {
		log.Panicf("all metheds must have 1-3 returns, but %s no return", m.Name)
	}

	if len(m.Results) > 3 {
		log.Panicf("all metheds must have 1-3 returns, but method '%s' has %d returns", m.Name, len(m.Results))
	}

	if m.Results[len(m.Results)-1].Name != "error" {
		log.Panicf("method '%s' last return must error", m.Name)
	}

	i := strings.IndexAny(m.Doc, " \t")
	if i == -1 {
		log.Panicf("sql error for method '%s', must has one or more space", m.Name)
	}

	head := strings.ToLower(m.Doc[:i])
	switch head {
	default:
		log.Panicf("sql error for method '%s', must has prefix insert/update/delete/select keyword", m.Name)

	case "insert":
		if len(m.Results) == 1 {
			return Insert
		} else if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return Batch
		} else {
			log.Panicf("method '%s' for insert must only return 'error'", m.Name)
		}

	case "update":
		if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return Update
		} else {
			log.Panicf("method '%s' for 'update' must only return '(int64, error)'", m.Name)
		}

	case "delete":
		if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return Delete
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, error)'", m.Name)
		}

	case "select":
		// get/count/list/page
	}

	if len(m.Results) == 2 {
		if m.Results[0].Slice != "" {
			return List
		} else if len(m.Results[0].Fields) > 0 {
			return Get
		} else {
			return Count
		}
	}

	if len(m.Results) == 3 {
		if m.Results[0].Name == "int64" && m.Results[1].Slice != "" {
			return Page
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, []<*struct>, error)'", m.Name)
		}
	}

	panic("unreachable code")
}

func fillResultsVar(m *Method) {
	m.Results[len(m.Results)-1].Var = "err"

	switch m.Kind {
	case Insert:

	case Batch, Update, Delete:
		m.Results[0].Var = "xa"

	case Get:
		m.Results[0].Var = "xobj"

	case Count:
		m.Results[0].Var = "xcnt"

	case List:
		m.Results[0].Var = "xdata"

	case Page:
		m.Results[0].Var = "xcnt"
		m.Results[1].Var = "xdata"

	default:
		log.Panicf("unimplements method kind %s", m.Kind)
	}
}
