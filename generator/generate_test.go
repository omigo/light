package generator

import "testing"

func TestGetGomodPath(t *testing.T) {
	paths := []string{
		".",
		"/a/b/c",
		"/Users/Arstd/Reposits/projects/light/example/store/user.go",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			t.Log(getGomodPath(path))
		})
	}
}
