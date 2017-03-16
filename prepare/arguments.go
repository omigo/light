package prepare

import (
	"strings"

	"github.com/arstd/light/domain"
	"github.com/arstd/log"
)

func prepareArgs(m *domain.Method, fs []*domain.Fragment) {
	for _, f := range fs {
		fillRange(m, f)
		fillArgs(m, f)
		prepareArgs(m, f.Fragments)
	}
}

func fillRange(m *domain.Method, f *domain.Fragment) {
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

						tmp := *param
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

func fillArgs(m *domain.Method, f *domain.Fragment) {
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
