package sqlparser

type Statement struct {
	Type   Token
	Table  string
	Fields []string

	Fragments []*Fragment `json:"fragments,omitempty"`
}

type Fragment struct {
	Condition string `json:"cond,omitempty"`
	Range     string `json:"range,omitempty"`
	Open      string `json:"open,omitempty"`
	Close     string `json:"close,omitempty"`

	Statement string   `json:"stmt,omitempty"`
	Replacers []string `json:"replacers,omitempty"`
	Variables []string `json:"variables,omitempty"`

	Fragments []*Fragment `json:"fragments,omitempty"`
}
