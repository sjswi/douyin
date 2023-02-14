package global

import (
	"douyin_rpc/client/kitex_gen/comment/commentservice"
	"douyin_rpc/client/kitex_gen/relation/relationservice"
	"douyin_rpc/client/kitex_gen/user/userservice"
	"douyin_rpc/client/kitex_gen/video/feedservice"
	common "douyin_rpc/common/config"
	"douyin_rpc/server/cmd/favorite/config"
	"github.com/dtm-labs/rockscache"
	"gorm.io/gorm"
)

var (
	DB               *gorm.DB
	ServerConfig     config.ServerConfig
	NacosConfig      common.NacosConfig
	RocksCacheClient *rockscache.Client
	UserClient       userservice.Client
	VideoClient      feedservice.Client
	RelationClient   relationservice.Client
	CommentClient    commentservice.Client
)
