package generator

//
// import (
// 	"bytes"
// 	"fmt"
// 	"strings"
// )
//
// func Generate0(store *Store) {
// 	var buf bytes.Buffer
// 	w := func(f string, args ...interface{}) {
// 		buf.WriteString(fmt.Sprintf(f, args...))
// 	}
// 	wln := func(f string, args ...interface{}) {
// 		buf.WriteString(fmt.Sprintf(f+"\n", args...))
// 	}
//
// 	for _, m := range store.Methods {
// 		w("func (*%sStore) %s(", store.Name, m.Name)
//
// 		for _, v := range m.Params.List() {
// 			w(v.String())
// 			w(", ")
// 		}
// 		buf.Truncate(buf.Len() - 2)
// 		w(")(")
// 		for _, v := range m.Results.List() {
// 			w(v.String())
// 			w(", ")
// 		}
// 		buf.Truncate(buf.Len() - 2)
// 		wln(") {")
// 		// body
//
// 		methodBody(m, &buf, w, wln)
//
// 		wln("}")
// 	}
//
// 	var head bytes.Buffer
// 	hln := func(f string, args ...interface{}) {
// 		head.WriteString(fmt.Sprintf(f+"\n", args...))
// 	}
//
// 	hln("package %s", store.Package)
// 	hln("import (")
// 	for k, v := range store.Imports {
// 		if v != "" {
// 			hln(`%s "%s"`, v, k)
// 		} else {
// 			hln(`"%s"`, k)
// 		}
// 	}
// 	hln(`)`)
// 	hln("type %sStore struct{}", store.Name)
//
// 	fmt.Println(head.String())
// 	fmt.Println(buf.String())
// }
//
// func methodBody(m *Method, b *bytes.Buffer, w, wln func(f string, args ...interface{})) {
// 	wln("var buf bytes.Buffer")
// 	wln("var args []interface{}")
//
// 	buf := bytes.NewBufferString(`buf.WriteString("`)
// 	args := bytes.NewBufferString("args = append(args")
//
// 	doc := m.Doc
// 	for doc != "" {
// 		i := strings.Index(doc, "${")
// 		if i == -1 {
// 			buf.WriteString(doc)
// 			break
// 		}
// 		buf.WriteString(doc[:i])
// 		doc = doc[i:]
// 		i = strings.Index(doc, "}")
// 		dollar := doc[2:i]
// 		doc = doc[i+1:]
// 		buf.WriteString("?")
// 		v := m.Params.LightByName(dollar)
//
// 		if v.Wrap() == "" {
// 			args.WriteString(", ")
// 			args.WriteString(dollar)
// 		} else {
// 			args.WriteString(", " + v.Wrap() + "(&")
// 			args.WriteString(dollar)
// 			args.WriteString(")")
// 		}
// 	}
// 	buf.WriteString(`")`)
// 	args.WriteString(`)`)
//
// 	wln(buf.String())
// 	wln(args.String())
//
// 	if strings.HasPrefix(m.Doc, "insert") {
// 		wln("query := buf.String()")
// 		wln("log.Debug(query)")
// 		wln("log.Debug(args...)")
// 		wln("	res, err := db.Exec(query, args...)")
// 		wln("if err != nil {")
// 		wln("		log.Error(query)")
// 		wln("	log.Error(args)")
// 		wln("	log.Error(err)")
// 		wln("	return 0, err")
// 		wln("}")
// 		wln("return res.LastInsertId()")
// 		return
// 	}
//
// 	panic("unimplemented")
// }
