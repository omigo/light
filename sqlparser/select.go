package sqlparser

import (
	"bytes"
	"fmt"
)

// Parse parses a SQL SELECT statement.
func (p *Parser) ParseSelect() (*Statement, error) {
	stmt := Statement{Type: SELECT}
	first := &Fragment{}
	var buf bytes.Buffer

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	buf.WriteString("SELECT ")

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit, field, f := p.scanSelectField()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}

		stmt.Fields = append(stmt.Fields, field)
		first.Replacers = append(first.Replacers, f.Replacers...)
		first.Variables = append(first.Variables, f.Variables...)
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

	first.Statement = buf.String()
	stmt.Fragments = append(stmt.Fragments, first)
	stmt.Fragments = append(stmt.Fragments, p.scanFragments()...)

	return &stmt, nil
}

func (p *Parser) scanSelectField() (tok Token, lit, field string, f *Fragment) {
	var buf bytes.Buffer
	p.s.scanSpace()

	f = &Fragment{}

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
				return IDENT, buf.String(), field, f
			}
			buf.WriteString(lit)

		case REPLACER:
			buf.WriteString("%s")
			f.Replacers = append(f.Replacers, lit)

		case VARIABLE:
			buf.WriteString("?")
			f.Variables = append(f.Variables, lit)

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

	return IDENT, buf.String(), field, f
}
