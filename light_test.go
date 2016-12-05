package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	tests := []testing.InternalTest{
		// {Name: "TestParsePackage", F: TestParsePackage},
		// {Name: "TestParseImports", F: TestParseImports},
		// {Name: "TestParseInterfaces", F: TestParseInterfaces},
		// {Name: "TestParseMethods", F: TestParseMethods},
		{Name: "TestParseFragment", F: TestParseFragment},
	}

	var run = func(pat string, str string) (bool, error) {
		return true, nil
	}
	m = testing.MainStart(run, tests, nil, nil)
	// setup()
	ret := m.Run()
	// teardown()
	os.Exit(ret)
}
