package log

import (
	"encoding/json"
	"io"
)

// 默认 debug 级别，方便调试，生产环境可以调用 LevelSet 设置 log 级别
var v Level = Ldebug

// 默认实现，输出到 os.Std 中，可以重定向到文件中，也可以调用 SetPrinter 其他方式输出
var std Printer

// SetLevel 设置日志级别
func SetLevel(l Level) {
	v = l
	if v > Ldebug {
		Colorized(false)
	}
}
func SetLevelString(s string) {
	l, err := ValueOfLevel(s)
	if err != nil {
		std.Tprintf(Lerror, "", "level value string %s invalid", s)
		return
	}
	SetLevel(l)
}

// Colorized 输出日志是否着色，默认着色，如果设置的级别高于 debug，不着色
func Colorized(c bool) { std.Colorized(c) }

// GetLevel 返回设置的日志级别
func GetLevel() (l Level) { return v }

// SetPrinter 切换 Printer 实现
func SetPrinter(p Printer) { std = p }

// SetWriter 改变输出位置，通过这个接口，可以实现日志文件按时间或按大小滚动
func SetWriter(w io.Writer) { std.SetWriter(w) }

// SetFormat 改变日志格式
func SetFormat(format string) { std.SetFormat(format) }

// 判断各种级别的日志是否会被输出
func IsTraceEnabled() bool { return v <= Ltrace }
func IsDebugEnabled() bool { return v <= Ldebug }
func IsInfoEnabled() bool  { return v <= Linfo }
func IsWarnEnabled() bool  { return v <= Lwarn }
func IsErrorEnabled() bool { return v <= Lerror }
func IsPanicEnabled() bool { return v <= Lpanic }
func IsFatalEnabled() bool { return v <= Lfatal }
func IsPrintEnabled() bool { return v <= Lprint }
func IsStackEnabled() bool { return v <= Lstack }

// 打印日志
func Trace(m ...interface{}) { std.Tprintf(Ltrace, "", "", m...) }
func Debug(m ...interface{}) { std.Tprintf(Ldebug, "", "", m...) }
func Info(m ...interface{})  { std.Tprintf(Linfo, "", "", m...) }
func Warn(m ...interface{})  { std.Tprintf(Lwarn, "", "", m...) }
func Error(m ...interface{}) { std.Tprintf(Lerror, "", "", m...) }
func Panic(m ...interface{}) { std.Tprintf(Lpanic, "", "", m...) }
func Fatal(m ...interface{}) { std.Tprintf(Lfatal, "", "", m...) }
func Print(m ...interface{}) { std.Tprintf(Lprint, "", "", m...) }
func Stack(m ...interface{}) { std.Tprintf(Lstack, "", "", m...) }

var nilerr = (error)(nil)

// Errorn check last argument, if error but nil, no print log
func Errorn(m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lerror, "", "", m...)
}
func Fataln(m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lfatal, "", "", m...)
}

// 按一定格式打印日志
func Tracef(format string, m ...interface{}) { std.Tprintf(Ltrace, "", format, m...) }
func Debugf(format string, m ...interface{}) { std.Tprintf(Ldebug, "", format, m...) }
func Infof(format string, m ...interface{})  { std.Tprintf(Linfo, "", format, m...) }
func Warnf(format string, m ...interface{})  { std.Tprintf(Lwarn, "", format, m...) }
func Errorf(format string, m ...interface{}) { std.Tprintf(Lerror, "", format, m...) }
func Panicf(format string, m ...interface{}) { std.Tprintf(Lpanic, "", format, m...) }
func Fatalf(format string, m ...interface{}) { std.Tprintf(Lfatal, "", format, m...) }
func Printf(format string, m ...interface{}) { std.Tprintf(Lprint, "", format, m...) }
func Stackf(format string, m ...interface{}) { std.Tprintf(Lstack, "", format, m...) }

