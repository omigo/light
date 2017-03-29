package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"golang.org/x/tools/imports"

	"github.com/arstd/log"
)

func usage() {
	fmt.Fprintln(os.Stderr, `usage: light [flags] [file.go]
	//go:generate light [flags] [file.go]
`)
	flag.PrintDefaults()

	fmt.Fprintln(os.Stderr, `
examples:
	light -force -dbvar=db.DB -dbpath=github.com/arstd/light/example/mapper
	light -force -dbvar=db2.DB -dbpath=github.com/arstd/light/example/mapper
`)
	os.Exit(2)
}

func main() {
	log.SetLevel(log.Linfo)
	log.SetFormat("2006-01-02 15:04:05.999 info examples/main.go:88 message")

	dbVar := flag.String("dbvar", "db", "variable of db to open transaction and execute SQL statements")
	dbPath := flag.String("dbpath", "", "path of db to open transaction and execute SQL statements")
	force := flag.Bool("force", false, "force to regenerate, even sourceimpl.go file newer than source.go file")
	version := flag.Bool("v", false, "variable of db to open transaction and execute SQL statements")
	flag.Usage = usage

	flag.Parse()
	if *version {
		fmt.Println("0.5.5")
		return
	}

	goFile := os.Getenv("GOFILE")
	if goFile == "" {
		if flag.NArg() > 1 {
			goFile = flag.Arg(0)
			if !strings.HasSuffix(goFile, ".go") {
				fmt.Println("file suffix must match *.go")
				return
			}
		} else {
			flag.Usage()
		}
	}
	fmt.Printf("Found  go file: %s\n", goFile)

	outFile := goFile[:len(goFile)-3] + "impl.go"
	if !*force {
		outStat, err := os.Stat(outFile)
		if err != nil {
			// log.Info(err)
		} else {
			goStat, _ := os.Stat(goFile)
			if !outStat.ModTime().Before(goStat.ModTime()) {
				fmt.Printf("Generated file: %s, skip!\n", outFile)
				return
			}
		}
	}

	pkg := &Package{
		Source:  goFile,
		DBVar:   *dbVar,
		Imports: map[string]string{},
	}
	if *dbPath != "" {
		ss := strings.Split(*dbVar, ".")
		if len(ss) != 2 {
			fmt.Println("arg 'dbvar' must be <package-name>:<variable-name")
			flag.Usage()
			return
		}
		pkg.Imports[ss[0]] = strings.Trim(*dbPath, `'"`)
	}

	ParseGoFile(pkg)
	Prepare(pkg)
	log.JSONIndent(pkg)

	paths := strings.Split(os.Getenv("GOPATH"), string(filepath.ListSeparator))
	tmplFile := filepath.Join(paths[0], "src", "github.com/arstd/light", "postgresql.pq.gotemplate")

	funcMap := template.FuncMap{
		"timestamp": func() string { return time.Now().Format("2006-01-02 15:04:05") },
	}

	tmpl, err := template.New("postgresql.pq.gotemplate").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		log.Panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pkg)
	if err != nil {
		log.Panic(err)
	}

	ioutil.WriteFile(outFile, buf.Bytes(), 0644)

	pretty, err := imports.Process(outFile, buf.Bytes(), nil)
	if err != nil {
		log.Panic(err)
	}
	err = ioutil.WriteFile(outFile, pretty, 0644)
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("Generated file: %s\n", outFile)
}
