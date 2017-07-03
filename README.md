light [![Build Status](https://travis-ci.org/arstd/light.svg?branch=master)](https://travis-ci.org/arstd/light)
================================================================================

`light` 是一个代码生成工具，帮助生成 DB 层代码。受 MyBatis/iBatis 启发，`light` 做到了最简，
只需要一个带参数的原始 SQL 语句和一个接口方法，无需其他代码。但和 MyBatis 不同的是，`light`
是编译前预先生成代码绑定参数和返回值，而不是运行中通过反射来绑定。`light` 是工具，生成代码，发生
在编译前，它不是 ORM 库。


支持各种 SQL 操作
--------------------------------------------------------------------------------

* 插入: insert into table(name) values('name') returning id
* 批量插入: insert into table(name) values('name'),('name2')
* 更新: update table set name='name' where id=1
* 删除: delete from table where id=1
* 单条查询: select id, name, (select 2) as other from table where id=1
* 统计: select count(id) from table where id < 1000 // count 专用，sum 等操作可以使用单条查询
* 列表查询: select id, name from table where id < 1000 order by id offset 10 limit 5
* 分页查询: select id, name from table where id < 1000 order by id offset 10 limit 5


用法
--------------------------------------------------------------------------------

0. 安装 `light`

```
go get -u -v github.com/arstd/light

# or

git clone git@github.com:arstd/light.git $GOPATH/src/github.com/arstd/light
cd $GOPATH/src/github.com/arstd/light
glide install # glide 安装依赖，无需翻墙
go install -v # 生成可执行程序到 $GOPATH/bin（需确保这个目录在 $PATH 中）
```

1. 按规范编写接口，并在方法上写 SQL 语句 [example/mapper/model.go](example/mapper/model.go)

```go
package persist

//go:generate light

// ModelMapper example model
type ModelMapper interface {

	// select id, name, flag, score, map, time, xarray, slice, status, pointer, struct_slice, uint32
	// from models
	// where name like ${m.Name}
	// [ [ and status in ( [{range ss}] ) ]
	//   and flag=${m.Flag}
	//   [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// ]
	// [ and time between ${from} and ${to} ]
	// [{ !from.IsZero() && to.IsZero() } and time >= ${from} ]
	// [{ from.IsZero() && !to.IsZero() } and time <= ${to} ]
	// [ and xarray && array[ [{range m.Array}] ] ]
	// [ and slice && ${m.Slice} ]
	// order by #{orderBy}
	// offset ${offset} limit ${limit}
	Page(m *domain.Model, sortBy string, ss []enum.Status, from, to time.Time, offset, limit int, tx ...*sql.Tx) (total int64, data []*domain.Model, err error)
}
```

2. 执行命令生成代码

    `go generate ./...`


more example: [example/mapper/model.go](example/mapper/model.go)

generated impl code: [example/mapper/modelimpl.go](example/mapper/modelimpl.go)


更多
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
	light -force -dbvar=alias.DB -dbpath=github.com/arstd/light/example/mapper
```


回答几个问题
--------------------------------------------------------------------------------

> 1. 为什么不用 ORM 库？

答：Github 上的确有很多相当好的开源 ORM，但是我个人不习惯使用：

（1）大部分 ORM 都提供类似 `Select` `Find` `SortBy` `Where` 等这样的方法，使用时要先把写
好的 SQL 语句强制拆解，转化成这些方法的链式调用，虽然这个转化非常简单，但也增加了心智负担，为了支持 SQL
各种条件和操作符，需要定义一堆方法和结构体。也有 ORM 可以
直接传入完整 SQL 语句的，但大多只是辅助，相对简单，复杂条件不能处理；

（2）ORM 大量使用反射，效率可能会降低两个数量级（网传，未确认）；

（3）ORM 错误只在运行时才抛出，`light` 生成代码，错误是很容易在编译前发现的。

实时上 database/sql 已经高度抽象了，再封装也是减少几行代码而已，并且付出的代价也是不小的，可能
得不尝失。如果说 SQL 是第四代编程语言， C/C++/Java/Go 是第三代编程语言，那么我们为什么要把更
高级的语言向下转化来编程呢？而且 SQL 标准和通用性超过任何一门其他语言，所以应当尽可能保留 SQL
原生的结构。

> 2. 为什么不直接连接到数据库生成代码？

答：对生成代码来说，需要知道的仅仅是字段和类型而已，结构体本身的字段和类型已经足够了，连接数据库
反而增加了开发和使用这个工具的复杂度。保持简单，灵活，无侵入，无依赖。很多
数据库查询工具也可以帮助导出 SQL 语句和字段。我写了个程序从建表语句生成符合 `light` 规范的 `insert/select/update`，代码不到 30 行。

> 3. 不支持 MySQL ？

答：从 2016 年 3 月开始写，5 月第一个可用版本出来，直到现在完全推翻重写了 4 次，但依然很乱。设计
上没有达到自己的理想，所以代码结构上也没优化，基本上都是一些上百行的大函数，测试没有，文档也没有。所以没
有打算推广，一直都是自己在用。暂时只支持 PostgreSQL。写一个 MySQL 的模板，改动不
大，如果确实需要使用，可以发邮件给我，有人使用，我会尽力去完善的。
