package models

import (
	"douyin/cache"
	"encoding/json"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"strconv"
)

//TODO
// 关联模式，将user，favorite，comment关联到user

type User struct {
	gorm.Model
	Name          string
	FollowCount   int
	FollowerCount int
	Password      string
	Salt          string
}

const UserCachePrefix string = "user:user_"

// queryUserByID 查询数据库的user
func queryUserByID(tx *gorm.DB, userID uint) (*User, error) {
	// 直接查询数据库
	var user User
	if err := tx.Model(User{}).Where("id=?", userID).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUserByIDWithCache 通过视频id查询视频信息
func QueryUserByIDWithCache(tx *gorm.DB, userId uint) (*User, error) {
	key := UserCachePrefix + "ID_" + strconv.Itoa(int(userId))
	// 查看key是否存在
	//不存在

	var result string
	var user *User
	var err error
	if !cache.Exist(key) {
		user, err = queryUserByID(tx, userId)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, *user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			user, err = queryUserByID(tx, userId)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, *user)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser 更新数据库的user
func UpdateUser(tx *gorm.DB, user User) error {
	if err := tx.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

// queryUserByName 通过userID查询user
func queryUserByName(tx *gorm.DB, name string) (*User, error) {
	var user User
	if err := tx.Model(User{}).Where("name=?", name).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUserByNameWithCache 先查缓存再查数据库
func QueryUserByNameWithCache(tx *gorm.DB, name string) (*User, error) {
	key := UserCachePrefix + "Name_" + name
	// 查看key是否存在
	//不存在

	var result string
	var user *User
	var err error
	if !cache.Exist(key) {
		user, err = queryUserByName(tx, name)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, user)
		if err != nil {
			return nil, err
		}
		return user, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			user, err = queryUserByName(tx, name)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, user)
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func CreateUser(tx *gorm.DB, user *User) error {
	if err := tx.Model(User{}).Create(&user).Error; err != nil {
		return err
	}
	return nil
}
