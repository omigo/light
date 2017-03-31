package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"

	"github.com/arstd/log"
)

func splitToFragments(doc string) (fs []*Fragment) {
	log.Debug(doc)

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
				part := strings.TrimSpace(doc[last+1 : i])
				last = i
				if part != "" { // ] 和 [ 中间没有内容
					// 如果有内容，只可能是简单语句，没有 {...} 和 [...]
					f := &Fragment{
						Cond: "",
						Stmt: part,
					}
					fs = extractArgs(fs, f)
				}
			}
			left++

		case ']':
			left--
			if left == 0 {
				part := strings.TrimSpace(doc[last+1 : i])
				last = i
				if part != "" { // [ 和 ] 之间没有内容
					// part 可能有 {...}, 可能嵌套有 [...]
					cond, sub := divideCond(part)
					f := checkRange(cond, sub)
					if sub == "" {
						if f.Range == nil {
							log.Panicf("miss statement: %s", part)
						}
						fs = extractArgs(fs, f)
					} else {
						nests := splitToFragments(sub)
						if len(nests) == 0 {
							log.Panicf("expect fragment(s), but no: %s", part)
						} else if len(nests) == 1 {
							f.Stmt = nests[0].Stmt
							f.Prepare = nests[0].Prepare
							f.Args = nests[0].Args
							f.Fragments = nests[0].Fragments
							fs = append(fs, f)
						} else {
							f.Fragments = nests
							fs = append(fs, f)
						}
					}
				}
			}
		}
	}
	if last != -1 && doc[last] == '[' {
		log.Panicf("miss `[` to match left bracket: %s", doc[last:])
	}

	part := strings.TrimSpace(doc[last+1:])
	if part != "" { // 已经到了末尾
		// 如果有内容，只可能是简单语句，没有 {...} 和 [...]
		f := &Fragment{
			Cond: "",
			Stmt: part,
		}
		fs = extractArgs(fs, f)
	}

	return fs
}

func divideCond(part string) (cond string, sub string) {
	if part[0] == '{' {
		var left int
		for i, c := range part {
			if c == '\'' {
				continue
			}
			if c == '{' {
				left++
			}
			if c == '}' {
				left--
				if left == 0 {
					return strings.TrimSpace(part[1:i]), strings.TrimSpace(part[i+1:])
				}
			}
		}
	}
	return "", part
}

var rre = regexp.MustCompile(`((\w+),\s*(\w+)\s*:=\s*)?range\s+([\w.]+)(\s*\|\s*(.+))?`)

func checkRange(cond, sub string) (f *Fragment) {
	f = &Fragment{Cond: cond, Stmt: sub, Bracket: true}
	if m := rre.FindStringSubmatch(cond); len(m) > 0 {
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
		f.Cond = fmt.Sprintf("len(%s) != 0", f.Range.Var)

		if f.Stmt == "" {
			f.Stmt = "${" + f.Iterator.Var + "}"
		}
	}
	return
}

func extractArgs(fs []*Fragment, f *Fragment) []*Fragment {
	fs = append(fs, f)
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
			if i > 0 {
				if f.Stmt[i-1] == '$' {
					buf.WriteString(f.Stmt[last+1 : i-1])
					last = i
					left++
				} else if f.Stmt[i-1] == '#' {
					buf.WriteString(f.Stmt[last+1 : i-1])
					last = i
					left--
					log.Debug(buf.String())
				}
			}

		case '}':
			if left == 1 {
				a := &VarType{
					Var: strings.TrimSpace(f.Stmt[last+1 : i]),
				}
				f.Args = append(f.Args, a)

				buf.WriteString("%s")
				last = i
			} else if left == -1 {
				fs = append(fs, &Fragment{
					Stmt:    strings.TrimSpace(f.Stmt[last-1 : i+1]),
					Hashtag: true,
					Args: []*VarType{
						{
							Var: strings.TrimSpace(f.Stmt[last+1 : i]),
						},
					},
				})
				fs = extractArgs(fs, &Fragment{
					Stmt: strings.TrimSpace(f.Stmt[i+1:]),
				})
				f.Stmt = strings.TrimSpace(f.Stmt[:last-1])
				f.Prepare = strings.TrimSpace(buf.String())
				return fs
			}

			left = 0
		}
	}
	buf.WriteString(strings.TrimSpace(f.Stmt[last+1:]))

	f.Prepare = buf.String()
	return fs
}
