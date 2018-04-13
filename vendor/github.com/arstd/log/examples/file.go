package main

import (
	"fmt"

	"github.com/arstd/log"
)

func execSourceFileExamples() {
	// 全路经
	format := fmt.Sprintf("%s %s %s %s:%d %s", "2006-1-2", "3:4:05.000",
		log.LevelToken, log.PathToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 11)

	// 包
	format = fmt.Sprintf("%s %s %s %s:%d %s", "2006-1-2", "3:4:05.000",
		log.LevelToken, log.PackageToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 12)

	// 项目
	format = fmt.Sprintf("%s %s %s %s:%d %s", "2006-1-2", "3:4:05.000",
		log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 13)

	// 文件
	format = fmt.Sprintf("%s %s %s %s:%d %s", "2006-1-2", "3:4:05.000",
		log.LevelToken, log.FileToken, log.LineToken, log.MessageToken)
	log.SetFormat(format)
	log.Infof("this is a test message, %d", 14)
}
