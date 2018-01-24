package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arstd/light/generator"
	"github.com/arstd/light/goparser"
	"github.com/arstd/log"
	"golang.org/x/tools/imports"
)

var (
	lg = flag.Bool("log", false, "Generated file with log")
)

// 每次执行 `go generate ./...` 前，先编译安装，保证代码最新
//go:generate go install

func main() {
	flag.Parse()

	src := getSourceFile()
	fmt.Printf("Source file    %s\n", src)
	dst := src[:len(src)-3] + ".light.go"
	os.Remove(dst)

	store := goparser.Parse(src)
	// log.JSONIndent(store)
	store.Log = *lg

	content := generator.Generate(store)

	err := ioutil.WriteFile(dst, content, 0666)
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
