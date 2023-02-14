package initialize

import (
	"douyin_rpc/server/cmd/relation/global"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {

	password := global.ServerConfig.RedisInfo.Password

	addr := global.ServerConfig.RedisInfo.Address

	db := global.ServerConfig.RedisInfo.Db

	redisClient := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    []string{addr},
		DB:       db,
		Password: password,
	})
	global.RocksCacheClient = rockscache.NewClient(redisClient, rockscache.NewDefaultOptions())

}
