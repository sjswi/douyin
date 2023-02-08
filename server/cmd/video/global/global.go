package global

import (
	"douyin_rpc/client/kitex_gen/comment/commentservice"
	"douyin_rpc/client/kitex_gen/favorite/favoriteservice"
	"douyin_rpc/client/kitex_gen/relation/relationservice"
	"douyin_rpc/client/kitex_gen/user/userservice"
	common "douyin_rpc/common/config"
	"douyin_rpc/server/cmd/video/config"
	"douyin_rpc/server/cmd/video/kitex_gen/video/feedservice"
	"douyin_rpc/server/cmd/video/storage"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	DB             *gorm.DB
	ServerConfig   config.ServerConfig
	NacosConfig    common.NacosConfig
	Redis          *redis.Client
	UserClient     userservice.Client
	VideoClient    feedservice.Client
	FavoriteClient favoriteservice.Client
	CommentClient  commentservice.Client
	RelationClient relationservice.Client
	OSS            *storage.OSSClient
)
