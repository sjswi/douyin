package cache

import (
	"context"
	"douyin/bootstrap/driver"
	"encoding/json"
	"github.com/u2takey/go-utils/strings"
	"time"
)

func Set(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if _, err = driver.RedisClient.Set(context.Background(), key, strings.BytesToString(data), 3*time.Hour).Result(); err != nil {
		return err
	}
	return nil
}

func Get(key string) (result string, err error) {

	if result, err = driver.RedisClient.Get(context.Background(), key).Result(); err != nil {
		return "", err
	}
	return
}

func Delete(keys []string) error {
	if _, err := driver.RedisClient.Del(context.Background(), keys...).Result(); err != nil {
		return err
	}
	return nil
}

func Exist(key string) bool {
	exist, err := driver.RedisClient.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	if exist == 0 {
		return false
	}
	return true
}
