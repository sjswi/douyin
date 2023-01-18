package models

import (
	"douyin/cache"
	"encoding/json"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type Comment struct {
	gorm.Model
	Content    string
	CreateTime time.Time
	VideoID    uint
	UserID     uint
}

const CommentCachePrefix string = "comment:comment_"

func queryCommentByUserID(tx *gorm.DB, userID uint) ([]Comment, error) {
	var Comments []Comment
	if err := tx.Model(Comment{}).Where("user_id=?", userID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByUserIDWithCache(tx *gorm.DB, userID uint) ([]Comment, error) {
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

func queryCommentByID(tx *gorm.DB, ID uint) (*Comment, error) {
	var Comments *Comment
	if err := tx.Model(Comment{}).Where("id=?", ID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByIDWithCache(tx *gorm.DB, ID uint) (*Comment, error) {
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

func queryCommentByVideoID(tx *gorm.DB, videoID uint) ([]Comment, error) {
	var Comments []Comment
	if err := tx.Model(Comment{}).Where("video_id=?", videoID).Order("created_at DESC").Find(&Comments).Error; err != nil {
		return nil, err
	}
	return Comments, nil
}

func QueryCommentByVideoIDWithCache(tx *gorm.DB, videoID uint) ([]Comment, error) {
	key := CommentCachePrefix + "VideoID_" + strconv.Itoa(int(videoID))
	// 查看key是否存在
	//不存在
	var result string
	var Comments []Comment
	var err error
	if !cache.Exist(key) {
		Comments, err = queryCommentByVideoID(tx, videoID)
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
			Comments, err = queryCommentByVideoID(tx, videoID)
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

func queryCommentByUserIDAndVideoID(tx *gorm.DB, userID, videoID uint) (*Comment, error) {
	var Comments Comment
	if err := tx.Model(Comment{}).Where("user_id=? and video_id=?", userID, videoID).Find(&Comments).Error; err != nil {
		return nil, err
	}
	return &Comments, nil
}

func QueryCommentByUserIDAndVideoIDWithCache(tx *gorm.DB, userID, videoID uint) (*Comment, error) {
	key := CommentCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_videoID_" + strconv.Itoa(int(videoID))
	var result string
	var Comment *Comment
	var err error
	// 查看key是否存在
	//不存在
	if !cache.Exist(key) {
		Comment, err = queryCommentByUserIDAndVideoID(tx, userID, videoID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Comment)
		if err != nil {
			return nil, err
		}
		return Comment, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Comment, err = queryCommentByUserIDAndVideoID(tx, userID, videoID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Comment)
			if err != nil {
				return nil, err
			}
			return Comment, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Comment)
	if err != nil {
		return nil, err
	}
	return Comment, nil
}

func UpdateComment(tx *gorm.DB, comment Comment) error {
	if err := tx.Save(&comment).Error; err != nil {
		return err
	}
	return nil
}

func CreateComment(tx *gorm.DB, comment Comment) error {
	if err := tx.Model(comment).Create(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteComment(tx *gorm.DB, comment Comment) error {
	if err := tx.Model(comment).Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}

func DeleteCommentByID(tx *gorm.DB, commentID uint) error {
	if err := tx.Model(Comment{}).Delete(&Comment{}, commentID).Error; err != nil {
		return err
	}
	return nil
}
