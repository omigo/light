package main

import (
	"fmt"
	"strings"

	"github.com/arstd/log"
)

func prepareArgs(m *Method, fs []*Fragment) {
	for _, f := range fs {
		fillRange(m, f)
		fillArgs(m, f)
		prepareArgs(m, f.Fragments)
		fillCond(f)
	}
}

func fillCond(f *Fragment) {
	if f.Cond == "" && f.Bracket {
		var cs []string
		for _, arg := range f.Args {
			cs = append(cs, getCond(arg))
		}
		for _, x := range f.Fragments {
			x.Bracket = true
			fillCond(x)
			if x.Cond != "" {
				cs = append(cs, x.Cond)
			}
		}

		// 去重
		if len(cs) > 1 {
			last := 1
			for i := 1; i < len(cs); i++ {
				j := 0
				for ; j < last; j++ {
					if cs[j] == cs[i] {
						break
					}
				}
				if j == last {
					if last != i {
						cs[last] = cs[i]
					}
					last++
				}
			}
			cs = cs[:last]
		}

		f.Cond = strings.Join(cs, " && ")
	}
}

func getCond(arg *VarType) string {
	if arg.Slice != "" || arg.Array != "" || arg.Name == "map" {
		return fmt.Sprintf("len(%s) != 0", arg.Var)
	} else if arg.Pointer != "" {
		cond := fmt.Sprintf("%s != nil", arg.Var)
		typ := arg.Alias
		aliasStart, aliasEnd := "", ""
		if typ == "" {
			typ = arg.Name
		} else {
			aliasStart = arg.Alias + "("
			aliasEnd = ")"
		}
		switch typ {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"byte", "rune", "float32", "float64":
			return cond + " && " + aliasStart + "*" + arg.Var + aliasEnd + " != 0"
		case "string":
			return cond + " && " + aliasStart + "*" + arg.Var + aliasEnd + ` != ""`
		case "bool":
			return cond + " && " + aliasStart + "*" + arg.Var + aliasEnd
		}
	} else if arg.Path == "time" && arg.Name == "Time" {
		return fmt.Sprintf("!%s.IsZero()", arg.Var)
	} else {
		typ := arg.Alias
		aliasStart, aliasEnd := "", ""
		if typ == "" {
			typ = arg.Name
		} else {
			aliasStart = arg.Alias + "("
			aliasEnd = ")"
		}
		switch typ {
		case "int", "int8", "int16", "int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64",
			"byte", "rune", "float32", "float64":
			return aliasStart + arg.Var + aliasEnd + " != 0"
		case "string":
			return aliasStart + arg.Var + aliasEnd + ` != ""`
		case "bool":
			return aliasStart + arg.Var + aliasEnd
		}
	}
	log.Panic("unimplemented")
	return ""
}

func fillRange(m *Method, f *Fragment) {
	if f.Range == nil {
		return
	}

	sel := strings.Split(f.Range.Var, ".")

	for _, param := range m.Params {

		if sel[0] == param.Var {
			switch len(sel) {
			case 1:
				if param.Slice == "" && param.Array == "" {
					log.Panicf("variable `%s` must be slice or array for method `%s`", param.Var, m.Name)
				}

				tmp := *param
				tmp.Var = f.Range.Var
				*f.Range = tmp

				iter := tmp
				iter.Var = f.Iterator.Var
				iter.Slice = ""
				iter.Array = ""
				*f.Iterator = iter

			case 2:
				if len(param.Fields) == 0 {
					log.Panicf("varible `%s` no field `%s` for method %s", sel[0], sel[1], m.Name)
				}
				for _, field := range param.Fields {
					if field.Var == sel[1] {
						if field.Slice == "" && field.Array == "" {
							log.Panicf("variable `%s` must be slice or array for method `%s`", field.Var, m.Name)
						}

						tmp := *field
						tmp.Var = f.Range.Var
						*f.Range = tmp

						iter := tmp
						iter.Var = f.Iterator.Var
						iter.Slice = ""
						iter.Array = ""
						*f.Iterator = iter
					}
				}

			default:
				log.Panicf("variable `%s` not found for method %s", f.Range.Var, m.Name)
			}
		}
	}

	if f.Range.Name == "" {
		log.Panicf("variable `%s` not found for method `%s`", f.Range.Var, m.Name)
	}

}

func fillArgs(m *Method, f *Fragment) {
	for _, vt := range f.Args {
		sel := strings.Split(vt.Var, ".")

		if f.Range != nil {
			if sel[0] == f.Index.Var {
				*vt = *f.Index

			} else if sel[0] == f.Iterator.Var {
				switch len(sel) {
				case 1:
					tmp := *f.Iterator
					tmp.Var = vt.Var
					*vt = tmp

				case 2:
					if len(f.Iterator.Fields) == 0 {
						log.Panicf("varible `%s` no field `%s` for method %s", sel[0], sel[1], m.Name)
					}
					for _, field := range f.Iterator.Fields {
						if field.Var == sel[1] {
							tmp := *field
							tmp.Var = vt.Var
							*vt = tmp
						}
					}

				default:
					log.Panicf("variable `%s` not found for method %s", vt.Var, m.Name)
				}
			}
		}
		if vt.Name != "" {
			continue
		}

		for _, param := range m.Params {
			if param.Var == sel[0] {
				switch len(sel) {
				case 1:
					tmp := *param
					tmp.Var = vt.Var
					*vt = tmp

				case 2:
					if len(param.Fields) == 0 {
						log.Panicf("varible `%s` no field `%s` for method %s", sel[0], sel[1], m.Name)
					}
					for _, field := range param.Fields {
						if field.Var == sel[1] {
							tmp := *field
							tmp.Var = vt.Var
							*vt = tmp
						}
					}

				default:
					log.Panicf("variable `%s` not found for method %s", vt.Var, m.Name)
				}
			}
		}
		if vt.Name == "" {
			log.Panicf("variable `%s` not found for method %s", vt.Var, m.Name)
		}
	}
}
