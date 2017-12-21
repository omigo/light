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
	if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	} else {
		stmt.Table = lit
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != LEFT_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected (", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != RIGHT_PARENTHESIS {
		return nil, fmt.Errorf("found %q, expected )", lit)
	}

	stmt.Fragments = p.scanFragments()

	// Return the successfully parsed statement.
	return &stmt, nil
}
