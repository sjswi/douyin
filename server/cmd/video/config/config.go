package config

import "douyin_rpc/common/config"

type ServerConfig struct {
	Name            string              `mapstructure:"name" json:"name"`
	Host            string              `mapstructure:"host" json:"host"`
	MysqlInfo       config.MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	OtelInfo        config.OtelConfig   `mapstructure:"otel" json:"otel"`
	RedisInfo       config.RedisConfig  `mapstructure:"redis" json:"redis"`
	OSSInfo         config.OSSConfig    `mapstructure:"oss" json:"oss"`
	FeedNumber      int                 `mapstructure:"feed_number" json:"feed_number"`
	UserSrvInfo     config.RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	RelationSrvInfo config.RPCSrvConfig `mapstructure:"relation_srv" json:"relation_srv"`
}
