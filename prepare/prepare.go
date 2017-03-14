package prepare

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/arstd/light/domain"
	"github.com/arstd/log"
)

func PrepareStmt(p *domain.Package) {
	// var err error
	for _, intf := range p.Interfaces {
		for _, m := range intf.Methods {
			m.Kind = getMethodKind(m)
			m.Fragments = getFragments(m.Doc)
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
		} else {
			log.Panicf("method '%s' for insert must only return 'error'", m.Name)
		}

	case "update":
		if len(m.Results) == 2 && m.Results[0].Type == "int64" {
			return domain.Update
		} else {
			log.Panicf("method '%s' for 'update' must only return '(int64, error)'", m.Name)
		}

	case "delete":
		if len(m.Results) == 2 && m.Results[0].Type == "int64" {
			return domain.Delete
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, error)'", m.Name)
		}

	case "select":
		// get/count/list/page
	}

	if len(m.Results) == 2 {
		if len(m.Results[0].Fields) > 0 {
			return domain.Get
		} else if m.Results[0].Slice != "" {
			return domain.List
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

func getFragments(doc string) (fs []*domain.Fragment) {
	log.Infof(doc)
	time.Sleep(200 * time.Millisecond)

	ignore, left, last := false, 0, -1
	for i, c := range doc {
		if ignore {
			if c == '\'' {
				ignore = false
			}
			continue
		}

		if c == '[' {
			if left == 0 {
				if strings.HasSuffix(doc[last+1:i], "array") {
					// array[ 之前没有 [，之后的会被认为普通字符
					continue
				}
				f := &domain.Fragment{
					Stmt: strings.TrimSpace(doc[last+1 : i]),
				}
				if parseFragment(f) {
					fs = append(fs, f)
				}
				last = i
			}
			left++
		} else if c == ']' {
			left--
			if left == 0 {
				f := &domain.Fragment{
					Stmt: strings.TrimSpace(doc[last+1 : i]),
				}

				if parseFragment(f) {
					fs = append(fs, f)
					last = i
				}
			}
		}
	}

	f := &domain.Fragment{
		Stmt: strings.TrimSpace(doc[last+1:]),
	}

	if parseFragment(f) {
		fs = append(fs, f)
	}

	return fs
}

var condRegexp = regexp.MustCompile(`((\w+),\s*(\w+)\s*:=\s*)?range\s+([\w.]+)(\s*\|\s*(.+))?`)

func parseFragment(f *domain.Fragment) bool {
	log.Debug(f.Stmt)

	f.Stmt = strings.TrimSpace(f.Stmt)
	if len(f.Stmt) == 0 {
		return false
	}

	if len(f.Stmt) < 1 {
		log.Panicf("sql error near by %s", f.Stmt)
	}
	if f.Stmt[0] == '{' {
		for i, c := range f.Stmt {
			// TODO must deal cond contain { or }
			if c == '}' {
				f.Cond = strings.TrimSpace(f.Stmt[1:i])
				f.Stmt = strings.TrimSpace(f.Stmt[i+1:])
				break
			}
		}
		if f.Cond == "" {
			log.Panicf("sql error near by %s", f.Stmt)
		}
	}

	if f.Cond != "" {
		if m := condRegexp.FindStringSubmatch(f.Cond); len(m) > 0 {
			f.Index, f.Iterator = "i", "v"
			if m[2] != "" {
				f.Index = m[2]
			}
			if m[3] != "" {
				f.Iterator = m[3]
			}
			f.Seperator = ","
			if m[6] != "" {
				f.Seperator = m[6]
			}
			f.Range = m[4]
			f.Cond = fmt.Sprintf("%s, %s := range %s", f.Index, f.Iterator, f.Range)

			if f.Stmt == "" {
				f.Stmt = "${" + f.Iterator + "}"
			}
		}
	}

	if f.Cond != "" || f.Stmt[0] == '[' {
		fs := getFragments(f.Stmt)
		if len(fs) > 1 || (len(fs) == 1 && fs[0].Cond != "") {
			f.Fragments = fs
			return true
		}
	}

	return parseArgs(f)
}

func parseArgs(f *domain.Fragment) bool {
	buf := &bytes.Buffer{}

	stack, top := make([]int, 32), 0

	last := -1
	for i, c := range f.Stmt {
		switch c {
		case '\'':
			if f.Stmt[stack[top]] == '\'' {
				stack[top] = 0
				top--
			}

		case '{':
			if f.Stmt[stack[top]] == '\'' {
				continue
			}
			if i >= 1 && f.Stmt[i-1] == '$' {
				if last+1 < i-2 {
					buf.WriteString(f.Stmt[last+1 : i-1])
				}

				top++
				stack[top] = i
			}

		case '}':
			if f.Stmt[stack[top]] == '\'' {
				continue
			}

			if f.Stmt[stack[top]] == '{' {
				a := &domain.VarType{
					Var: strings.TrimSpace(f.Stmt[stack[top]+1 : i]),
				}
				f.Args = append(f.Args, a)

				buf.WriteString("%s")
				last = i

				stack[top] = 0
				top--
			} else {
				log.Panicf("unexpected symbol `}`, not pair symbol `{`: %s", f.Stmt)
			}

		default: // do nothing
		}
	}
	if top != 0 {
		log.Panicf("parentheses do not match, expect `}`")
	}
	if last < len(f.Stmt) {
		buf.WriteString(f.Stmt[last+1:])
	}

	f.Prepare = buf.String()
	return true
}
