package config

type NacosConfig struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint"`
	AccessKeyId     string `mapstructure:"access_key_id" json:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret" json:"access_key_secret"`
	Bucket          string `mapstructure:"bucket" json:"bucket"`
}
type JWTConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
}

type OtelConfig struct {
	EndPoint string `mapstructure:"endpoint" json:"endpoint"`
}

type ServerConfig struct {
	Name            string       `mapstructure:"name" json:"name"`
	Host            string       `mapstructure:"host" json:"host"`
	Port            int          `mapstructure:"port" json:"port"`
	JWTInfo         JWTConfig    `mapstructure:"jwt" json:"jwt"`
	OtelInfo        OtelConfig   `mapstructure:"otel" json:"otel"`
	UserSrvInfo     RPCSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	VideoSrvInfo    RPCSrvConfig `mapstructure:"video_srv" json:"video_srv"`
	MessageSrvInfo  RPCSrvConfig `mapstructure:"message_srv" json:"message_srv"`
	RelationSrvInfo RPCSrvConfig `mapstructure:"relation_srv" json:"relation_srv"`
	FavoriteSrvInfo RPCSrvConfig `mapstructure:"favorite_srv" json:"favorite_srv"`
	CommentSrvInfo  RPCSrvConfig `mapstructure:"comment_srv" json:"comment_srv"`
	OSSInfo         OSSConfig    `mapstructure:"oss" json:"oss"`
}

type RPCSrvConfig struct {
	Name      string `mapstructure:"name" json:"name"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
}
