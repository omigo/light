package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseInsertStmt(t *testing.T) {
	sql := `
	insert into users_#{u.Username}(username, phone, address, status, birthday, created, updated)
	values (?, ${u.Phone}, ${u.Address}, ${u.Status}, ${u.Birthday},
	  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
