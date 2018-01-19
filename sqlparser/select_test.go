package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseSelectStmt(t *testing.T) {
	sql := "select (select id from users where id=1) as id, `username`, phone as phone, address, status, birthday, created, updated" + `
    	from users
    	where id != -1 and username > ''
		username like ?
    	[
			and address = ?
			[and phone like ${u.Phone}]
	    	and created > ${u.Created}
		]
		and status != ?
    	[and updated > ${u.Updated}]
		and birthday is not null
		order by updated desc
    	limit ${(page-1)*size}, ${size}`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
