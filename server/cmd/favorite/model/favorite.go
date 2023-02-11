package model

import (
	"douyin_rpc/common/cache"
	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
)

type Favorite struct {
	gorm.Model
	ID      int64 `gorm:"primary_key; not null"`
	UserID  int64 `gorm:"index; not null;"` //用户id
	VideoID int64 `gorm:"index; not null;"` //视频id
	Exist   bool  //是否存在，避免重复的点赞取消点赞使得该表内容变得很大
}

func (b *Favorite) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const FavoriteCachePrefix string = "favorite:favorite_"

func queryFavoriteByUserID(tx *gorm.DB, userID int64) ([]Favorite, error) {
	var Favorites []Favorite
	if err := tx.Table("favorite").Where("exist=1").Where("user_id=?", userID).Find(&Favorites).Error; err != nil {
		return nil, err
	}
	return Favorites, nil
}

func QueryFavoriteByUserIDWithCache(tx *gorm.DB, userID int64) ([]Favorite, error) {
	key := FavoriteCachePrefix + "UserID_" + strconv.Itoa(int(userID))
	// 查看key是否存在
	//不存在

	var result string
	var Favorites []Favorite
	var err error
	if !cache.Exist(key) {
		Favorites, err = queryFavoriteByUserID(tx, userID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Favorites)
		if err != nil {
			return nil, err
		}
		return Favorites, nil
	}
	//TODO
	// lua脚本优化，保证原子性
	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Favorites, err = queryFavoriteByUserID(tx, userID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Favorites)
			if err != nil {
				return nil, err
			}
			return Favorites, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Favorites)
	if err != nil {
		return nil, err
	}
	return Favorites, nil
}

func queryFavoriteByVideoID(tx *gorm.DB, videoID int64) ([]Favorite, error) {
	var Favorites []Favorite
	if err := tx.Table("favorite").Where("exist=1").Where("video_id=?", videoID).Find(&Favorites).Error; err != nil {
		return nil, err
	}
	return Favorites, nil
}
func CountFavoriteByVideoID(tx *gorm.DB, videoID int64) (int64, error) {
	var count int64
	if err := tx.Table("favorite").Where("exist=1").Where("video_id=?", videoID).Count(&count).Error; err != nil {
		return -1, err
	}
	return count, nil
}
func QueryFavoriteByVideoIDWithCache(tx *gorm.DB, videoID int64) ([]Favorite, error) {
	key := FavoriteCachePrefix + "videoID_" + strconv.Itoa(int(videoID))
	// 查看key是否存在
	//不存在
	var result string
	var Favorites []Favorite
	var err error
	if !cache.Exist(key) {
		Favorites, err = queryFavoriteByVideoID(tx, videoID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Favorites)
		if err != nil {
			return nil, err
		}
		return Favorites, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Favorites, err = queryFavoriteByVideoID(tx, videoID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Favorites)
			if err != nil {
				return nil, err
			}
			return Favorites, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Favorites)
	if err != nil {
		return nil, err
	}
	return Favorites, nil
}

func queryFavoriteByUserIDAndVideoID(tx *gorm.DB, userID, videoID int64) (*Favorite, error) {
	var Favorites Favorite
	if err := tx.Table("favorite").Where("exist=1").Where("user_id=? and video_id=?", userID, videoID).Find(&Favorites).Error; err != nil {
		return nil, err
	}
	return &Favorites, nil
}

func QueryFavoriteByUserIDAndVideoIDWithCache(tx *gorm.DB, userID, videoID int64) (*Favorite, error) {
	key := FavoriteCachePrefix + "UserID_" + strconv.Itoa(int(userID)) + "_VideoID_" + strconv.Itoa(int(videoID))
	var result string
	var Favorite *Favorite
	var err error
	// 查看key是否存在
	//不存在
	if !cache.Exist(key) {
		Favorite, err = queryFavoriteByUserIDAndVideoID(tx, userID, videoID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, Favorite)
		if err != nil {
			return nil, err
		}
		return Favorite, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			Favorite, err = queryFavoriteByUserIDAndVideoID(tx, userID, videoID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, Favorite)
			if err != nil {
				return nil, err
			}
			return Favorite, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &Favorite)
	if err != nil {
		return nil, err
	}
	return Favorite, nil
}

func UpdateFavorite(tx *gorm.DB, favorite Favorite) error {
	//TODO
	// 删除缓存
	if err := tx.Save(&favorite).Error; err != nil {
		return err
	}
	return nil
}

func CreateFavorite(tx *gorm.DB, favorite Favorite) error {
	if err := tx.Table("favorite").Create(&favorite).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOrCreateFavorite(tx *gorm.DB, favorite Favorite) error {
	if err := tx.Clauses(clause.OnConflict{
		Columns:      []clause.Column{{Name: "user_id"}, {Name: "video_id"}},
		Where:        clause.Where{},
		TargetWhere:  clause.Where{},
		OnConstraint: "",
		DoNothing:    false,
		DoUpdates:    clause.Assignments(map[string]interface{}{"exist": true}),
		UpdateAll:    false,
	}).Create(&favorite).Error; err != nil {
		return err
	}
	return nil
}
