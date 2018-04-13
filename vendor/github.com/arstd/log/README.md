log [![Build Status](https://travis-ci.org/arstd/log.svg?branch=master)](https://travis-ci.org/arstd/log) [![Godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/arstd/log) [![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/arstd/log/master/LICENSE)
================================================================================

`log` 提供一
个类型 `slf4j` 的标准接口，同时也提供了一个基本的实现，可以自定义日志格式，输出各种类型的日志，如
csv/json/xml，同时支持 Tag（TraceId/RequestId)。


Usage
-----

安装：`go get -v -u github.com/arstd/log`

使用：
``` go
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
```
日志输出：
```
2016-01-16 20:28:34 debug examples/main.go:10 this is a test message, 1111
2016-01-16 20:28:34.280601 6ba7b814-9dad-11d1-80b4-00c04fd430c8 info examples/main.go:15 this is a test message, 1111
{"date": "2016-01-16", "time": "20:28:34.28", "level": "info", "file": "examples/main.go", "line": 20, "log": "this is a test message, 1111"}
<log><date>2016-01-16</date><time>20:28:34.280</time><level>info</level><file>examples/main.go</file><line>25</line><msg>this is a test message, 1111</msg><log>
```

更多用法 [examples](examples/main.go)
着色示例 ![color.png](color.png)

Go Doc and API
--------------

所有可调用的接口 API 和 文档都在 [log.go](log.go)


log/Printer/Standard
--------------------

Golang 不同于 Java，非面向对象语言（没有继承，只有组合，不能把组合实例赋给被组合的实例，即 Java
说的 子对象 赋给 父对象），为了方便使用，很多函数都是包封装的，无需创建 struct ，就可以直接调用。
（一般把裸露的方法称为函数，结构体和其他类型的方法才称为某某的方法）

log 包也一样，使用时，无需 new ，直接用。log 包有所有级别的函数可以调用，所有函数最终都调用了
print 函数。print 函数又调用了包内部变量的 std 的 Print 方法。这个 std 是一个 Printer 接
口类型，定义了打印接口。用不同的实现改变 std 就可以打印出不同格式的日志，也可以输出到不同位置。
（这个接口貌似还没有抽象好，再想想）

Printer 有个基本的实现 Standard，如果不改变，默认使用这个实现打印日志。

Standard 实现了的 Printer 接口，把日志打印到 Stdout。


性能测试
-------

环境：MacBookPro 15，4 核 8 线程 16G 内存

实际测试结果（把日志重定向到文件）：
该库平均每秒可输出 16w 行日志；
Go 语言标准库平均每秒输出 36.5w 行日志。

模板方式输出日志对性能有一定影响，其他 New Record 等也可能造成性能下降。但是实现上比标准库略微
复杂，输出格式可以灵活配置，所以整体上可以接受，后期随着对 Go 语言的学习更深入再不断优化。


TODO
----

* 测试是否支持各种格式的日期
* 处理秒和毫秒，如1:1:02.9
* 实现日志文件按一定规则自动滚动
