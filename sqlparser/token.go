// Package token defines constants representing the lexical tokens of the
// light and basic operations on tokens (printing, predicates).
//
package sqlparser

import (
	"strconv"
)

// Token is the set of lexical tokens of the Go programming language.
type Token int

// The list of tokens.
const (
	// Special tokens
	EOF Token = iota
	COMMENT

	literal_beg
	// Identifiers and basic type literals
	// (these tokens stand for classes of literals)
	IDENT  // main
	INT    // 12345
	FLOAT  // 123.45
	STRING // 'abc'
	literal_end

	VARIABLE // ${...}
	REPLACER // #{...}

	// Special Character
	POUND    // #
	DOLLAR   // $
	LBRACKET // [
	RBRACKET // ]
	LBRACES  // {
	RBRACES  // }
	QUESTION // ?

	operator_beg
	// Operator
	EQ      // =
	NE      // !=
	LG      // <>
	GT      // >
	GE      // >=
	LT      // <
	LE      // <=
	BETWEEN // between ... and ...
	operator_end

	// Misc characters
	SPACE       // SPACE
	EXCLAMATION // !
	DOT         // .
	ASTERISK    // *
	COMMA       // ,
	LPAREN      // (
	RPAREN      // )
	APOSTROPHE  // '
	BACKQUOTE   // `
	MINUS       // -

	keyword_beg
	// Keywords
	INSERT
	IGNORE
	REPLACE
	INTO
	VALUES
	UPDATE
	SET
	DELETE
	CREATE
	TABLE
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
	CURRENT_TIMESTAMP
	ON
	DUPLICATE
	KEY
	AS
	keyword_end
)

var tokens = [...]string{
	// Special tokens
	EOF:     "EOF",
	COMMENT: "COMMENT",

	// Literals
	IDENT: "IDENT", // fields, table_name

	// Special Character
	POUND:    "#",
	DOLLAR:   "$",
	LBRACKET: "[",
	RBRACKET: "]",
	LBRACES:  "{",
	RBRACES:  "}",
	QUESTION: "?",

	// Operator
	EQ:      "=",
	NE:      "!=",
	LG:      "<>",
	GT:      ">",
	GE:      ">=",
	LT:      "<",
	LE:      "<=",
	BETWEEN: "BETWEEN", // between ... and ...

	// Misc characters
	SPACE:      " ", //
	DOT:        ".",
	ASTERISK:   "*",
	COMMA:      ",",
	LPAREN:     "(",
	RPAREN:     ")",
	APOSTROPHE: "'",
	BACKQUOTE:  "`",
	MINUS:      "-",

	// Keywords
	INSERT:            "INSERT",
	IGNORE:            "IGNORE",
	REPLACE:           "REPLACE",
	INTO:              "INTO",
	VALUES:            "VALUES",
	UPDATE:            "UPDATE",
	SET:               "SET",
	DELETE:            "DELETE",
	CREATE:            "CREATE",
	TABLE:             "TABLE",
	SELECT:            "SELECT",
	FROM:              "FROM",
	WHERE:             "WHERE",
	AND:               "AND",
	OR:                "OR",
	LIKE:              "LIKE",
	NOT:               "NOT",
	EXISTS:            "EXISTS",
	GROUP:             "GROUP",
	BY:                "BY",
	ORDER:             "ORDER",
	HAVING:            "HAVING",
	IS:                "IS",
	NULL:              "NULL",
	ASC:               "ASC",
	DESC:              "DESC",
	LIMIT:             "LIMIT",
	UNION:             "UNION",
	ALL:               "ALL",
	CURRENT_TIMESTAMP: "CURRENT_TIMESTAMP",
	ON:                "ON",
	DUPLICATE:         "DUPLICATE",
	KEY:               "KEY",
	AS:                "AS",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token ADD, the string is
// "+"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
//
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := EOF; i < keyword_end; i++ {
		keywords[tokens[i]] = Token(i)
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
//
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

// Predicates

// IsLiteral returns true for tokens corresponding to identifiers
// and basic type literals; it returns false otherwise.
func (tok Token) IsLiteral() bool {
	return literal_beg < tok && tok < literal_end
}

// IsOperator returns true for tokens corresponding to operators and
// delimiters; it returns false otherwise.
func (tok Token) IsOperator() bool {
	return operator_beg < tok && tok < operator_end
}

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool {
	return keyword_beg < tok && tok < keyword_end
}
