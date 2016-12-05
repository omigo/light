package main

import (
	"strings"

	"github.com/arstd/log"
)

func (a *Analyzer) parse() {
	// var err error
	for _, intf := range a.Interfaces {
		for _, m := range intf.Methods {
			log.Debug(m.Comment)
			m.Kind = getMethodKind(m)
			log.Debug(m.Kind)
		}
	}
}

func getMethodKind(m *Method) MethodKind {
	i := strings.IndexAny(m.Comment, " \t")
	if i == -1 {
		panic("sql error for method '" + m.Name + "', must has one or more space")
	}

	switch strings.ToLower(m.Comment[:i]) {
	case "insert":
		return KindInsert
	case "update":
		return KindUpdate
	case "delete":
		return KindDelete
	case "select":
		// get count list page
	default:
		panic("sql error for method '" + m.Name + "', must has prefix insert/update/delete/select keyword")
	}

	if len(m.Results) == 2 {
		if m.Results[0].IsPrimitive() {
			return KindCount
		} else if m.Results[0].IsStruct() {
			return KindGet
		} else if m.Results[0].IsArray() {
			return KindList
		} else {
			panic("error result type for '" + m.Name + "'")
		}
	}

	if len(m.Results) == 3 {
		if m.Results[0].IsPrimitive() && m.Results[1].IsArray() {
			return KindList
		} else {
			panic("error result type for '" + m.Name + "', page method must return (total int64, data []*Struct, err error)")
		}
	}

	if len(m.Results) > 3 {
		panic("results error for method '" + m.Name + "', too many arguments to return")
	}

	if len(m.Results) < 2 {
		panic("results error for method '" + m.Name + "', not enough arguments to return")
	}

	panic("unreachable code")
}
