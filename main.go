package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/omigo/light/generator"
	"github.com/omigo/light/goparser"
	"github.com/omigo/log"
	"golang.org/x/tools/imports"
)

var (
	withLog = flag.Bool("log", false, "Generated file with log")
	timeout = flag.Int64("timeout", 10, "Timeout(s) of SQL execution canceled context")
)

func main() {
	flag.Parse()

	src := getSourceFile()
	fmt.Printf("Source file    %s\n", src)
	dst := src[:len(src)-3] + ".light.go"
	// TODO must remove all *.light.go files
	os.Remove(dst)

	store, err := goparser.Parse(src, nil)
	if err != nil {
		log.Fatal(err)
	}
	// log.JSONIndent(store)
	store.Log = *withLog
	store.Timeout = *timeout

	content := generator.Generate(store)

	err = ioutil.WriteFile(dst, content, 0666)
	log.Fataln(err)
	fmt.Printf("Generated file %s\n", dst)

	pretty, err := imports.Process(dst, content, nil)
	log.Fataln(err)
	err = ioutil.WriteFile(dst, pretty, 0666)
	log.Fataln(err)
}

func getSourceFile() string {
	var src string
	if len(flag.Args()) > 0 {
		src = flag.Arg(0)
	} else {
		src = os.Getenv("GOFILE")
	}
	if src == "" {
		fmt.Println("source file must not blank")
		os.Exit(1)
	}
	if src[0] != '/' {
		wd, err := os.Getwd()
		log.Fataln(err)
		src = wd + "/" + src
	}
	return src
}
