package model

import (
	"douyin_rpc/server/cmd/user/global"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type User struct {
	gorm.Model

	ID       int64  `gorm:"primary_key; not null"`
	Name     string `gorm:"type:varchar(32); not null; default '';index: unique"`
	Password string `gorm:"type:varchar(32); not null; default '';"`
	Salt     string `gorm:"type:varchar(8); not null; default '';"`
}

func (b *User) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const UserCachePrefix string = "user:user:"

var mutex sync.Mutex
var update_keys []string

// queryUserByID 查询数据库的user
func queryUserByID(tx *gorm.DB, userID int64) (*User, error) {
	// 直接查询数据库
	var user User
	if err := tx.Table("user").Where("id=?", userID).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUserByIDWithCache 通过视频id查询视频信息
func QueryUserByIDWithCache(tx *gorm.DB, userId int64) (*User, error) {
	key := UserCachePrefix + "ID:" + strconv.Itoa(int(userId))
	var user *User
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		user, err = queryUserByID(tx, userId)
		if err != nil {
			return "", nil
		}
		data, err := sonic.Marshal(user)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

// queryUserByName 通过userID查询user
func queryUserByName(tx *gorm.DB, name string) (*User, error) {
	var user User
	if err := tx.Table("user").Where("name=?", name).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// QueryUserByNameWithCache 先查缓存再查数据库
func QueryUserByNameWithCache(tx *gorm.DB, name string) (*User, error) {
	//return queryUserByName(tx, name)
	key := UserCachePrefix + "Name:" + name
	var user *User
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		user, err = queryUserByName(tx, name)
		if err != nil {
			return "", nil
		}
		data, err := sonic.Marshal(user)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &user)
	if err != nil {
		return nil, err
	}
	return user, nil

}

func CreateUser(tx *gorm.DB, user *User) error {
	if err := tx.Table("user").Create(&user).Error; err != nil {

		return err
	}
	updateKeys(user)
	return nil
}

func DeleteCache() {
	mutex.Lock()
	temp := update_keys[:]
	update_keys = update_keys[:0]
	mutex.Unlock()
	for {
		err := global.RocksCacheClient.TagAsDeletedBatch(temp)
		if err == nil {
			break
		}
	}

	return
}

func updateKeys(user *User) {
	key1 := UserCachePrefix + "ID:" + strconv.FormatInt(user.ID, 10)
	key2 := UserCachePrefix + "Name:" + user.Name
	mutex.Lock()
	update_keys = append(update_keys, []string{key1, key2}...)
	mutex.Unlock()
}
