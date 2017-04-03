light [![Build Status](https://travis-ci.org/arstd/light.svg?branch=master)](https://travis-ci.org/arstd/light)
================================================================================

Generate Go database query code, sprit from MyBatis/ibatis, GoBatis?

8 kinds of methods
--------------------------------------------------------------------------------

* insert: insert into table(name) values('name') returning id
* batch: batch insert into table(name) values('name'),('name2')
* update: update table set name='name' where id=1
* delete: delete from table where id=1
* get: select id, name from table where id=1
* count: select count(id) from table where id < 1000
* list: select id, name from table where id < 1000 order by id offset 10 limit 5
* page: select id, name from table where id < 1000 order by id offset 10 limit 5


Usage
--------------------------------------------------------------------------------

1. Code interface, Comment methods with SQL statement

```go
package persist

//go:generate light

// ModelMapper example model
type ModelMapper interface {

	// select id, name, flag, score, map, time, slice, status, pointer, struct_slice, uint32
	// from models
	// where name like ${m.Name}
	// [
	//   [ and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag}
	// ]
	// [ and slice && ${m.Slice} ]
	// [ and time between ${from} and ${to} ]
	// [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// [{ from.IsZero() && !to.IsZero() } and time <= ${to} ]
	// order by id
	// offset ${offset} limit ${limit}
	Page(m *domain.Model, ss []enum.Status, from, to time.Time, offset, limit int, tx ...*sql.Tx) (total int64, data []*domain.Model, err error)
}
```

2. Execute go generate tool

    `go generate ./...`


more example: [example/mapper/model.go](example/mapper/model.go)

generated impl code: [example/mapper/modelimpl.go](example/mapper/modelimpl.go)


More
--------------------------------------------------------------------------------

```
# light -h
usage: light [flags] [file.go]
	//go:generate light [flags] [file.go]

  -dbpath string
    	path of db to open transaction and execute SQL statements
  -dbvar string
    	variable of db to open transaction and execute SQL statements (default "db")
  -quick
    	if true, use go/types to parse dependences, much fast when built pkg cached; 
        if false, use go/loader parse source and dependences, much slow (default true)
  -skip
    	skip generate if sourceimpl.go file newer than source.go file (default true)
  -v	variable of db to open transaction and execute SQL statements

examples:
	light -force -dbvar=db.DB -dbpath=github.com/arstd/light/example/mapper
	light -force -dbvar=db2.DB -dbpath=github.com/arstd/light/example/mapper
```
