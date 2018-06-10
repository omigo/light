package goparser

import (
	"go/types"
)

type Results struct {
	*Params

	Result *Variable
}

func NewResults(tuple *types.Tuple) *Results {
	rs := &Results{Params: NewParams(tuple)}
	switch tuple.Len() {
	case 1:
		// ddl
	case 2:
		rs.Result = rs.List[0]
	case 3:
		rs.Result = rs.List[1]
	default:
		panic(len(rs.List))
	}
	return rs
}
