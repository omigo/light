package sqlparser

import (
	"bytes"
	"testing"

	"github.com/omigo/log"
)

func TestParseInsertStmt(t *testing.T) {
	sql := "insert into users(`username`, phone, address, _status, birthday, created, updated)" + `
		values (?,?,?,?,?,CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)on duplicate key update
		  username=values(?), phone=values(?), address=values(?),
		  status=values(?), birthday=values(?), update=CURRENT_TIMESTAMP
	`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
