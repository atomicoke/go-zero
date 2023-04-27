```txt
syntax = "v1"

info(
    title: "靓号"
	desc: "靓号相关接口"
    curdPrefix: "QerCommon" // 生成的CURD代码前缀 可选
)

@server (
	jwt: JwtAuth
	group: users
	prefix: /admin/users/code
	middleware: PermMenuAuth
	curd: true
)
service  admin-api {}
```

```txt
@server (
    curd: true // 标记生成 curd 的路由放在这个 service 下
)
```

```shell
go run . curd -api example/apis/desc/common.api \
              -dir example/apis/admin/ \
              -url "shop:U2FsdGVkX19g@tcp(192.168.8.88:3306)/shop" \
              -table test_think_recharge_v2 > qeqwe.api
```