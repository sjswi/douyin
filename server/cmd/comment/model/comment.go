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

type Comment struct {
	gorm.Model
	ID         int64  `gorm:"primary_key; not null"`
	CreateTime int64  `gorm:"index; not null"`
	Content    string `gorm:"type:text; default '';"`
	VideoID    int64  `gorm:"index; not null"`
	UserID     int64  `gorm:"index; not null"`
}

func (b *Comment) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const CommentCachePrefix string = "comment:comment_"

func queryCommentByUserID(tx *gorm.DB, userID int64) ([]Comment, error) {
	var Comments []Comment
	if err := tx.Table("comment").Where("user_id=?", userID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByUserIDWithCache(tx *gorm.DB, userID int64) ([]Comment, error) {
	key := CommentCachePrefix + "UserID_" + strconv.Itoa(int(userID))
	// 查看key是否存在
	//不存在

	var result string
	var Comments []Comment
	var err error
	if !cache.Exist(key) {
		Comments, err = queryCommentByUserID(tx, userID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Comments)
		if err != nil {
			return nil, err
		}
		return Comments, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Comments, err = queryCommentByUserID(tx, userID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Comments)
			if err != nil {
				return nil, err
			}
			return Comments, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Comments)
	if err != nil {
		return nil, err
	}
	return Comments, nil
}

func queryCommentByID(tx *gorm.DB, ID int64) (*Comment, error) {
	var Comments *Comment
	if err := tx.Table("comment").Where("id=?", ID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByIDWithCache(tx *gorm.DB, ID int64) (*Comment, error) {
	key := CommentCachePrefix + "ID_" + strconv.Itoa(int(ID))
	// 查看key是否存在
	//不存在

	var result string
	var Comments *Comment
	var err error
	if !cache.Exist(key) {
		Comments, err = queryCommentByID(tx, ID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Comments)
		if err != nil {
			return nil, err
		}
		return Comments, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Comments, err = queryCommentByID(tx, ID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Comments)
			if err != nil {
				return nil, err
			}
			return Comments, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Comments)
	if err != nil {
		return nil, err
	}
	return Comments, nil
}

func queryCommentByVideoID(tx *gorm.DB, videoID int64) ([]Comment, error) {
	var Comments []Comment
	if err := tx.Table("comment").Where("video_id=?", videoID).Order("created_at DESC").Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}
func CountCommentByVideoID(tx *gorm.DB, videoID int64) (int64, error) {
	var count int64
	if err := tx.Table("comment").Where("video_id=?", videoID).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
func QueryCommentByVideoIDWithCache(tx *gorm.DB, videoID int64) ([]Comment, error) {
	return queryCommentByVideoID(tx, videoID)
	//key := CommentCachePrefix + "VideoID_" + strconv.Itoa(int(videoID))
	//// 查看key是否存在
	////不存在
	//var result string
	//var Comments []Comment
	//var err error
	//if !cache.Exist(key) {
	//	Comments, err = queryCommentByVideoID(tx, videoID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	// 从数据库查出，放进redis
	//	err := cache.Set(key, Comments)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return Comments, nil
	//}
	////TODO
	//// lua脚本优化，保证原子性
	//
	////查询redis
	//if result, err = cache.Get(key); err != nil {
	//	// 极端情况：在判断存在后查询前过期了
	//	if err.Error() == "redis: nil" {
	//		Comments, err = queryCommentByVideoID(tx, videoID)
	//		if err != nil {
	//			return nil, err
	//		}
	//		// 从数据库查出，放进redis
	//		err := cache.Set(key, Comments)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return Comments, nil
	//	}
	//	return nil, err
	//}
	//// 反序列化
	//err = json.Unmarshal(strings.StringToBytes(result), &Comments)
	//if err != nil {
	//	return nil, err
	//}
	//return Comments, nil
}

func queryCommentByUserIDAndVideoID(tx *gorm.DB, userID, videoID int64) ([]Comment, error) {
	var Comments []Comment
	if err := tx.Table("comment").Where("user_id=? and video_id=?", userID, videoID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByUserIDAndVideoIDWithCache(tx *gorm.DB, userID, videoID int64) ([]Comment, error) {
	return queryCommentByUserIDAndVideoID(tx, userID, videoID)
	//key := CommentCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_videoID_" + strconv.Itoa(int(videoID))
	//var result string
	//var Comment *Comment
	//var err error
	//// 查看key是否存在
	////不存在
	//if !cache.Exist(key) {
	//	comments, err = queryCommentByUserIDAndVideoID(tx, userID, videoID)
	//	if err != nil {
	//		return nil, err
	//	}
	//	// 从数据库查出，放进redis
	//	err := cache.Set(key, Comment)
	//	if err != nil {
	//		return nil, err
	//	}
	//	return Comment, nil
	//}
	////TODO
	//// lua脚本优化，保证原子性
	//
	////查询redis
	//if result, err = cache.Get(key); err != nil {
	//	// 极端情况：在判断存在后查询前过期了
	//	if err.Error() == "redis: nil" {
	//		Comment, err = queryCommentByUserIDAndVideoID(tx, userID, videoID)
	//		if err != nil {
	//			return nil, err
	//		}
	//		// 从数据库查出，放进redis
	//		err := cache.Set(key, Comment)
	//		if err != nil {
	//			return nil, err
	//		}
	//		return Comment, nil
	//	}
	//	return nil, err
	//}
	//// 反序列化
	//err = json.Unmarshal(strings.StringToBytes(result), &Comment)
	//if err != nil {
	//	return nil, err
	//}
	//return Comment, nil
}

func UpdateComment(tx *gorm.DB, comment Comment) error {
	if err := tx.Save(&comment).Error; err != nil {
		return err
	}
	return nil
}

func CreateComment(tx *gorm.DB, comment *Comment) error {
	if err := tx.Table("comment").Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(tx *gorm.DB, comment Comment) error {
	if err := tx.Table("comment").Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommentByID(tx *gorm.DB, commentID int64) error {
	if err := tx.Table("comment").Delete(&Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}
