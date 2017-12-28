package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/arstd/light/generator"
	"github.com/arstd/light/goparser"
	"github.com/arstd/log"
	"golang.org/x/tools/imports"
)

func main() {
	src := os.Getenv("GOFILE")
	if src == "" {
		src = "/Users/Arstd/Reposits/src/github.com/arstd/light/example/store/user.go"
	}
	if src[0] != '/' {
		wd, err := os.Getwd()
		log.Errorn(err)
		src = wd + "/" + src
	}
	fmt.Printf("Source file    %s\n", src)
	dst := src[:len(src)-3] + ".light.go"
	os.Remove(dst)

	store := goparser.Parse(src)
	// log.JSONIndent(store)

	content := generator.Generate(store)

	err := ioutil.WriteFile(dst, content, 0666)
	log.Fataln(err)
	fmt.Printf("Generated file %s\n", dst)

	pretty, err := imports.Process(dst, content, nil)
	log.Fataln(err)
	err = ioutil.WriteFile(dst, pretty, 0666)
	log.Fataln(err)
}
