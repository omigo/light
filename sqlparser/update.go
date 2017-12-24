package sqlparser

import "fmt"

// Parse parses a SQL UPDATE statement.
func (p *Parser) ParseUpdate() (*Statement, error) {
	stmt := Statement{Type: UPDATE}

	// First token should be a "UPDATE" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != UPDATE {
		return nil, fmt.Errorf("found %q, expected UPDATE", lit)
	}
	p.unscan()

	stmt.Fragments = p.scanFragments()

	stmt.Fragments[0].Statement = "UPDATE " + stmt.Fragments[0].Statement
	// Return the successfully parsed statement.
	return &stmt, nil
}
