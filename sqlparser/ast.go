package sqlparser

type Stmt interface {
	Type() Token
}

type SelectStmt struct {
	Fields []string

	Fragments []*Fragment

	Offset string
	Limit  string
}

func (*SelectStmt) Type() Token { return SELECT }

type Fragment struct {
	cond      bool
	Cond      string
	Stmt      string
	Variables []string
	Fragments []*Fragment
}
