package log

import (
	"io"
	"os"
)

func init() {
	// 默认实现标准格式标准输出
	SetPrinter(NewStandard(os.Stdout, DefaultFormatTag))
}

// Printer 定义了打印接口
type Printer interface {

	// 所有方法最终归为这个方法，真正打印日志
	Tprintf(l Level, tag string, format string, m ...interface{})

	// SetFormat 改变日志格式
	SetFormat(format string)

	// 输出日志是否着色，默认着色
	Colorized(c bool)

	// SetWriter 改变输出流
	SetWriter(w io.Writer)
}
