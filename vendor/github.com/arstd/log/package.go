/*
log 实现了一个像 slf4j(Simple Logging Facade for Java)
一样标准的可以自定义级别的 log 库。这个 log 库的需要完成得任务就是提供一个标准统一的接口，同时也提供了一个基本的实现。
使用这个 log 库打印日志，可以随时切换日志级别，可以更换不同的 logger 实现，以打印不同格式的日
志，也可以改变日志输出位置，输出到数据库、消息队列等。

安装：
   go get -v -u github.com/arstd/log

使用：
    package main

    import "github.com/arstd/log"

    func main() {
        log.Debugf("this is a test message, %d", 1111)

        format := fmt.Sprintf("%s %s %s %s:%d %s", "2006-01-02 15:04:05.000000", log.TagToken,
        	log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
        log.SetFormat(format)
        log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 1111)

        format = fmt.Sprintf(`{"date": "%s", "time": "%s", "level": "%s", "file": "%s", "line": %d, "log": "%s"}`,
        	"2006-01-02", "15:04:05.999", log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
        log.SetFormat(format)
        log.Infof("this is a test message, %d", 1111)

        format = fmt.Sprintf(`<log><date>%s</date><time>%s</time><level>%s</level><file>%s</file><line>%d</line><msg>%s</msg><log>`,
        	"2006-01-02", "15:04:05.000", log.LevelToken, log.ProjectToken, log.LineToken, log.MessageToken)
        log.SetFormat(format)
        log.Tinfof("6ba7b814-9dad-11d1-80b4-00c04fd430c8", "this is a test message, %d", 1111)
    }

日志输出：

   2016-01-16 20:28:34 debug examples/main.go:10 this is a test message, 1111
   2016-01-16 20:28:34.280601 6ba7b814-9dad-11d1-80b4-00c04fd430c8 info examples/main.go:15 this is a test message, 1111
   {"date": "2016-01-16", "time": "20:28:34.28", "level": "info", "file": "examples/main.go", "line": 20, "log": "this is a test message, 1111"}
   <log><date>2016-01-16</date><time>20:28:34.280</time><level>info</level><file>examples/main.go</file><line>25</line><msg>this is a test message, 1111</msg><log>

*/
package log
