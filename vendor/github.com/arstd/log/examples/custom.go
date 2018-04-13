package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arstd/log"
)

func execCustomPrinterExample() {
	c := NewCustomPrinter()
	c.SetFormat(log.DefaultFormat)
	log.SetPrinter(c)

	log.Infof(msgFmt, 16)

	// 查看缓冲区
	// 每条日志结尾会加一个换行符
	line, err := c.ReadLog()
	if err != nil {
		fmt.Println(err) // 这里不能再调用 log.Error()
		os.Exit(1)
	}
	fmt.Print(line)
}

// CustomPrinter 自定义实现 Printer
type CustomPrinter struct {
	mu  sync.Mutex
	buf *bytes.Buffer

	dateFmt, timeFmt string
	tag              bool
	prefixLen        int
	line             bool
}

// NewCustomPrinter 创建 CustomPrinter
func NewCustomPrinter() *CustomPrinter {
	return &CustomPrinter{
		buf: bytes.NewBuffer(make([]byte, 4096)),
	}
}

// Tprintf 简单实现打印方法
func (p *CustomPrinter) Tprintf(l log.Level, tag string, format string, m ...interface{}) {
	if log.GetLevel() > l {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if p.dateFmt != "" {
		now := time.Now()
		p.buf.WriteString(now.Format(p.dateFmt))
		p.buf.WriteByte(' ')
		if p.timeFmt != "" {
			p.buf.WriteString(now.Format(p.timeFmt))
			p.buf.WriteByte(' ')
		}
	}

	if p.tag {
		if tag == "" {
			tag = "-"
		}
		p.buf.WriteString(tag)
		p.buf.WriteByte(' ')
	}

	if p.prefixLen > -1 {
		_, file, line, ok := runtime.Caller(2) // expensive
		if ok && p.prefixLen < len(file) {
			p.buf.WriteString(file[p.prefixLen:])
			p.buf.WriteByte(':')
			p.buf.WriteString(strconv.Itoa(line))
			p.buf.WriteByte(' ')
		} else {
			p.buf.WriteString("???:0 ")
		}
	}

	if format == "" {
		p.buf.WriteString(fmt.Sprint(m...))
	} else {
		p.buf.WriteString(fmt.Sprintf(format, m...))
	}

	p.buf.WriteByte('\n')
}

// Colorized 着色
func (p *CustomPrinter) Colorized(c bool) {}

// SetFormat 设置格式
func (p *CustomPrinter) SetFormat(format string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.dateFmt, p.timeFmt = log.ExtactDateTime(format)
	if strings.Contains(format, log.TagToken) {
		p.tag = true
	}
	p.line = strings.Contains(format, strconv.Itoa(log.LineToken))
}

// SetWriter 未实现
func (p *CustomPrinter) SetWriter(w io.Writer) {}

// ReadLog 查看日志
func (p *CustomPrinter) ReadLog() (line string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.buf.ReadString('\n')
}
