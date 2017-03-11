纠结是用 go/types 库 还是 go/ast 库

go/types 取类型等稍微简单，但需要先编译一下

go/ast 能取到所有信息，但是稍微复杂一些

怎么定义一种合适数据模型，无冗余而且使用方便？？？

使用 package 结构吧

发现 go/loader 可能更方便，结合了 go/ast 和 go/types
