package model

import (
	"douyin_rpc/server/cmd/favorite/global"
	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"sync"
	"time"
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

const FavoriteCachePrefix string = "favorite:favorite:"

func queryFavoriteByUserID(tx *gorm.DB, userID int64) ([]Favorite, error) {
	var Favorites []Favorite
	if err := tx.Table("favorite").Where("exist=1").Where("user_id=?", userID).Find(&Favorites).Error; err != nil {
		return nil, err
	}
	return Favorites, nil
}

func QueryFavoriteByUserIDWithCache(tx *gorm.DB, userID int64) ([]Favorite, error) {
	key := FavoriteCachePrefix + "User_ID_::" + strconv.FormatInt(userID, 10)
	var favorites []Favorite
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		favorites, err = queryFavoriteByUserID(tx, userID)
		if err != nil {
			return "", err
		}
		data, err := sonic.Marshal(favorites)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if favorites != nil {
		return favorites, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &favorites)
	if err != nil {
		return nil, err
	}
	return favorites, nil

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
	key := FavoriteCachePrefix + "Video_ID:" + strconv.FormatInt(videoID, 10)
	var favorites []Favorite
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		favorites, err = queryFavoriteByVideoID(tx, videoID)
		if err != nil {
			return "", err
		}
		data, err := sonic.Marshal(favorites)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if favorites != nil {
		return favorites, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &favorites)
	if err != nil {
		return nil, err
	}
	return favorites, nil

}

func queryFavoriteByUserIDAndVideoID(tx *gorm.DB, userID, videoID int64) (*Favorite, error) {
	var Favorites Favorite
	if err := tx.Table("favorite").Where("exist=1").Where("user_id=? and video_id=?", userID, videoID).Find(&Favorites).Error; err != nil {
		return nil, err
	}
	return &Favorites, nil
}

func QueryFavoriteByUserIDAndVideoIDWithCache(tx *gorm.DB, userID, videoID int64) (*Favorite, error) {
	key := FavoriteCachePrefix + "User_ID:" + strconv.FormatInt(userID, 10) + ":Video_ID:" + strconv.FormatInt(videoID, 10)
	var favorite *Favorite
	var err error
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		favorite, err = queryFavoriteByUserIDAndVideoID(tx, userID, videoID)
		if err != nil {
			return "", err
		}
		data, err := sonic.Marshal(favorite)
		if err != nil {
			return "", err
		}
		return string(data), nil
	})
	if err != nil {
		return nil, err
	}
	if favorite != nil {
		return favorite, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &favorite)
	if err != nil {
		return nil, err
	}
	return favorite, nil

}

var mutex sync.Mutex
var update_keys []string

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
func UpdateFavorite(tx *gorm.DB, favorite Favorite) error {
	//TODO
	// 删除缓存
	if err := tx.Save(&favorite).Error; err != nil {
		return err
	}
	updateKeys(favorite)
	return nil
}
func updateKeys(favorite Favorite) {
	key1 := FavoriteCachePrefix + "User_ID:" + strconv.FormatInt(favorite.UserID, 10) + ":Video_ID:" + strconv.FormatInt(favorite.VideoID, 10)
	key2 := FavoriteCachePrefix + "User_ID_:" + strconv.FormatInt(favorite.UserID, 10)
	key3 := FavoriteCachePrefix + "User_ID_:" + strconv.FormatInt(favorite.VideoID, 10)
	key4 := FavoriteCachePrefix + "Video_ID:" + strconv.FormatInt(favorite.VideoID, 10)
	key5 := FavoriteCachePrefix + "Video_ID:" + strconv.FormatInt(favorite.UserID, 10)
	mutex.Lock()
	update_keys = append(update_keys, []string{key1, key4, key3, key2, key5}...)
	mutex.Unlock()
}
func CreateFavorite(tx *gorm.DB, favorite Favorite) error {
	if err := tx.Table("favorite").Create(&favorite).Error; err != nil {
		return err
	}
	updateKeys(favorite)
	return nil
}

func UpdateOrCreateFavorite(tx *gorm.DB, favorite Favorite) error {
	if err := tx.Table("favorite").Clauses(clause.OnConflict{
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
	updateKeys(favorite)
	return nil
}
