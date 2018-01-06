package sqlparser

import (
	"bytes"
	"fmt"
	"strings"
)

// Parse parses a SQL INSERT statement.
func (p *Parser) ParseInsert() (*Statement, error) {
	stmt := Statement{Type: INSERT}
	f := Fragment{}
	stmt.Fragments = append(stmt.Fragments, &f)

	var buf bytes.Buffer
	// First token should be a "INSERT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != INSERT {
		return nil, fmt.Errorf("found %q, expected INSERT", lit)
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != INTO {
		return nil, fmt.Errorf("found %q, expected INTO", lit)
	}
	buf.WriteString("INSERT INTO ")

	// table name
	for {
		if tok, lit := p.scanIgnoreWhitespace(); tok == IDENT {
			buf.WriteString(lit)
		} else if tok == POUND {
			p.unscan()
			v := p.scanReplacer()
			f.Replacers = append(f.Replacers, v)
			buf.WriteString("%v")
		} else {
			return nil, fmt.Errorf("found %q, expected IDENT, at `%s`", lit, buf.String())
		}
		if tok, _ := p.scanIgnoreWhitespace(); tok != LPAREN {
			p.unscan()
		} else {
			buf.WriteByte('(')
			break
		}
	}

	for {
		if tok, lit := p.scanIgnoreWhitespace(); tok != IDENT {
			return nil, fmt.Errorf("found %q, expected IDENT", lit)
		} else {
			buf.WriteString(lit)
			stmt.Fields = append(stmt.Fields, lit)
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok == COMMA {
			buf.WriteByte(',')
		} else if tok == RPAREN {
			buf.WriteByte(')')
			break
		} else {
			return nil, fmt.Errorf("found %q, expected `,` or `)`", lit)
		}
	}
	if tok, lit := p.scanIgnoreWhitespace(); tok != VALUES {
		return nil, fmt.Errorf("found %q, expected `VALUES`", lit)
	}
	buf.WriteString("VALUES")
	if tok, lit := p.scanIgnoreWhitespace(); tok != LPAREN {
		return nil, fmt.Errorf("found %q, expected `(`", lit)
	}
	buf.WriteByte('(')

	// values
	for i := 0; ; i++ {
		tok, lit := p.scanIgnoreWhitespace()
		if tok == QUESTION {
			f.Variables = append(f.Variables, stmt.Fields[i])
			buf.WriteByte('?')
		} else if tok == DOLLAR {
			p.unscan()
			v := p.scanVariable()
			f.Variables = append(f.Variables, v)
			buf.WriteByte('?')
		} else {
			buf.WriteString(lit)
		}

		if tok, lit := p.scanIgnoreWhitespace(); tok == COMMA {
			buf.WriteByte(',')
		} else if tok == RPAREN {
			buf.WriteByte(')')
			break
		} else {
			return nil, fmt.Errorf("found %q, expected `,` or `)`", lit)
		}
	}

	f.Statement = strings.TrimSpace(buf.String())
	return &stmt, nil
}
