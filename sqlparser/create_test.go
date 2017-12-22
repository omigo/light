package sqlparser

import (
	"bytes"
	"testing"

	"github.com/arstd/log"
)

func TestParseCreate(t *testing.T) {
	sql := `create table if not exists #{dev.Platform}_#{dev.Cid} (
			cid text, platform text, version text
		) ENGINE=InnoDB DEFAULT CHARSET=utf8 `

	p := Parser{s: NewScanner(bytes.NewBufferString(sql))}
	stmt, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}
	log.JsonIndent(stmt)
}
