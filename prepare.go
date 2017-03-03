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
	if len(m.Results) < 0 {
		log.Panicf("all metheds must have 1-3 returns, but %s no return", m.Name)
	}

	if len(m.Results) > 3 {
		log.Panicf("all metheds must have 1-3 returns, but method '%s' has %d returns", m.Name, len(m.Results))
	}

	if m.Results[len(m.Results)-1].Type != "error" {
		log.Panicf("method '%s' last return must error", m.Name)
	}

	i := strings.IndexAny(m.Comment, " \t")
	if i == -1 {
		log.Panicf("sql error for method '%s', must has one or more space", m.Name)
	}

	head := strings.ToLower(m.Comment[:i])
	switch head {
	default:
		log.Panicf("sql error for method '%s', must has prefix insert/update/delete/select keyword", m.Name)

	case "insert":
		if len(m.Results) == 1 {
			return KindInsert
		} else {
			log.Panicf("method '%s' for insert must only return 'error'", m.Name)
		}

	case "update":
		if len(m.Results) == 2 && m.Results[0].Type == "int64" {
			return KindUpdate
		} else {
			log.Panicf("method '%s' for 'update' must only return '(int64, error)'", m.Name)
		}

	case "delete":
		if len(m.Results) == 2 && m.Results[0].Type == "int64" {
			return KindDelete
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, error)'", m.Name)
		}

	case "select":
		// get/count/list/page
	}

	if len(m.Results) == 2 {
		if m.Results[0].IsStruct() {
			return KindGet
		} else if m.Results[0].IsArray() {
			return KindList
		} else {
			return KindCount
		}
	}

	if len(m.Results) == 3 {
		if m.Results[0].Type == "int64" && m.Results[1].IsArray() {
			return KindList
		} else {
			log.Panicf("method '%s' for 'delete' must only return '(int64, []<*struct>, error)'", m.Name)
		}
	}

	panic("unreachable code")
}
