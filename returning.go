package main

import (
	"strings"

	"github.com/arstd/log"
)

func getReturnings(m *Method) (rs []*VarType) {
	switch m.Kind {
	case Count:
		return nil
	case Batch, Update, Delete:
		return nil
	case Insert:
		return getInsertReturnings(m)
	case Get, List:
		return getFieldsReturings(m, 0)
	case Page:
		return getFieldsReturings(m, 1)
	}
	return rs
}

func getInsertReturnings(m *Method) (rs []*VarType) {
	// log.JSONIndent(m)
	stmt := m.Fragments[len(m.Fragments)-1].Stmt

	log.Print(m.Name, stmt)

	idx := strings.Index(stmt, "returning ")
	if idx == -1 {
		return nil
	}

	fs := strings.Split(stmt[(idx+len("returning ")):], ",")
	rs = make([]*VarType, len(fs))
	for i, f := range fs {
		f = strings.TrimSpace(f)
		for _, p := range m.Params {
			if p.Fields == nil {
				break
			}
			for _, vt := range p.Fields {
				if strings.HasPrefix(vt.Tag, f) {
					rs[i] = vt
					break
				}
			}
		}
		if rs[i] == nil {
			log.Panicf("returning `%s` no matched field for method `%s`", f, m.Name)
		}
	}

	return rs
}

func getFieldsReturings(m *Method, idx int) (rs []*VarType) {
	stmt := m.Fragments[0].Stmt

	stmt = stmt[len("select "):strings.Index(stmt, " from ")]
	fs := strings.Split(stmt, ",")
	rs = make([]*VarType, len(fs))
	for i, f := range fs {
		fs := strings.Split(f, " ")
		f = fs[len(fs)-1]
		fs = strings.Split(f, ".")
		f = fs[len(fs)-1]
		f = strings.TrimSpace(f)
		// TODO model index
		for _, vt := range m.Results[idx].Fields {
			if strings.HasPrefix(vt.Tag, f) {
				rs[i] = vt
				break
			}
		}
		if rs[i] == nil {
			log.Panicf("returning `%s` no matched field for method `%s`", f, m.Name)
		}
	}

	return rs
}
