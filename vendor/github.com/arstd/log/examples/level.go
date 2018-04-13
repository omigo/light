package main

import "github.com/arstd/log"

func execLevelExamples() {
	// 默认日志级别 debug
	log.Printf("default log level: %s", log.GetLevel())
	log.Tracef("IsTraceEnabled? %t", log.IsTraceEnabled())
	log.Debugf("IsDebugEnabled? %t", log.IsDebugEnabled())
	log.Infof("IsInfoEnabled? %t", log.IsInfoEnabled())

	// trace 级别
	log.SetLevel(log.Ltrace)
	log.Tracef(msgFmt, 1)

	// info 级别
	log.SetLevel(log.Linfo)
	log.Debugf(msgFmt, 2)
	log.Infof(msgFmt, 2)

	// warn 级别
	log.SetLevel(log.Lwarn)
	log.Infof(msgFmt, 3)
	log.Warnf(msgFmt, 3)

	// error 级别
	log.SetLevel(log.Lerror)
	log.Warnf(msgFmt, 4)
	log.Errorf(msgFmt, 4)

	// 恢复默认级别，防止影响其他测试
	// debug 级别
	log.SetLevel(log.Ldebug)
	log.Tracef(msgFmt, 5)
	log.Debugf(msgFmt, 5)
}
