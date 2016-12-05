package main

import (
	"testing"

	"github.com/arstd/log"
)

func TestParsePackage(t *testing.T) {
	filename := "./testdata/package.go"
	a := NewAnalyzer(filename)

	a.Analyze()
	log.JSONIndent(a)
}

func TestParseImports(t *testing.T) {
	filename := "./testdata/imports.go"
	a := NewAnalyzer(filename)

	a.Analyze()
	log.JSONIndent(a)
}

func TestParseInterfaces(t *testing.T) {
	filename := "./testdata/interfaces.go"
	a := NewAnalyzer(filename)

	a.Analyze()
	log.JSONIndent(a)
}

func TestParseMethods(t *testing.T) {
	filename := "./testdata/methods.go"
	a := NewAnalyzer(filename)

	a.Analyze()
	log.JSONIndent(a)
}
