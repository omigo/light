package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/arstd/log"
)

func execChangeWriterExample() {
	// 改变 Writer，把日志打印到缓冲区，也可以打印到文件，定时切换文件，可以实现日志滚动
	buf := bytes.NewBuffer(make([]byte, 255))
	log.SetWriter(buf)

	log.Infof(msgFmt, 15)

	// 查看缓冲区
	// 每条日志结尾会加一个换行符
	line, err := buf.ReadString('\n')
	if err != nil {
		fmt.Println(err) // 这里不能再调用 log.Error()
		os.Exit(1)
	}

	fmt.Print(line)
}
