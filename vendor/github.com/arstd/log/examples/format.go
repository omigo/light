package main

import (
	"fmt"

	"github.com/arstd/log"
)

func execFormatExamples() {
	// 默认简洁格式
	log.Infof("this is a test message, %d", 6)

	// 带标签的格式
	log.SetFormat(log.DefaultFormatTag)
	log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 7)

	// 自定义其他格式的日志
	format := fmt.Sprintf("%s %s %s %s:%d %s", "2006-1-2", "3:4:05.000",
		log.LevelToken, log.PathToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 8)

	// 自定义 json 格式的日志
	format = fmt.Sprintf(`{"date": "%s", "time": "%s", "level": "%s", "file": "%s", "line": %d, "log": "%s"}`,
		"2006-01-02", "15:04:05.999", log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 9)

	// 自定义 xml 格式的日志
	format = fmt.Sprintf(`<log><date>%s</date><time>%s</time><tid>%s</tid><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
		"2006-01-02", "15:04:05.000", log.TagToken, log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 10)
}
