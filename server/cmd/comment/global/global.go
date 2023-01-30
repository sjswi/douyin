package global

import (
	"douyin_rpc/server/cmd/user/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  config.NacosConfig
	Redis        *redis.Client
)
