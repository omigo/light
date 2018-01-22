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
		tok, lit, field := p.scanSelectField()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}

		stmt.Fields = append(stmt.Fields, field)
		buf.WriteString(lit)

		// If the next token is not a comma then break the loop.
		if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}

		buf.WriteString(", ")
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

func (p *Parser) scanSelectField() (tok Token, lit, field string) {
	var buf bytes.Buffer
	p.s.scanSpace()

	var deep int
	for {
		tok, lit = p.scan()
		switch tok {
		case SPACE:
			buf.WriteByte(' ')

		case LPAREN:
			deep++
			buf.WriteString(LPAREN.String())

		case RPAREN:
			deep--
			buf.WriteString(RPAREN.String())

		case COMMA, FROM, EOF:
			if deep == 0 {
				p.unscan()
				return IDENT, buf.String(), field
			}
			buf.WriteString(lit)

		case BACKQUOTE:
			p.unscan()
			_, lit = p.s.scanBackQuoteIdent()
			buf.WriteString(lit)
			field = lit

		default:
			buf.WriteString(lit)
			field = lit
		}
	}

	return IDENT, buf.String(), field
}
