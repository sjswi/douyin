package cache

import (
	"context"
	"douyin_rpc/server/cmd/video/global"
	"encoding/json"
	"github.com/u2takey/go-utils/strings"
	"time"
)

func Set(key string, value interface{}) error {
	//TODO
	// 设置和删除需要加锁，需要锁续期等等操作
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = global.Redis.Set(context.Background(), key, strings.BytesToString(data), 3*time.Hour).Result(); err != nil {
		return err
	}
	return nil
}

func Get(key string) (result string, err error) {

	if result, err = global.Redis.Get(context.Background(), key).Result(); err != nil {
		return "", err
	}
	return
}

func Delete(keys []string) error {
	if _, err := global.Redis.Del(context.Background(), keys...).Result(); err != nil {
		return err
	}
	return nil
}

func Exist(key string) bool {
	exist, err := global.Redis.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if exist == 0 {
		return false
	}
	return true
}
