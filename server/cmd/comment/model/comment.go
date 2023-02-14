package model

import (
	"douyin_rpc/server/cmd/comment/global"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
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
	return queryCommentByUserID(tx, userID)
	//key := CommentCachePrefix + "UserID_" + strconv.Itoa(int(userID))
	//// 查看key是否存在
	////不存在
	//
	//var result string
	//var Comments []Comment
	//var err error
	//if !cache.Exist(key) {
	//	Comments, err = queryCommentByUserID(tx, userID)
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
	////查询redis
	//if result, err = cache.Get(key); err != nil {
	//	// 极端情况：在判断存在后查询前过期了
	//	if err.Error() == "redis: nil" {
	//		Comments, err = queryCommentByUserID(tx, userID)
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

var mutex sync.Mutex
var update_keys []string

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

	var Comments *Comment
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		Comments, err = queryCommentByID(tx, ID)
		if err != nil {
			return "", err
		}
		// 从数据库查出，放进redis
		tss, err := sonic.Marshal(Comments)
		if err != nil {
			return "", err
		}
		return string(tss), nil
	})
	if err != nil {
		return nil, err
	}
	if Comments != nil {
		return Comments, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &Comments)
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
	key := CommentCachePrefix + "Video_ID_" + strconv.FormatInt(videoID, 10) + "_count"

	var count *int64
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		if err := tx.Table("comment").Where("video_id=?", videoID).Count(count).Error; err != nil {
			return "", err
		}
		return strconv.FormatInt(*count, 10), nil
	})
	if err != nil {
		return -1, err
	}
	if count != nil {
		return *count, nil
	}
	*count, err = strconv.ParseInt(fetch, 0, 64)
	return *count, nil
}
func QueryCommentByVideoIDWithCache(tx *gorm.DB, videoID int64) ([]Comment, error) {
	key := CommentCachePrefix + "Video_ID_" + strconv.FormatInt(videoID, 10)

	var Comments []Comment
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		Comments, err = queryCommentByVideoID(tx, videoID)
		if err != nil {
			return "", nil
		}
		data, err := sonic.Marshal(Comments)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if Comments != nil {
		return Comments, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &Comments)
	if err != nil {
		return nil, err
	}
	return Comments, nil

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
	// 当删除一个评论时，需要删除对应的缓存，目前来说只有VideoId对应的缓存和ID对应的缓存
	if err := tx.Table("comment").Delete(&comment).Error; err != nil {
		return err
	}
	return nil
}
func DeleteCache() error {
	mutex.Lock()
	temp := update_keys[:]
	update_keys = update_keys[:0]
	mutex.Unlock()
	err := global.RocksCacheClient.TagAsDeletedBatch(temp)
	if err != nil {
		return err
	}
	return nil
}

func DeleteCommentByID(tx *gorm.DB, commentID int64) error {
	var temp Comment
	if err := tx.Table("comment").Delete(&temp, commentID).Error; err != nil {
		return err
	}

	key1 := CommentCachePrefix + "ID_" + strconv.FormatInt(commentID, 10)
	key2 := CommentCachePrefix + "Video_ID_" + strconv.FormatInt(temp.VideoID, 10)
	key3 := CommentCachePrefix + "Video_ID_" + strconv.FormatInt(temp.VideoID, 10) + "_count"
	mutex.Lock()
	update_keys = append(update_keys, []string{key3, key2, key1}...)
	mutex.Unlock()
	return nil
}
