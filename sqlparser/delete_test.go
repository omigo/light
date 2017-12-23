package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseDeleteStmt(t *testing.T) {
	sql := `DELETE FROM users WHERE id=${id}`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
