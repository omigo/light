package sqlparser

import (
	"bytes"
	"fmt"
)

// Parse parses a SQL SELECT statement.
func (p *Parser) ParseSelect() (*Statement, error) {
	stmt := Statement{Type: SELECT}

	var buf bytes.Buffer

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	buf.WriteString("SELECT ")

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)
		buf.WriteString(lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
		buf.WriteByte(',')
	}
	buf.WriteByte(' ')

	// First token should be a "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}
	p.unscan()

	stmt.Fragments = append(stmt.Fragments, &Fragment{Statement: buf.String()})
	stmt.Fragments = append(stmt.Fragments, p.scanFragments()...)

	return &stmt, nil
}
