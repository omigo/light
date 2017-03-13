package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arstd/light/parse"
	"github.com/arstd/light/prepare"
	"github.com/arstd/log"
)

func init() {
	// log.SetLevel(log.Lwarn)
	log.SetFormat("2006-01-02 15:04:05.999 info examples/main.go:88 message")

	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "usage: light [flags] [file.go]\n\t//go:generate light [flags] [file.go]")
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()
}

func main() {
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

	p := parse.ParseGoFile(goFile)

	prepare.PrepareStmt(p)

	log.JSONIndent(p)
}
