package prepare

import (
	"strings"

	"github.com/arstd/light/domain"
	"github.com/arstd/log"
)

func getReturnings(m *domain.Method) (rs []*domain.VarType) {
	switch m.Kind {
	case domain.Batch, domain.Update, domain.Delete:
		return nil
	case domain.Insert:
		return getInsertReturnings(m)
	case domain.Get, domain.List, domain.Page:
		return getFieldsReturings(m)
	}
	return rs
}

func getInsertReturnings(m *domain.Method) (rs []*domain.VarType) {
	stmt := m.Fragments[len(m.Fragments)-1].Stmt

	fs := strings.Split(stmt[(strings.Index(stmt, "returning ")+len("returning ")):], ",")
	log.Debug(fs)
	rs = make([]*domain.VarType, len(fs))
	for i, f := range fs {
		f = strings.TrimSpace(f)
		// TODO model index ?= 1
		for _, vt := range m.Params[1].Fields {
			setTag(vt)
			if vt.Tag == f {
				rs[i] = vt
			}
		}
	}

	return rs
}

func getFieldsReturings(m *domain.Method) (rs []*domain.VarType) {
	stmt := m.Fragments[0].Stmt

	stmt = stmt[len("select "):strings.Index(stmt, " from ")]
	fs := strings.Split(stmt, ",")
	rs = make([]*domain.VarType, len(fs))
	for i, f := range fs {
		f = strings.TrimSpace(f)
		// TODO model index
		for _, vt := range m.Results[0].Fields {
			setTag(vt)
			if vt.Tag == f {
				rs[i] = vt
			}
		}
	}

	return rs
}

func setTag(vt *domain.VarType) {
	if vt.Var == "" || vt.Tag != "" {
		return
	}
	last := 0
	for i := 1; i < len(vt.Var); i++ {
		if vt.Var[i] >= 'A' && vt.Var[i] <= 'Z' {
			vt.Tag += vt.Var[last+1:i] + "_" + strings.ToLower(vt.Var[i:i+1])
			last = i
		}
	}
	vt.Tag = strings.ToLower(vt.Var[:1]) + vt.Tag + vt.Var[last+1:]
}
