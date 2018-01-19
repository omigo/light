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
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
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
	tok, lit = p.scanIgnoreWhitespace()
	buf.WriteString(lit)
	if tok != LPAREN {
		field = lit
		for {
			tok, lit = p.scanIgnoreWhitespace()
			if tok == COMMA || tok == eof || tok == FROM {
				p.unscan()
				return IDENT, buf.String(), field
			} else {
				buf.WriteByte(' ')
				buf.WriteString(lit)
				field = lit
			}
		}
	}

	var deep int
OUTTER:
	for {
		tok, lit = p.scan()
		buf.WriteString(lit)
		switch tok {
		case EOF:
			break

		case LPAREN:
			deep++

		case RPAREN:
			if deep == 0 {
				break OUTTER
			}
			deep--
		default:
		}
	}

	tok, lit = p.scanIgnoreWhitespace()
	if tok == AS {
		buf.WriteString(" AS ")
		tok, lit = p.scanIgnoreWhitespace()
	}

	if tok != IDENT {
		panic("require fields ")
	}
	buf.WriteString(lit)
	return IDENT, buf.String(), lit
}
