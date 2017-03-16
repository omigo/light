package prepare

import (
	"strings"

	"github.com/arstd/light/domain"
	"github.com/arstd/log"
)

func Prepare(pkg *domain.Package) {
	// var err error
	for _, intf := range pkg.Interfaces {
		for _, m := range intf.Methods {

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

			m.Kind = getMethodKind(m)
			m.Fragments = getFragments(m.Doc)
			prepareArgs(m, m.Fragments)
			m.Returnings = getReturnings(m)
		}
	}
}

func getMethodKind(m *domain.Method) domain.MethodKind {
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
			return domain.Insert
		} else if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return domain.Batch
		} else {
			log.Panicf("method '%s' for insert must only return 'error'", m.Name)
		}

	case "update":
		if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return domain.Update
		} else {
			log.Panicf("method '%s' for 'update' must only return '(int64, error)'", m.Name)
		}

	case "delete":
		if len(m.Results) == 2 && m.Results[0].Name == "int64" {
			return domain.Delete
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, error)'", m.Name)
		}

	case "select":
		// get/count/list/page
	}

	if len(m.Results) == 2 {
		if m.Results[0].Slice != "" {
			return domain.List
		} else if len(m.Results[0].Fields) > 0 {
			return domain.Get
		} else {
			return domain.Count
		}
	}

	if len(m.Results) == 3 {
		if m.Results[0].Name == "int64" && m.Results[1].Slice != "" {
			return domain.Page
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, []<*struct>, error)'", m.Name)
		}
	}

	panic("unreachable code")
}
