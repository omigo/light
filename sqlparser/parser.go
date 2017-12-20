package sqlparser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

const (
	eof      = 0
	space    = ' '
	question = '?'
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}
func isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// read reads the next rune from the bufferred reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (tok Token, lit string) {
	// Read the next rune.
	ch := s.read()

	// If we see whitespace then consume all contiguous whitespace.
	// If we see a letter then consume as an ident or reserved word.
	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isLetter(ch) {
		s.unread()
		return s.scanIdent()
	}

	// Otherwise read the individual character.
	switch Token(ch) {
	case EOF:
		return EOF, ""
	case LEFT_PARENTHESIS, LEFT_BRACKET, LEFT_BRACES,
		RIGHT_PARENTHESIS, RIGHT_BRACKET, RIGHT_BRACES,
		ASTERISK, COMMA, DOLLAR, DOT, EQUAL:
		return Token(ch), string(ch)
	default:
		return IDENT, string(ch)
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanWhitespace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return WS, buf.String()
}

// scanIdent consumes the current rune and all contiguous ident runes.
func (s *Scanner) scanIdent() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent ident character into the buffer.
	// Non-ident characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isLetter(ch) && !isDigit(ch) && ch != '_' {
			s.unread()
			break
		} else {
			_, _ = buf.WriteRune(ch)
		}
	}

	// If the string matches a keyword then return that keyword.
	ident := buf.String()
	kw := strings.ToUpper(ident)
	if tok, ok := keywords[kw]; ok {
		return tok, kw
	}

	// Otherwise return as a regular identifier.
	return IDENT, ident
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
	return &Parser{s: NewScanner(r)}
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

func (p *Parser) Parse() (Stmt, error) {
	tok, _ := p.scanIgnoreWhitespace()
	p.unscan()
	switch tok {
	case SELECT:
		return p.ParseSelectStmt()

	default:
		panic("unimpemented")
	}

}

// Parse parses a SQL SELECT statement.
func (p *Parser) ParseSelectStmt() (*SelectStmt, error) {
	stmt := SelectStmt{}

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

	stmt.Fragments = p.scanFragments()

	// Return the successfully parsed statement.
	return &stmt, nil
}

func (p *Parser) scanFragments() []*Fragment {
	var fs []*Fragment
	// scan fragment
	for {
		tok, _ := p.scanIgnoreWhitespace()
		if tok == EOF || tok == RIGHT_BRACKET {
			break
		}
		var f *Fragment
		if tok == LEFT_BRACKET {
			f = p.parseFragment(true)
		} else {
			p.unscan()
			f = p.parseFragment(false)
		}
		fs = append(fs, f)
	}
	return fs
}

func (p *Parser) parseFragment(inner bool) *Fragment {
	f := Fragment{}
	f.Cond = p.scanCond()
	if f.Cond == "" && inner {
		f.Cond = "-"
	}

	var buf bytes.Buffer
	for {
		tok, lit := p.scan()

		switch tok {
		default:
			buf.WriteString(lit)
		case WS:
			buf.WriteRune(space)

		case DOLLAR:
			p.unscan()
			lit := p.scanVariable()
			f.Variables = append(f.Variables, lit)
			buf.WriteRune(question)

		case LEFT_BRACKET:
			p.unscan()
			if inner {
				out := Fragment{Cond: f.Cond}
				f.Stmt = strings.TrimSpace(buf.String())
				if len(f.Stmt) > 0 {
					out.Fragments = append(out.Fragments, &f)
				}
				out.Fragments = append(out.Fragments, p.scanFragments()...)
				return &out
			}
			f.Stmt = strings.TrimSpace(buf.String())
			return &f

		case RIGHT_BRACKET:
			if !inner {
				p.unscan()
			}
			f.Stmt = strings.TrimSpace(buf.String())
			return &f

		case EOF:
			p.unscan()
			f.Stmt = strings.TrimSpace(buf.String())
			return &f
		}
	}

}
