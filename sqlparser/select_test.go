package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseSelectStmt(t *testing.T) {
	sql := `select id, username, phone, address, status, birthday, created, updated
    	from users
    	where username like ${u.Username}
    	[
			and address = ${u.Address}
			[and phone like ${u.Phone}]
	    	and created > ${u.Created}
		]
		and status != ${u.Status}
    	[and updated > ${u.Updated}]
		and birthday is not null
		order by updated desc
    	limit ${(page-1)*size}, ${size}`

	p := Parser{s: NewScanner(bytes.NewBufferString(sql))}
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
