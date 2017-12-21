package sqlparser

type Statement struct {
	Type   Token
	Table  string
	Fields []string

	Fragments []*Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	Condition string `json:"cond,omitempty"`

	Statement string   `json:"stmt,omitempty"`
	Variables []string `json:"variables,omitempty"`

	Fragments []*Fragment `json:"fragments,omitempty"`
}
