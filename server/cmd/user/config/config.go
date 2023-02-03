package config

import "douyin_rpc/common/config"

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
	MaxConn  int    `mapstructure:"max_conn" json:"max_conn"`
	MaxIdle  int    `mapstructure:"max_idle" json:"max_idle"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}
type RedisConfig struct {
	Address    string `mapstructure:"address" json:"address"`
	Db         int    `mapstructure:"db" json:"db"`
	Password   string `mapstructure:"password" json:"password"`
	MaxConnAge string `mapstructure:"max_conn_age" json:"max_conn_age"`
}
type WXConfig struct {
	AppId     string `mapstructure:"app_id" json:"app_id"`
	AppSecret string `mapstructure:"app_secret" json:"app_secret"`
}

type ServerConfig struct {
	Name        string              `mapstructure:"name" json:"name"`
	Host        string              `mapstructure:"host" json:"host"`
	MysqlInfo   MysqlConfig         `mapstructure:"mysql" json:"mysql"`
	OtelInfo    OtelConfig          `mapstructure:"otel" json:"otel"`
	RedisInfo   RedisConfig         `mapstructure:"redis" json:"redis"`
	FeedNumber  int                 `mapstructure:"feed_number" json:"feed_number"`
	UserSrvInfo config.RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
}
