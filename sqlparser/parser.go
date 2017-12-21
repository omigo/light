package sqlparser

import (
	"io"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) Parse() (s *Statement, err error) {
	tok, _ := p.scanIgnoreWhitespace()
	p.unscan()
	switch tok {
	case SELECT:
		s, err = p.ParseSelect()

	case INSERT:
		s, err = p.ParseInsert()

	case UPDATE:
		s, err = p.ParseUpdate()

	case DELETE:
		s, err = p.ParseDelete()

	default:
		panic("sql error, must start with SELECT/INSERT/UPDATE/DELETE")
	}

	if len(s.Fragments) > 0 {
		f := s.Fragments[len(s.Fragments)-1]
		f.Statement = strings.TrimSpace(f.Statement)
	}
	return s, err
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) scanVariable() (v string) {
	tok, lit := p.scanIgnoreWhitespace()
	if tok != DOLLAR {
		panic("variable must start with $")
	}
	tok, lit = p.scanIgnoreWhitespace()
	if tok != LEFT_BRACES {
		panic("variable must wraped by ${...}")
	}

	for {
		tok, lit = p.scan()
		switch tok {
		default:
			v += lit
		case WS:
			// ingnore
		case RIGHT_BRACES:
			return
		case EOF:
			panic("expect more words")
		}
	}
}

func (p *Parser) scanCond() (v string) {
	tok, lit := p.scan()
	if tok != LEFT_BRACES {
		p.unscan()
		return ""
	}
	for {
		tok, lit = p.scan()
		switch tok {
		default:
			v += lit
		case WS:
			v += " "
		case RIGHT_BRACES:
			return
		case EOF:
			panic("expect more words")
		}
	}
}
