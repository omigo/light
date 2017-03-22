package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/arstd/log"
)

func getFragments(doc string) (fs []*Fragment) {
	// log.Infof(doc)
	// time.Sleep(200 * time.Millisecond)

	ignore, left, last := false, 0, -1
	for i, c := range doc {
		if ignore {
			if c == '\'' {
				ignore = false
			}
			continue
		}

		switch c {
		default:
		case '\'':
			ignore = true

		case '[':
			if left == 0 {
				if strings.HasSuffix(doc[last+1:i], "array") {
					// array[ 之前没有 [，之后的会被认为普通字符
					break
				}
				f := &Fragment{
					Stmt: strings.TrimSpace(doc[last+1 : i]),
				}
				if parseFragment(f) {
					fs = append(fs, f)
				}
				last = i
			}
			left++

		case ']':
			left--
			if left == 0 {
				f := &Fragment{
					Stmt: strings.TrimSpace(doc[last+1 : i]),
				}

				if parseFragment(f) {
					fs = append(fs, f)
					last = i
				}
			}
		}
	}

	f := &Fragment{
		Stmt: strings.TrimSpace(doc[last+1:]),
	}

	if parseFragment(f) {
		fs = append(fs, f)
	}

	return fs
}

var condRegexp = regexp.MustCompile(`((\w+),\s*(\w+)\s*:=\s*)?range\s+([\w.]+)(\s*\|\s*(.+))?`)

func parseFragment(f *Fragment) bool {
	// log.Debug(f.Stmt)
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
			f.Index = &VarType{Var: "i", Name: "int"}
			f.Iterator = &VarType{Var: "v"}
			if m[2] != "" {
				f.Index.Var = m[2]
			}
			if m[3] != "" {
				f.Iterator.Var = m[3]
			}
			f.Seperator = ","
			if m[6] != "" {
				f.Seperator = m[6]
			}
			f.Range = &VarType{Var: m[4]}
			f.Cond = fmt.Sprintf("%s, %s := range %s", f.Index.Var, f.Iterator.Var, f.Range.Var)

			if f.Stmt == "" {
				f.Stmt = "${" + f.Iterator.Var + "}"
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

func parseArgs(f *Fragment) bool {
	buf := &bytes.Buffer{}

	quote, left, last := false, 0, -1
	for i, c := range f.Stmt {
		if quote {
			if c == '\'' {
				quote = false
			}
			continue
		}

		switch c {
		default:
		case '\'':
			quote = true

		case '{':
			if i > 0 && f.Stmt[i-1] == '$' {
				buf.WriteString(f.Stmt[last+1 : i-1])
				last = i
				left++
			}

		case '}':
			left--
			if left == 0 {
				a := &VarType{
					Var: strings.TrimSpace(f.Stmt[last+1 : i]),
				}
				f.Args = append(f.Args, a)

				buf.WriteString("%s")
				last = i
			}
		}
	}
	buf.WriteString(f.Stmt[last+1:])

	f.Prepare = buf.String()
	return true
}
