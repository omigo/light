package sqlparser

import "fmt"

// Parse parses a SQL DELETE statement.
func (p *Parser) ParseDelete() (*Statement, error) {
	stmt := Statement{Type: DELETE}

	// First token should be a "DELETE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != DELETE {
		return nil, fmt.Errorf("found %q, expected DELETE", lit)
	}
	p.unscan()

	stmt.Fragments = p.scanFragments()
	return &stmt, nil
}
