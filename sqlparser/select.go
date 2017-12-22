package sqlparser

import (
	"fmt"
)

// Parse parses a SQL SELECT statement.
func (p *Parser) ParseSelect() (*Statement, error) {
	stmt := Statement{Type: SELECT}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
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
