package sqlparser

import "strings"

// Token represents a lexical token.
type Token uint8

const (
	// Special tokens
	EOF Token = iota
	WS

	// Literals
	IDENT // fields, table_name

	// Misc characters
	DOT               = '.'
	DOLLAR            = '$'
	ASTERISK          = '*'
	COMMA             = ','
	LEFT_PARENTHESIS  = '('
	RIGHT_PARENTHESIS = ')'
	LEFT_BRACKET      = '['
	RIGHT_BRACKET     = ']'
	LEFT_BRACES       = '{'
	RIGHT_BRACES      = '}'
	EQUAL             = '='

	// Keywords
	INSERT = 128 + iota
	INTO
	VALUES
	UPDATE
	SET
	DELETE
	SELECT
	FROM
	WHERE
	AND
	OR
	LIKE
	NOT
	EXISTS
	GROUP
	BY
	ORDER
	HAVING
	IS
	NULL
	ASC
	DESC
	LIMIT
)

var tokens = []string{
	INSERT: "INSERT",
	INTO:   "INTO",
	VALUES: "VALUES",
	UPDATE: "UPDATE",
	SET:    "SET",
	DELETE: "DELETE",
	SELECT: "SELECT",
	FROM:   "FROM",
	WHERE:  "WHERE",
	AND:    "AND",
	OR:     "OR",
	LIKE:   "LIKE",
	NOT:    "NOT",
	EXISTS: "EXISTS",
	GROUP:  "GROUP",
	BY:     "BY",
	ORDER:  "ORDER",
	HAVING: "HAVING",
	IS:     "IS",
	NULL:   "NULL",
	ASC:    "ASC",
	DESC:   "DESC",
	LIMIT:  "LIMIT",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok, kw := range tokens {
		keywords[strings.ToUpper(kw)] = Token(tok)
	}
}
