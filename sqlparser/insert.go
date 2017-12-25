package sqlparser

import "fmt"

// Parse parses a SQL INSERT statement.
func (p *Parser) ParseInsert() (*Statement, error) {
	stmt := Statement{Type: INSERT}

	// First token should be a "INSERT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INSERT {
		return nil, fmt.Errorf("found %q, expected INSERT", lit)
	}
	p.unscan()

	stmt.Fragments = p.scanFragments()
	return &stmt, nil
}
