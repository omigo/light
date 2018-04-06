package sqlparser

import (
	"bytes"
	"io"
	"strings"
)

func Parse(doc string) (s *Statement, err error) {
	return NewParser(bytes.NewBufferString(doc)).Parse()
}

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
	return &Parser{s: NeSPACEcanner(r)}
}

func (p *Parser) Parse() (s *Statement, err error) {
	tok, _ := p.scanIgnoreWhitespace()
	p.unscan()
	switch tok {
	case SELECT:
		s, err = p.ParseSelect()

	case INSERT:
		s, err = p.ParseInsert()

	case REPLACE:
		s, err = p.ParseReplace()

	case UPDATE:
		s, err = p.ParseUpdate()

	case DELETE:
		s, err = p.ParseDelete()

	case CREATE:
		s, err = p.ParseCreate()

	default:
		panic("sql error, must start with SELECT/INSERT/UPDATE/DELETE")
	}
	if err != nil {
		return nil, err
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
	if tok == SPACE {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) scanVariable() (v string) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok != DOLLAR {
		panic("variable must start with $")
	}
	tok, _ = p.scanIgnoreWhitespace()
	if tok != LBRACES {
		panic("variable must wraped by ${...}")
	}

	var lit string
	for {
		tok, lit = p.scan()
		switch tok {
		default:
			v += lit
		case SPACE:
			// ingnore
		case RBRACES:
			return
		case EOF:
			panic("expect more words")
		}
	}
}

func (p *Parser) scanReplacer() (v string) {
	tok, _ := p.scanIgnoreWhitespace()
	if tok != POUND {
		panic("replacer must start with #")
	}
	tok, _ = p.scanIgnoreWhitespace()
	if tok != LBRACES {
		panic("replacer must wraped by #{...}")
	}

	var lit string
	for {
		tok, lit = p.scan()
		switch tok {
		default:
			v += lit
		case SPACE:
			// ingnore
		case RBRACES:
			return
		case EOF:
			panic("expect more words")
		}
	}
}

func (p *Parser) scanCondition() (v string) {
	tok, _ := p.scan()
	if tok != LBRACES {
		p.unscan()
		return ""
	}

	var buf bytes.Buffer
	for {
		tok, lit := p.scan()
		switch tok {
		default:
			buf.WriteString(lit)
		case SPACE:
			buf.WriteString(" ")
		case RBRACES:
			return buf.String()
		case EOF:
			panic("expect more words")
		}
	}
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
	if tok == LBRACKET {
		inner = true
	} else if tok == RBRACKET {
		p.unscan()
		return nil, EOF
	} else if tok == ORDER {
		buf.WriteString(strings.ToUpper(lit))
	} else {
		p.unscan()
	}

	f := Fragment{}
	f.Condition = p.scanCondition()
	if f.Condition == "" && inner {
		f.Condition = "-"
	}

	var last string
	for {
		tok, lit = p.scan()
		switch tok {
		default:
			buf.WriteString(lit)

		case IDENT:
			buf.WriteString(lit)
			last = lit

		case SPACE:
			buf.WriteString(SPACE.String())

		case QUESTION:
			f.Variables = append(f.Variables, last)
			buf.WriteString(QUESTION.String())

		// case DOLLAR:
		// 	p.unscan()
		// 	lit = p.scanVariable()
		// 	f.Variables = append(f.Variables, lit)
		// 	buf.WriteString(QUESTION.String())
		case VARIABLE:
			f.Variables = append(f.Variables, lit)
			buf.WriteString(QUESTION.String())

		// case POUND:
		// 	p.unscan()
		// 	lit = p.scanReplacer()
		// 	f.Replacers = append(f.Replacers, lit)
		// 	buf.WriteString("%v")
		case REPLACER:
			f.Replacers = append(f.Replacers, lit)
			buf.WriteString("%v")

		case LBRACKET:
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

		case ORDER, LIMIT:
			if inner {
				buf.WriteString(lit)
				last = lit
			} else {
				p.unscan()
				goto END
			}

		case RBRACKET, EOF:
			p.unscan()
			goto END
		}
	}

END:
	tok, lit = p.scanIgnoreWhitespace()
	if inner {
		if tok != RBRACKET {
			panic("expect ], but got " + lit + ", " + buf.String())
		}
	} else {
		p.unscan()
		if tok == RBRACKET {
			tok = EOF
		}
	}
	f.Statement = strings.TrimSpace(buf.String())
	if strings.TrimSpace(f.Condition) == "range" {
		f.Condition = "-"
		f.Range = f.Replacers[0]
	}
	return &f, tok
}
