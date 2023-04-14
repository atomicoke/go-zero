package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/tools/goctl/config"
	"github.com/zeromicro/go-zero/tools/goctl/example/out/internal/model"
)

type ServiceContext struct {
	Version          int
	Redis            *redis.Redis
	PermMenuAuth     rest.Middleware
	DeleteVerifyCode rest.Middleware
	LoginLog         rest.Middleware
	UsersModel       model.ThinkUsersModel
	AppLiveModel        model.AppLiveModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	mysqlConn := sqlx.NewMysql("")
	redisClient := redis.MustNewRedis(redis.RedisConf{
		Host: "",
		Type: "",
		Pass: "",
	})

	var s = &ServiceContext{
		Redis:      redisClient,
		UsersModel: model.NewThinkUsersModel(mysqlConn, c.Cache),
		AppLiveModel: model.NewAppLiveModel(mysqlConn, c.Cache),
	}

	return s
}
