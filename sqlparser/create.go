package sqlparser

import "fmt"

// Parse parses a SQL CREATE statement.
func (p *Parser) ParseCreate() (*Statement, error) {
	stmt := Statement{Type: CREATE}

	// First token should be a "CREATE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != CREATE {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}
	// First token should be a "TABLE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != TABLE {
		return nil, fmt.Errorf("found %q, expected TABLE", lit)
	}

	stmt.Fragments = p.scanFragments()

	stmt.Fragments[0].Statement = "CREATE TABLE " + stmt.Fragments[0].Statement

	// Return the successfully parsed statement.
	return &stmt, nil
}
