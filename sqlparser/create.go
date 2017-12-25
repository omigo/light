package sqlparser

import "fmt"

// Parse parses a SQL CREATE statement.
func (p *Parser) ParseCreate() (*Statement, error) {
	stmt := Statement{Type: CREATE}

	// First token should be a "CREATE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != CREATE {
		return nil, fmt.Errorf("found %q, expected CREATE", lit)
	}
	p.unscan()

	stmt.Fragments = p.scanFragments()

	// Return the successfully parsed statement.
	return &stmt, nil
}
