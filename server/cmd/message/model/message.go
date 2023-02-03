package model

import (
	"douyin_rpc/common/cache"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"strconv"
)

type Message struct {
	gorm.Model
	ID         int64  `gorm:"primary_key; not null"`
	Content    string `gorm:"type:text; not null"`
	UserID     int64  `gorm:"index; not null"`
	TargetID   int64  `gorm:"index; not null"`
	CreateTime int64  `gorm:"index; not null"`
}

func (b *Message) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const MessageCachePrefix string = "message:message_"

func queryMessageByUserID(tx *gorm.DB, userID uint) ([]Message, error) {
	var Messages []Message
	if err := tx.Table("message").Where("user_id=?", userID).Find(&Messages).Error; err != nil {
		return nil, err
	}
	return Messages, nil
}

func QueryMessageByUserIDWithCache(tx *gorm.DB, userID uint) ([]Message, error) {
	key := MessageCachePrefix + "UserID_" + strconv.Itoa(int(userID))
	// 查看key是否存在
	//不存在

	var result string
	var Messages []Message
	var err error
	if !cache.Exist(key) {
		Messages, err = queryMessageByUserID(tx, userID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Messages)
		if err != nil {
			return nil, err
		}
		return Messages, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Messages, err = queryMessageByUserID(tx, userID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Messages)
			if err != nil {
				return nil, err
			}
			return Messages, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Messages)
	if err != nil {
		return nil, err
	}
	return Messages, nil
}

func queryMessageByTargetID(tx *gorm.DB, targetID uint) ([]Message, error) {
	var Messages []Message
	if err := tx.Table("message").Where("target_id=?", targetID).Find(&Messages).Error; err != nil {
		return nil, err
	}
	return Messages, nil
}

func QueryMessageByTargetIDWithCache(tx *gorm.DB, targetID uint) ([]Message, error) {
	key := MessageCachePrefix + "TargetID_" + strconv.Itoa(int(targetID))
	// 查看key是否存在
	//不存在
	var result string
	var Messages []Message
	var err error
	if !cache.Exist(key) {
		Messages, err = queryMessageByTargetID(tx, targetID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Messages)
		if err != nil {
			return nil, err
		}
		return Messages, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Messages, err = queryMessageByTargetID(tx, targetID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Messages)
			if err != nil {
				return nil, err
			}
			return Messages, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Messages)
	if err != nil {
		return nil, err
	}
	return Messages, nil
}

func queryMessageByUserIDAndTargetID(tx *gorm.DB, userID, targetID uint) ([]Message, error) {
	var messages []Message
	if err := tx.Table("message").Where("user_id=? and target_id=?", userID, targetID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func QueryMessageByUserIDAndTargetIDWithCache(tx *gorm.DB, userID, targetID uint) ([]Message, error) {
	key := MessageCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_TargetID_" + strconv.Itoa(int(targetID))
	var result string
	var messages []Message
	var err error
	// 查看key是否存在
	//不存在
	if !cache.Exist(key) {
		messages, err = queryMessageByUserIDAndTargetID(tx, userID, targetID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, messages)
		if err != nil {
			return nil, err
		}
		return messages, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			messages, err = queryMessageByUserIDAndTargetID(tx, userID, targetID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, messages)
			if err != nil {
				return nil, err
			}
			return messages, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

func UpdateMessage(tx *gorm.DB, message Message) error {
	if err := tx.Save(&message).Error; err != nil {
		return err
	}
	return nil
}

func CreateMessage(tx *gorm.DB, message Message) error {
	if err := tx.Model(message).Create(&message).Error; err != nil {
		return err
	}
	return nil
}
