package generator

import (
	"bytes"

	"github.com/arstd/light/goparser"
)

func writeHeader(store *goparser.Store) *bytes.Buffer {
	var header bytes.Buffer

	w := header.WriteString
	wln := func(s string) { header.WriteString(s + "\n") }

	w("package ")
	wln(store.Package)
	wln("import (")
	for k, v := range store.Imports {
		w(v)
		w(` "`)
		w(k)
		wln(`"`)
	}
	wln(")")

	w("type ")
	w(store.Name)
	wln("Store struct{}")

	return &header
}
