package config

import "douyin_rpc/common/config"

type ServerConfig struct {
	Name         string              `mapstructure:"name" json:"name"`
	Host         string              `mapstructure:"host" json:"host"`
	MysqlInfo    config.MysqlConfig  `mapstructure:"mysql" json:"mysql"`
	OtelInfo     config.OtelConfig   `mapstructure:"otel" json:"otel"`
	RedisInfo    config.RedisConfig  `mapstructure:"redis" json:"redis"`
	UserSrvInfo  config.RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	VideoSrvInfo config.RPCSrvConfig `mapstructure:"video_srv" json:"video_srv"`
}
