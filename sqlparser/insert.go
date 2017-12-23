package sqlparser

import "fmt"

// Parse parses a SQL INSERT statement.
func (p *Parser) ParseInsert() (*Statement, error) {
	stmt := Statement{Type: INSERT}

	// First token should be a "INSERT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INSERT {
		return nil, fmt.Errorf("found %q, expected INSEST", lit)
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}

	stmt.Fragments = p.scanFragments()

	stmt.Fragments[0].Statement = "INSEST INTO " + stmt.Fragments[0].Statement

	// Return the successfully parsed statement.
	return &stmt, nil
}
