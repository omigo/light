package sqlparser

import "fmt"

// Parse parses a SQL UPDATE statement.
func (p *Parser) ParseUpdate() (*Statement, error) {
	stmt := Statement{}

	// First token should be a "UPDATE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != UPDATE {
		return nil, fmt.Errorf("found %q, expected UPDATE", lit)
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	} else {
		stmt.Table = lit
	}

	stmt.Fragments = p.scanFragments()

	// Return the successfully parsed statement.
	return &stmt, nil
}
