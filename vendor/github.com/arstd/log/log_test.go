package log

import (
	"bytes"
	"fmt"
	"testing"
	"time"
)

const uuid = "6ba7b814-9dad-11d1-80b4-00c04fd430c8"

func TestLogLevel(t *testing.T) {
	l := GetLevel()
	SetLevel(Linfo)
	if IsDebugEnabled() || !IsInfoEnabled() || !IsWarnEnabled() {
		t.FailNow()
	}
	SetLevel(l) // 恢复现场，避免影响其他单元测试
}

func TestSetWriter(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 4096))
	SetWriter(buf)

	rand := time.Now().String()
	Info(rand)
	if !bytes.Contains(buf.Bytes(), ([]byte)(rand)) {
		t.FailNow()
	}
}

func TestSetFormat(t *testing.T) {
	format := fmt.Sprintf(`<log><date>%s</date><time>%s</time><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
		"2006-01-02", "15:04:05.000", LevelToken, ProjectToken, LineToken, MessageToken)
	SetFormat(format)

	buf := bytes.NewBuffer(make([]byte, 4096))
	SetWriter(buf)

	rand := time.Now().String()
	Debug(rand)
	if bytes.HasPrefix(buf.Bytes(), ([]byte)("<log><date>")) &&
		!bytes.HasSuffix(buf.Bytes(), ([]byte)("</msg><log>")) {
		t.FailNow()
	}
}

func TestPanicLog(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Fail()
		}
	}()
	Panic("test panic")
}

func TestNormalLog(t *testing.T) {
	SetLevel(Lall)

	Trace(Lall)
	Trace(Ltrace)
	Debug(Ldebug)
	Info(Linfo)
	Warn(Lwarn)
	Error(Lerror)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Panic(Lpanic)
	}()
	// Fatal( LevelFatal)
	Print(Lprint)
	Stack(Lstack)
}

func TestFormatLog(t *testing.T) {
	SetLevel(Lall)

	Tracef("%d %s", Lall, Lall)
	Tracef("%d %s", Ltrace, Ltrace)
	Debugf("%d %s", Ldebug, Ldebug)
	Infof("%d %s", Linfo, Linfo)
	Warnf("%d %s", Lwarn, Lwarn)
	Errorf("%d %s", Lerror, Lerror)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Panicf("%d %s", Lpanic, Lpanic)
	}()
	// Fatalf("%d %s", Lfatal, Lfatal)
	Printf("%d %s", Lprint, Lprint)
	Stackf("%d %s", Lstack, Lstack)
}

func TestFormatLogWithTag(t *testing.T) {
	format := "2006-01-02 15:04:05 tag info examples/main.go:88 message"
	SetFormat(format)

	SetLevel(Lall)

	Tracef("%d %s", Lall, Lall)
	Tracef("%d %s", Ltrace, Ltrace)
	Debugf("%d %s", Ldebug, Ldebug)
	Infof("%d %s", Linfo, Linfo)
	Warnf("%d %s", Lwarn, Lwarn)
	Errorf("%d %s", Lerror, Lerror)
	func() {
		defer func() {
			if err := recover(); err == nil {
				t.Fail()
			}
		}()
		Panicf("%d %s", Lpanic, Lpanic)
	}()
	// Fatalf("%d %s", Lfatal, Lfatal)
	Printf("%d %s", Lprint, Lprint)
	Stackf("%d %s", Lstack, Lstack)
}
