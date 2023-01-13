package driver

import (
	"douyin/config"
	"douyin/storage"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var Db *gorm.DB

var RedisClient *redis.Client

func InitConn(str string) {
	switch str {
	case "mysql":
		Db = mySql()
	default:
		mySql()
	}
}

// InitRedis 初始化redis
func InitRedis() {

	password := config.Config.GetString("douyin.redis.password")

	host := config.Config.GetString("douyin.redis.host")

	db := config.Config.GetInt("douyin.redis.db")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:       host,
		Password:   password,
		DB:         db,
		MaxConnAge: 30 * time.Second,
	})
	//cache.CacheDb = cache.NewCacheDb(Redis)
}

// mysql connect
func mySql() *gorm.DB {

	user := config.Config.GetString("douyin.mysql.user")

	password := config.Config.GetString("douyin.mysql.password")

	host := config.Config.GetString("douyin.mysql.host")

	dbName := config.Config.GetString("douyin.mysql.dbName")

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", user, password, host, dbName)
	//dsn := "root:YYL521wxl@@tcp(47.105.50.53:3306)/go_blog?charset=utf8&parseTime=True&loc=Local"
	//db, err := gorm.Open("mysql", "root:YYL521wxl@@tcp(47.105.50.53:3306)/go_blog?charset=utf8&parseTime=True&loc=Local")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	//defer db.Close()
	mode := config.Config.GetString("douyin.log.mode")
	if mode == "DEBUG" {
		db.Logger.LogMode(2)
	}

	return db
}

// InitOSS 初始化阿里云oss
func InitOSS() {

	Endpoint := config.Config.GetString("Endpoint")

	AccessKeyId := config.Config.GetString("AccessKeyId")

	AccessKeySecret := config.Config.GetString("AccessKeySecret")

	Bucket := config.Config.GetString("Bucket")

	client, err := oss.New(Endpoint, AccessKeyId, AccessKeySecret)
	if err != nil {
		panic(err)
	}
	storage.OSS = storage.NewOSSClient(client, Bucket, Endpoint)
}
