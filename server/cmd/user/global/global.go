package global

import (
	"douyin_rpc/client/kitex_gen/user/userservice"
	common "douyin_rpc/common/config"
	"douyin_rpc/server/cmd/user/config"
	"douyin_rpc/server/cmd/video/kitex_gen/video/feedservice"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  common.NacosConfig
	Redis        *redis.Client
	UserClient   userservice.Client
	VideoClient  feedservice.Client
)
