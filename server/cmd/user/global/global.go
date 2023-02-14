package global

import (
	"douyin_rpc/client/kitex_gen/relation/relationservice"
	common "douyin_rpc/common/config"
	"douyin_rpc/server/cmd/user/config"
	"github.com/dtm-labs/rockscache"
	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	ServerConfig     config.ServerConfig
	NacosConfig      common.NacosConfig
	RocksCacheClient *rockscache.Client
	RelationClient   relationservice.Client
)
