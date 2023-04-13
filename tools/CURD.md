```api
syntax = "v1"

info(
	title: "common"
	desc: "common"
	author: "Trevor"
	email: "trevorlan@163.com"
	tableName: "12313"
	curd: "true"
)
```

1. 一个新命令, 用于生成CRUD代码, `info`中的`curd`字段为`true`时, 生成CRUD代码, 生成完后删除`curd`字段, 以免重复生成
2. 通过 `flag` 传入数据库信息
