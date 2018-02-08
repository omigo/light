package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseSelectStmt(t *testing.T) {
	sql := "select (select id from users where id=1) as id, sum(status) status,`username`, phone as phone, address, birthday, created, updated" + `
    	from users` +
		"where `from`=${from} id != -1 and username > ''" +
		`username like ?
    	[
			and address = ?
			[and phone like ${u.Phone}]
	    	and created > ${u.Created}
		]
		and status != ?
		[{range} and status in (${ss})]
    	[and updated > ${u.Updated}]
		and birthday is not null
		order by updated desc
    	limit ${page*size}, ${size}`

	p := NewParser(bytes.NewBufferString(sql))
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
