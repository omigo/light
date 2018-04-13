package test

import (
	"bytes"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/arstd/log"
)

func TestDefaultFormat(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 4096))
	log.SetWriter(buf)

	msg := "this is a test message"
	log.Info(msg)
	if ok, _ := regexp.Match(`\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}.\d+ -? ?info .*test/deeper/log_test.go:\d+ this is a test message`, buf.Bytes()); !ok {
		t.Logf("%s", buf.Bytes()) // 2016-01-24 19:41:19 info test/deeper/log_test.go:16 this is a test message
		t.FailNow()
	}
}

func TestSetFormatFile(t *testing.T) {
	format := fmt.Sprintf(`<log><date>%s</date><time>%s</time><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
		"2006-01-02", "15:04:05.000", log.LevelToken, log.PackageToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)

	buf := bytes.NewBuffer(make([]byte, 4096))
	log.SetWriter(buf)

	rand := time.Now().String()
	log.Debug(rand)
	if bytes.HasPrefix(buf.Bytes(), ([]byte)("<file>github.com/arstd/log/test/deeper/log_test.go</file>")) {
		t.FailNow()
	}
}
