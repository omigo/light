package main

import (
	"testing"

	"github.com/arstd/log"
)

func TestPrepareFragment(t *testing.T) {
	filename := "./testdata/methods.go"
	a := NewAnalyzer(filename)

	a.Analyze()

	a.parse()

	log.JSONIndent(a)
}
