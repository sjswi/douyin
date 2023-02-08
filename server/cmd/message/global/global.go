package global

import (
	common "douyin_rpc/common/config"
	"douyin_rpc/server/cmd/user/config"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  common.NacosConfig
	Redis        *redis.Client
)
