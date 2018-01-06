package sqlparser

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

const eof = 0

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
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

// NeSPACEcanner returns a new instance of Scanner.
func NeSPACEcanner(r io.Reader) *Scanner {
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

	switch {
	case isSpace(ch):
		s.unread()
		return s.scanSpace()

	case isLetter(ch):
		s.unread()
		return s.scanIdent()

	default:
	}

	switch ch {
	case eof:
		return EOF, ""

	case '#', '$', '[', ']', '{', '}', '?', '=', '.', '*', ',', '(', ')':
		return Lookup(string(ch)), string(ch)

	case '!':
		next := s.read()
		if next == '=' {
			return NE, "!="
		}
		s.unread()
		return EXCLAMATION, "!"

	case '<':
		next := s.read()
		if next == '=' {
			return LE, "<="
		} else if next == '>' {
			return NE, "<>"
		}
		s.unread()
		return LT, "<"

	case '>':
		next := s.read()
		if next == '=' {
			return GE, ">="
		}
		s.unread()
		return GT, ">"

	default:
		return IDENT, string(ch)
	}
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanSpace() (tok Token, lit string) {
	// Create a buffer and read the current character into it.
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	// Read every subsequent whitespace character into the buffer.
	// Non-whitespace characters and EOF will cause the loop to exit.
	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isSpace(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return SPACE, buf.String()
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
