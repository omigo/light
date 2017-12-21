package sqlparser

import (
	"bytes"
	"fmt"
	"strings"
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

func (p *Parser) scanFragments() (fs []*Fragment) {
	// scan fragment
	for {
		f, lastToken := p.parseFragment()
		if f != nil {
			fs = append(fs, f)
		}
		if lastToken == EOF {
			break
		}
	}
	return fs
}

func (p *Parser) parseFragment() (*Fragment, Token) {
	var inner bool
	var buf bytes.Buffer

	tok, lit := p.scanIgnoreWhitespace()
	if tok == LEFT_BRACKET {
		inner = true
	} else if tok == RIGHT_BRACKET {
		p.unscan()
		return nil, EOF
	} else if tok == ORDER {
		buf.WriteString(strings.ToUpper(lit))
	} else {
		p.unscan()
	}

	f := Fragment{}
	f.Condition = p.scanCond()
	if f.Condition == "" && inner {
		f.Condition = "-"
	}

	for {
		tok, lit = p.scan()

		switch tok {
		default:
			buf.WriteString(lit)

		case WS:
			buf.WriteRune(space)

		case DOLLAR:
			p.unscan()
			lit = p.scanVariable()
			f.Variables = append(f.Variables, lit)
			buf.WriteRune(question)

		case LEFT_BRACKET:
			p.unscan()
			if inner {
				stmt := strings.TrimSpace(buf.String())
				buf.Reset()
				if len(stmt) > 0 {
					innerFirst := Fragment{Statement: stmt, Variables: f.Variables}
					f.Variables = nil
					f.Fragments = append(f.Fragments, &innerFirst)
				}
				f.Fragments = append(f.Fragments, p.scanFragments()...)
			}
			goto END

		case RIGHT_BRACKET, ORDER, EOF:
			p.unscan()
			goto END
		}
	}

END:
	tok, lit = p.scanIgnoreWhitespace()
	if inner {
		if tok != RIGHT_BRACKET {
			panic("expect ], but got " + lit + ", " + buf.String())
		}
	} else {
		p.unscan()
		if tok == RIGHT_BRACKET {
			tok = EOF
		}
	}
	f.Statement = strings.TrimSpace(buf.String())
	if len(f.Statement) > 0 {
		f.Statement += " "
	}
	return &f, tok
}
