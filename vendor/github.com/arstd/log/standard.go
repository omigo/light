package log

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

type record struct {
	Start, End string
	Date, Time string
	Tag        string
	Level      string
	File       string
	Line       int
	Message    string
	Stack      []byte
}

// Standard 日志输出基本实现
type Standard struct {
	mu  sync.Mutex    // ensures atomic writes; protects the following fields
	out *bufio.Writer // destination for output

	format    string
	pattern   string
	colorized bool

	tpl     *template.Template
	dateFmt string
	timeFmt string
}

// NewStandard 返回标准实现
func NewStandard(w io.Writer, format string) *Standard {
	std := &Standard{out: bufio.NewWriter(w), colorized: true}

	std.SetFormat(format)
	return std
}

// SetWriter 改变输出流
func (s *Standard) SetWriter(w io.Writer) {
	s.mu.Lock()
	s.out = bufio.NewWriter(w)
	s.mu.Unlock()
}

// Colorized 输出日志是否着色，默认着色
func (s *Standard) Colorized(c bool) {
	// 没改变
	if c == s.colorized {
		return
	}

	s.colorized = c

	s.mu.Lock()
	defer s.mu.Unlock()

	p := s.pattern
	if s.colorized {
		p = "{{.Start}}" + p + "{{.End}}"
	}
	s.tpl = template.Must(template.New("record").Parse(p))
}

// SetFormat 改变日志格式
func (s *Standard) SetFormat(format string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.format = format

	s.dateFmt, s.timeFmt = ExtactDateTime(format)

	p := parseFormat(format, s.dateFmt, s.timeFmt)

	s.pattern = p
	if s.colorized {
		p = "{{.Start}}" + p + "{{.End}}"
	}
	s.tpl = template.Must(template.New("record").Parse(p))
}

// Tprintf 打印日志
func (s *Standard) Tprintf(l Level, tag string, format string, m ...interface{}) {
	if v > l {
		return
	}

	if tag == "" {
		tag = "-"
	}
	r := record{
		Level: l.String(),
		Tag:   tag,
	}

	if s.dateFmt != "" {
		now := time.Now() // get this early.
		r.Date = now.Format(s.dateFmt)
		if s.timeFmt != "" {
			r.Time = now.Format(s.timeFmt)
		}
	}

	var ok bool
	_, r.File, r.Line, ok = runtime.Caller(2) // expensive
	if ok {
		if i := strings.LastIndex(r.File, "/github.com/"); i > -1 {
			r.File = r.File[i+12:]
			if i = strings.Index(r.File, "/"); i > -1 {
				r.File = r.File[i+1:]
			}
		} else if i := strings.LastIndex(r.File, "/vendor/"); i > -1 {
			r.File = r.File[i+8:]
		} else if i := strings.LastIndex(r.File, "/src/"); i > -1 {
			r.File = r.File[i+5:]
		}
	} else {
		r.File = "???"
	}

	if format == "" {
		for _, x := range m {
			if _, ok := x.([]byte); ok {
				format += ", %s"
			} else {
				format += ", %v"
			}
		}
		if len(format) > 2 {
			format = format[2:]
		}
	}
	r.Message = fmt.Sprintf(format, m...)
	r.Message = strings.TrimSpace(r.Message)

	if l == Lstack {
		r.Stack = make([]byte, 4096)
		n := runtime.Stack(r.Stack, true)
		if n == 4096 {
			r.Stack[n-1] = '\n'
		} else {
			r.Stack[n] = '\n'
			r.Stack = r.Stack[:n+1]
		}
	}

	if s.colorized {
		r.Start, r.End = calculateColor(l)
	}

	s.mu.Lock()
	defer func() {
		s.mu.Unlock()

		if l == Lpanic {
			panic(m)
		}

		if l == Lfatal {
			os.Exit(-1)
		}
	}()

	s.tpl.Execute(s.out, r)
	s.out.WriteByte('\n')

	if l == Lstack {
		s.out.Write(r.Stack)
	}

	s.out.Flush()
}

// 格式解析，把格式串替换成 token 串
func parseFormat(format string, dateFmt, timeFmt string) (pattern string) {
	// 顺序最好不要变，从最长的开始匹配
	pattern = strings.Replace(format, PathToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, PackageToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, ProjectToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, FileToken, "{{ .File }}", -1)
	pattern = strings.Replace(pattern, TagToken, "{{ .Tag }}", -1)
	pattern = strings.Replace(pattern, LevelToken, "{{ .Level }}", -1)
	pattern = strings.Replace(pattern, strconv.Itoa(LineToken), "{{ .Line }}", -1)
	pattern = strings.Replace(pattern, MessageToken, "{{ .Message }}", -1)

	// 提取出日期和时间的格式化模式字符串
	if dateFmt != "" {
		pattern = strings.Replace(pattern, dateFmt, "{{ .Date }}", -1)
	}
	if timeFmt != "" {
		pattern = strings.Replace(pattern, timeFmt, "{{ .Time }}", -1)
	}
	return pattern
}

func calculateColor(l Level) (start, end string) {
	// all, trace, debug,          info,   warn,   error,  panic,  fatal, print, stack
	colors := []string{"", "", "", "0;32", "0;33", "0;31", "0;35", "0;35", "", ""}
	if colors[l] != "" {
		start = "\033[" + colors[l] + "m"
		end = "\033[0m"
	}
	return start, end
}

/*

Bash Shell定义文字颜色有三个参数：Style，Frontground和Background，每个参数有7个值，意义如下：

0：黑色
1：蓝色
2：绿色
3：青色
4：红色
5：洋红色
6：黄色
7：白色
其中，+30表示前景色，+40表示背景色
这里提供一段代码可以打印颜色表：

#/bin/bash
for STYLE in 0 1 2 3 4 5 6 7; do
  for FG in 30 31 32 33 34 35 36 37; do
    for BG in 40 41 42 43 44 45 46 47; do
      CTRL="\033[${STYLE};${FG};${BG}m"
      echo -en "${CTRL}"
      echo -n " ${STYLE};${FG};${BG} "
      echo -en "\033[0m"
    done
    echo
  done
  echo
done
# Reset
echo -e "\033[0m"


代码               意义
 -------------------------
 0                 OFF
 1                  高亮显示
 4                 underline
 5                  闪烁
 7                  反 白显示
 8                  不可见
*/
