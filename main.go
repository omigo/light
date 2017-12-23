package main

import (
	"os"

	"github.com/arstd/light/generator"
	"github.com/arstd/light/goparser"
	"github.com/arstd/log"
)

func main() {
	src := os.Getenv("GOFILE")
	if src == "" {
		src = "/Users/Arstd/Reposits/src/github.com/arstd/light/example/store/user.go"
	}

	store := goparser.Parse(src)

	log.JSONIndent(store)

	generator.Generate(store)
}
