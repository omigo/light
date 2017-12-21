package sqlparser

import "fmt"

// Parse parses a SQL DELETE statement.
func (p *Parser) ParseDelete() (*Statement, error) {
	stmt := Statement{Type: DELETE}

	// First token should be a "DELETE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != DELETE {
		return nil, fmt.Errorf("found %q, expected DELETE", lit)
	}

	// First token should be a "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}

	// First token should be a "<table>" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	} else {
		stmt.Table = lit
	}

	stmt.Fragments = p.scanFragments()

	// Return the successfully parsed statement.
	return &stmt, nil
}
