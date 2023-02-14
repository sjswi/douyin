package model

import (
	"context"
	"douyin_rpc/server/cmd/message/global"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/go-redis/redis/v8"
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

const MessageCachePrefix string = "message:message:"

//var mutex sync.Mutex
//var update_keys []string
//
//func DeleteCache() {
//	mutex.Lock()
//	temp := update_keys[:]
//	update_keys = update_keys[:0]
//	mutex.Unlock()
//	for {
//		err := global.RocksCacheClient.TagAsDeletedBatch(temp)
//		if err == nil {
//			break
//		}
//	}
//
//	return
//}

//func updateKeys(message *Message) {
//	key1 := MessageCachePrefix + "User_ID:" + strconv.FormatInt(message.ID, 10)
//	key2 := MessageCachePrefix + "Name:" + message.Name
//	mutex.Lock()
//	update_keys = append(update_keys, []string{key1, key2}...)
//	mutex.Unlock()
//}

//
//func queryMessageByUserID(tx *gorm.DB, userID int64) ([]Message, error) {
//	var Messages []Message
//	if err := tx.Table("message").Where("user_id=?", userID).Find(&Messages).Error; err != nil {
//		return nil, err
//	}
//	return Messages, nil
//}
//
//func QueryMessageByUserIDWithCache(tx *gorm.DB, userID int64) ([]Message, error) {
//	key := MessageCachePrefix + "UserID_" + strconv.Itoa(int(userID))
//	// 查看key是否存在
//	//不存在
//
//	var result string
//	var Messages []Message
//	var err error
//	if !cache.Exist(key) {
//		Messages, err = queryMessageByUserID(tx, userID)
//		if err != nil {
//			return nil, err
//		}
//		// 从数据库查出，放进redis
//		err := cache.Set(key, Messages)
//		if err != nil {
//			return nil, err
//		}
//		return Messages, nil
//	}
//	//TODO
//	// lua脚本优化，保证原子性
//	//查询redis
//	if result, err = cache.Get(key); err != nil {
//		// 极端情况：在判断存在后查询前过期了
//		if err.Error() == "redis: nil" {
//			Messages, err = queryMessageByUserID(tx, userID)
//			if err != nil {
//				return nil, err
//			}
//			// 从数据库查出，放进redis
//			err := cache.Set(key, Messages)
//			if err != nil {
//				return nil, err
//			}
//			return Messages, nil
//		}
//		return nil, err
//	}
//	// 反序列化
//	err = json.Unmarshal(strings.StringToBytes(result), &Messages)
//	if err != nil {
//		return nil, err
//	}
//	return Messages, nil
//}
//
//func queryMessageByTargetID(tx *gorm.DB, targetID int64) ([]Message, error) {
//	var Messages []Message
//	if err := tx.Table("message").Where("target_id=?", targetID).Find(&Messages).Error; err != nil {
//		return nil, err
//	}
//	return Messages, nil
//}
//
//func QueryMessageByTargetIDWithCache(tx *gorm.DB, targetID int64) ([]Message, error) {
//	key := MessageCachePrefix + "TargetID_" + strconv.Itoa(int(targetID))
//	// 查看key是否存在
//	//不存在
//	var result string
//	var Messages []Message
//	var err error
//
//	if !cache.Exist(key) {
//		Messages, err = queryMessageByTargetID(tx, targetID)
//		if err != nil {
//			return nil, err
//		}
//		// 从数据库查出，放进redis
//		err := cache.Set(key, Messages)
//		if err != nil {
//			return nil, err
//		}
//		return Messages, nil
//	}
//	//TODO
//	// lua脚本优化，保证原子性
//
//	//查询redis
//	if result, err = cache.Get(key); err != nil {
//		// 极端情况：在判断存在后查询前过期了
//		if err.Error() == "redis: nil" {
//			Messages, err = queryMessageByTargetID(tx, targetID)
//			if err != nil {
//				return nil, err
//			}
//			// 从数据库查出，放进redis
//			err := cache.Set(key, Messages)
//			if err != nil {
//				return nil, err
//			}
//			return Messages, nil
//		}
//		return nil, err
//	}
//	// 反序列化
//	err = json.Unmarshal(strings.StringToBytes(result), &Messages)
//	if err != nil {
//		return nil, err
//	}
//	return Messages, nil
//}

// 问题：
//	 1、前端没有传时间的情况下，每次获取消息都是获取全部数据，而前端中还有之前的数据只是需要最新的而已
//   2、
//
func queryMessageByUserIDAndTargetID(tx *gorm.DB, userID, targetID int64) ([]Message, error) {
	var messages1, messages2 []Message
	//messageTime:=int64(-1)
	// 登入用户查询聊天记录
	//var messages []Message
	_, err := global.RocksCacheClient.RawGet(context.Background(), "message:message:User_ID:"+strconv.FormatInt(userID, 10)+":start")
	key := MessageCachePrefix + "User_ID:" + strconv.FormatInt(userID, 10) + ":Target_ID:" + strconv.FormatInt(targetID, 10) + ":time"
	if err == redis.Nil {
		// 非第一次登录，查询上一次查询的时间

		fetch, err := global.RocksCacheClient.Fetch(key, -1, func() (string, error) {
			// 没有查到
			if err := tx.Table("message").Where("user_id=? and target_id=?", userID, targetID).Find(&messages1).Order("create_time").Error; err != nil {
				return "", err
			}
			if err := tx.Table("message").Where("user_id=? and target_id=?", targetID, userID).Find(&messages2).Order("create_time").Error; err != nil {
				return "", err
			}
			if len(messages1) > 0 && len(messages2) > 0 {
				if messages1[len(messages1)-1].CreateTime > messages2[len(messages2)-1].CreateTime {
					return strconv.FormatInt(messages1[len(messages1)-1].CreateTime, 10), nil
				} else {
					return strconv.FormatInt(messages2[len(messages2)-1].CreateTime, 10), nil
				}
			} else if len(messages1) == 0 {
				return strconv.FormatInt(messages2[len(messages2)-1].CreateTime, 10), nil
			} else if len(messages2) == 0 {
				return strconv.FormatInt(messages1[len(messages1)-1].CreateTime, 10), nil
			} else {
				return "", nil
			}

		})

		if err != nil {
			return nil, err
		}
		messages1 = append(messages1, messages2...)

		if messages1 != nil {
			return messages1, nil
		}
		//查询成功
		messageTime, err := strconv.ParseInt(fetch, 0, 64)
		if err := tx.Table("message").Where("user_id=? and target_id=?", userID, targetID).Where("create_time > ?", messageTime).Find(&messages1).Order("create_time").Error; err != nil {
			return nil, err
		}
		if err := tx.Table("message").Where("user_id=? and target_id=?", userID, targetID).Where("create_time > ?", messageTime).Find(&messages2).Order("create_time").Error; err != nil {
			return nil, err
		}
		count := 0
		for {
			if global.RocksCacheClient.RawSet(context.Background(), key, strconv.FormatInt(messages1[len(messages1)-1].CreateTime, 10), -1) == nil {
				break
			}
			if count > 10 {
				break
			}
			count += 1
		}

	}
	// 第一次登录，不需要时间

	// 1、删除第一次登录标记
	global.RocksCacheClient.TagAsDeleted("message:message:User_ID:" + strconv.FormatInt(userID, 10) + ":start")
	// 2、删除上一次的时间标记
	global.RocksCacheClient.TagAsDeleted(key)
	// 查询数据库不加时间

	if err := tx.Table("message").Where("user_id=? and target_id=?", userID, targetID).Find(&messages1).Order("create_time").Error; err != nil {
		return nil, err
	}
	if err := tx.Table("message").Where("user_id=? and target_id=?", targetID, userID).Find(&messages2).Order("create_time").Error; err != nil {
		return nil, err
	}
	latestime := "-1"
	if len(messages1) > 0 && len(messages2) > 0 {
		if messages1[len(messages1)-1].CreateTime > messages2[len(messages2)-1].CreateTime {
			latestime = strconv.FormatInt(messages1[len(messages1)-1].CreateTime, 10)
		} else {
			latestime = strconv.FormatInt(messages2[len(messages2)-1].CreateTime, 10)
		}
	} else if len(messages1) == 0 {
		latestime = strconv.FormatInt(messages2[len(messages2)-1].CreateTime, 10)
	} else if len(messages2) == 0 {
		latestime = strconv.FormatInt(messages1[len(messages1)-1].CreateTime, 10)
	}

	count := 0
	for {
		if err = global.RocksCacheClient.RawSet(context.Background(), key, latestime, -1); err == nil {
			break
		}
		if count > 10 {
			break
		}
		count += 1
	}
	messages1 = append(messages1, messages2[:]...)
	// 更新数据库时间
	return messages1, nil
}

func QueryMessageByUserIDAndTargetIDWithCache(tx *gorm.DB, userID, targetID int64) ([]Message, error) {
	return queryMessageByUserIDAndTargetID(tx, userID, targetID)
	//key := MessageCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_TargetID_" + strconv.Itoa(int(targetID))
	//var result string
	//var messages []Message
	//var err error
	//// 查看key是否存在
	////不存在
	//if !cache.Exist(key) {
	//	messages, err = queryMessageByUserIDAndTargetID(tx, userID, targetID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	// 从数据库查出，放进redis
	//	err := cache.Set(key, messages)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return messages, nil
	//}
	////TODO
	//// lua脚本优化，保证原子性
	//
	////查询redis
	//if result, err = cache.Get(key); err != nil {
	//	// 极端情况：在判断存在后查询前过期了
	//	if err.Error() == "redis: nil" {
	//		messages, err = queryMessageByUserIDAndTargetID(tx, userID, targetID)
	//		if err != nil {
	//			return nil, err
	//		}
	//		// 从数据库查出，放进redis
	//		err := cache.Set(key, messages)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return messages, nil
	//	}
	//	return nil, err
	//}
	//// 反序列化
	//err = json.Unmarshal(strings.StringToBytes(result), &messages)
	//if err != nil {
	//	return nil, err
	//}
	//return messages, nil
}

func UpdateMessage(tx *gorm.DB, message Message) error {
	if err := tx.Save(&message).Error; err != nil {
		return err
	}
	return nil
}

func CreateMessage(tx *gorm.DB, message *Message) error {
	if err := tx.Table("message").Create(&message).Error; err != nil {
		return err
	}
	return nil
}
