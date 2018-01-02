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
	SPACE             = ' '
	DOT               = '.'
	POUND             = '#'
	DOLLAR            = '$'
	ASTERISK          = '*'
	COMMA             = ','
	EQUAL             = '='
	LEFT_PARENTHESIS  = '('
	RIGHT_PARENTHESIS = ')'
	LEFT_BRACKET      = '['
	RIGHT_BRACKET     = ']'
	LEFT_BRACES       = '{'
	RIGHT_BRACES      = '}'
	QUESTION          = '?'
)

func isSymbol(ch rune) bool {
	switch Token(ch) {
	case DOT, POUND, DOLLAR, ASTERISK, COMMA, EQUAL,
		LEFT_PARENTHESIS, RIGHT_PARENTHESIS,
		LEFT_BRACKET, RIGHT_BRACKET,
		LEFT_BRACES, RIGHT_BRACES, QUESTION:
		return true
	default:
		return false
	}
}

const (
	// Keywords
	INSERT = 128 + iota
	INTO
	VALUES
	UPDATE
	SET
	DELETE
	CREATE
	TABLE
	IF
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
	UNION
	ALL
)

var tokens = []string{
	INSERT: "INSERT",
	INTO:   "INTO",
	VALUES: "VALUES",
	UPDATE: "UPDATE",
	SET:    "SET",
	DELETE: "DELETE",
	CREATE: "CREATE",
	TABLE:  "TABLE",
	IF:     "IF",
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
	UNION:  "UNION",
	ALL:    "ALL",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for tok, kw := range tokens {
		keywords[strings.ToUpper(kw)] = Token(tok)
	}
}