// Errorn check last argument, if error but nil, no print log
func Errornf(format string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lerror, "", format, m...)
}
func Fatalnf(format string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lfatal, "", format, m...)
}

// 打印日志时带上 tag
func Ttrace(tag string, m ...interface{}) { std.Tprintf(Ltrace, tag, "", m...) }
func Tdebug(tag string, m ...interface{}) { std.Tprintf(Ldebug, tag, "", m...) }
func Tinfo(tag string, m ...interface{})  { std.Tprintf(Linfo, tag, "", m...) }
func Twarn(tag string, m ...interface{})  { std.Tprintf(Lwarn, tag, "", m...) }
func Terror(tag string, m ...interface{}) { std.Tprintf(Lerror, tag, "", m...) }
func Tpanic(tag string, m ...interface{}) { std.Tprintf(Lpanic, tag, "", m...) }
func Tfatal(tag string, m ...interface{}) { std.Tprintf(Lfatal, tag, "", m...) }
func Tprint(tag string, m ...interface{}) { std.Tprintf(Lprint, tag, "", m...) }
func Tstack(tag string, m ...interface{}) { std.Tprintf(Lstack, tag, "", m...) }

// Errorn check last argument, if error but nil, no print log
func Terrorn(tag string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lerror, tag, "", m...)
}
func Tfataln(tag string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lfatal, tag, "", m...)
}

// 按一定格式打印日志，并在打印日志时带上 tag
func Ttracef(tag string, format string, m ...interface{}) {
	std.Tprintf(Ltrace, tag, format, m...)
}
func Tdebugf(tag string, format string, m ...interface{}) {
	std.Tprintf(Ldebug, tag, format, m...)
}
func Tinfof(tag string, format string, m ...interface{}) { std.Tprintf(Linfo, tag, format, m...) }
func Twarnf(tag string, format string, m ...interface{}) { std.Tprintf(Lwarn, tag, format, m...) }
func Terrorf(tag string, format string, m ...interface{}) {
	std.Tprintf(Lerror, tag, format, m...)
}
func Tpanicf(tag string, format string, m ...interface{}) {
	std.Tprintf(Lpanic, tag, format, m...)
}
func Tfatalf(tag string, format string, m ...interface{}) {
	std.Tprintf(Lfatal, tag, format, m...)
}
func Tprintf(tag string, format string, m ...interface{}) {
	std.Tprintf(Lprint, tag, format, m...)
}
func Tstackf(tag string, format string, m ...interface{}) {
	std.Tprintf(Lstack, tag, format, m...)
}

// Errorn check last argument, if error but nil, no print log
func Terrornf(tag string, format string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lerror, tag, format, m...)
}
func Tfatalnf(tag string, format string, m ...interface{}) {
	if m[len(m)-1] == nilerr {
		return
	}
	std.Tprintf(Lfatal, tag, format, m...)
}

func Json(m ...interface{}) {
	if v > Ldebug {
		return
	}
	js, err := json.Marshal(m)
	if err != nil {
		std.Tprintf(Ldebug, "", "%s", err)
	} else {
		std.Tprintf(Ldebug, "", "%s", js[1:len(js)-1])
	}
}

// 先转换成 JSON 格式，然后打印
func JSON(m ...interface{}) {
	if v > Ldebug {
		return
	}
	js, err := json.Marshal(m)
	if err != nil {
		std.Tprintf(Ldebug, "", "%s", err)
	} else {
		std.Tprintf(Ldebug, "", "%s", js[1:len(js)-1])
	}
}
func JsonIndent(m ...interface{}) {
	if v > Ldebug {
		return
	}
	js, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		std.Tprintf(Ldebug, "", "%s", err)
	} else {
		std.Tprintf(Ldebug, "", "%s", js)
	}
}

func JSONIndent(m ...interface{}) {
	if v > Ldebug {
		return
	}
	js, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		std.Tprintf(Ldebug, "", "%s", err)
	} else {
		std.Tprintf(Ldebug, "", "%s", js)
	}
}
