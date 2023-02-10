package initialize

import (
	"douyin_rpc/server/cmd/comment/global"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {

	user := global.ServerConfig.MysqlInfo.User

	password := global.ServerConfig.MysqlInfo.Password

	host := global.ServerConfig.MysqlInfo.Host

	dbName := global.ServerConfig.MysqlInfo.Name

	maxConn := global.ServerConfig.MysqlInfo.MaxConn

	maxIdle := global.ServerConfig.MysqlInfo.MaxIdle
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	//defer db.Close()
	//mode := config.Config.GetString("douyin.log.mode")
	//if mode == "DEBUG" {
	//	db.Logger.LogMode(2)
	//}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxIdleConns(maxIdle)
	sqlDb.SetMaxOpenConns(maxConn)
	global.DB = db
}
