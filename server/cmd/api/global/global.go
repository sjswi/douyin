package global

import (
	"douyin_rpc/server/cmd/api/config"
	"douyin_rpc/server/cmd/api/kitex_gen/comment/commentservice"
	"douyin_rpc/server/cmd/api/kitex_gen/favorite/favoriteservice"
	"douyin_rpc/server/cmd/api/kitex_gen/message/messageservice"
	"douyin_rpc/server/cmd/api/kitex_gen/relation/relationservice"
	"douyin_rpc/server/cmd/api/kitex_gen/user/userservice"
	"douyin_rpc/server/cmd/api/kitex_gen/video/feedservice"
	"douyin_rpc/server/cmd/api/storage"
)

var (
	ServerConfig   config.ServerConfig
	NacosConfig    config.NacosConfig
	UserClient     userservice.Client
	MessageClient  messageservice.Client
	VideoClient    feedservice.Client
	CommentClient  commentservice.Client
	FavoriteClient favoriteservice.Client
	RelationClient relationservice.Client
	OSS            *storage.OSSClient
)
