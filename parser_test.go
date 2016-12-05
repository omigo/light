package main

import "testing"

func TestParseFragment(t *testing.T) {
	filename := "./testdata/methods.go"
	a := NewAnalyzer(filename)

	a.Analyze()
	// log.JSONIndent(a)

	a.parse()
}
