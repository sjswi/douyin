package config

import (
	"github.com/spf13/viper"
)

var Config *viper.Viper

var Number int

func InitConfig() {
	configPATH := "./douyin.yaml"
	viper.SetConfigFile(configPATH)
	viper.SetConfigName("douyin")
	viper.SetConfigType("yaml") // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到错误；如果需要可以忽略
		} else {
			// 配置文件被找到，但产生了另外的错误
		}
		panic(err)
	}
	//TempPath = viper.GetString("douyin.video.tempPath")
	Number = viper.GetInt("douyin.feed.number")
	Config = viper.GetViper()
}
