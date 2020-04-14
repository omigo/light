package sqlparser

import (
	"bytes"
	"testing"

	"github.com/omigo/log"
)

func TestParseUpdateStmt(t *testing.T) {
	sql := `UPDATE users
	SET [username=${u.Username},]
	    [phone=${u.Phone},]
	    [address=${u.Address},]
	    [status=${u.Status},]
	    [birthday=${u.Birthday},]
	    updated=CURRENT_TIMESTAMP
	WHERE id=${u.Id}`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
