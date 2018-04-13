package main

import (
	g "log"

	"github.com/arstd/log"
)

func main() {
	arstdlog() //  大约 16w 行每秒

	// golog() // 大约 36.5w 行每秒
}

func arstdlog() {
	for i := 0; i < 200e4; i++ {
		log.Print("can't load package: package lib: cannot find package `xxx` in any of")
	}
}

func golog() {
	g.SetFlags(g.Ldate | g.Ltime | g.Lshortfile)
	for i := 0; i < 200e4; i++ {
		g.Print("can't load package: package lib: cannot find package `xxx` in any of")
	}
}
