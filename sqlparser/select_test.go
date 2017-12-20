package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseSelectStmt(t *testing.T) {
	sql := `select *
    	from users
    	where username like ${u.Username}
    	[
			[and phone like ${u.Phone}]
			and address = ${u.Address}
		]
		and status != ${u.Status}
    	[and updated > ${u.Updated}]
    	limit ${(page-1)*size}, ${size}`

	p := Parser{s: NewScanner(bytes.NewBufferString(sql))}

	stmt, err := p.ParseSelectStmt()
	if err != nil {
		t.Fatal(err)
	}
	log.Debug(sql)
	log.JsonIndent(stmt)
}
