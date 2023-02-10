package initialize

import (
	"douyin_rpc/server/cmd/relation/global"
	"github.com/go-redis/redis/v8"
	"time"
)

func InitRedis() {

	password := global.ServerConfig.RedisInfo.Password

	addr := global.ServerConfig.RedisInfo.Address

	db := global.ServerConfig.RedisInfo.Db

	global.Redis = redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   password,
		DB:         db,
		MaxConnAge: 30 * time.Second,
	})

}
