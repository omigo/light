yan
================================================================================

根据接口和 SQL 生成数据库 CRUD 实现方法

文档 [doc/yan.slide](doc/yan.slide)（暂时未更新）

支持 6 种操作
--------------------------------------------------------------------------------

* add: insert into table(name) values('name') returning id
* modify: update table set name='name' where id=1
* remove: delete from table where id=1
* get: select id,name from table where id=1
* list: select id,name from table where id < 1000 offset 10 limit 5
* count/sum: select count(id) from table where id < 1000


Usage
--------------------------------------------------------------------------------

1. 编写接口

```go
package persist

//go:generate yan

// ModelMapper 示例接口
type ModelMapper interface {

	// select id, name, third_field, status, content
	// from demos
	// where name=${d.Name}
	//   [?{d.ThirdField != false} and third_field=${d.ThirdField} ]
	//   [?{d.Content != nil} and content=${d.Content} ]
	//   [?{len(statuses) != 0} and status in (${statuses}) ]
	// order by id
	// offset ${offset} limit ${limit}
	List(tx *sql.Tx, d *domain.Demo, statuses []enums.Status, offset, limit int) ([]*domain.Demo, error)
}
```

2. 生成代码

    `go generate ./...`

更多示例见： [example/mapper/model.go](example/mapper/model.go)

生成的文件： [example/mapper/modelimpl.go](example/mapper/modelimpl.go)


更多参数
--------------------------------------------------------------------------------

```
# yan -h
Usage of yan:
  -db string
    	variable of prefix Query/QueryRow/Exec (default "db")
  -force
    	not skip, force to rewrite impl file even if it newer than go file
  -path string
    	path variable db
  -v	version

//go:generate yan -force -db "db.DB" -path "github.com/wothing/17mei/db"
```
