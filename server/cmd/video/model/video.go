package model

import (
	"douyin_rpc/common/cache"
	"douyin_rpc/server/cmd/video/global"

	"encoding/json"
	"github.com/bwmarrin/snowflake"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/u2takey/go-utils/strings"
	"gorm.io/gorm"
	"strconv"
)

type Video struct {
	gorm.Model
	ID       int64  `gorm:"primary_key;not null"`
	AuthorID int64  `gorm:"index; not null"`                          //作者id
	Title    string `gorm:"type:varchar(50); not null; default='';"`  //视频标题
	PlayURL  string `gorm:"type:varchar(150); not null; default='';"` //播放地址
	CoverURL string `gorm:"type:varchar(150); not null; default='';"` //封面地址
}

func (b *Video) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

const VideoCachePrefix string = "video:video_"

// queryVideoByID 查询数据库的video
func queryVideoByID(tx *gorm.DB, videoId int64) (*Video, error) {
	// 直接查询数据库
	var video Video
	if err := tx.Model(Video{}).Where("id=?", videoId).Find(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// QueryVideoByIDWithCache 通过视频id查询视频信息
func QueryVideoByIDWithCache(tx *gorm.DB, videoID int64) (*Video, error) {
	key := VideoCachePrefix + "ID_" + strconv.Itoa(int(videoID))
	// 查看key是否存在
	//不存在
	var result string
	var err error
	var video *Video
	if !cache.Exist(key) {
		video, err = queryVideoByID(tx, videoID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, *video)
		if err != nil {
			return nil, err
		}
		return video, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			video, err = queryVideoByID(tx, videoID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, *video)
			if err != nil {
				return nil, err
			}
			return video, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &video)
	if err != nil {
		return nil, err
	}
	return video, nil
}

// UpdateVideo 更新数据库的video
func UpdateVideo(tx *gorm.DB, video *Video) error {
	if err := tx.Save(&video).Error; err != nil {
		return err
	}
	return nil
}

// queryVideoByAuthorID 通过videoID查询video
func queryVideoByAuthorID(tx *gorm.DB, authorID int64) ([]Video, error) {
	var videos []Video
	if err := tx.Model(Video{}).Where("video_id=?", authorID).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// QueryVideoByAuthorIDWithCache 先查缓存再查数据库
func QueryVideoByAuthorIDWithCache(tx *gorm.DB, authorID int64) ([]Video, error) {
	key := VideoCachePrefix + "AuthorID_" + strconv.Itoa(int(authorID))
	// 查看key是否存在

	//不存在
	var result string
	var err error
	var videos []Video
	if !cache.Exist(key) {
		videos, err = queryVideoByAuthorID(tx, authorID)
		if err != nil {
			return nil, err
		}
		// 从数据库查出，放进redis
		err := cache.Set(key, videos)
		if err != nil {
			return nil, err
		}
		return videos, nil
	}
	//TODO
	// lua脚本优化，保证原子性

	//查询redis
	if result, err = cache.Get(key); err != nil {
		// 极端情况：在判断存在后查询前过期了
		if err.Error() == "redis: nil" {
			videos, err = queryVideoByAuthorID(tx, authorID)
			if err != nil {
				return nil, err
			}
			// 从数据库查出，放进redis
			err := cache.Set(key, videos)
			if err != nil {
				return nil, err
			}
			return videos, nil
		}
		return nil, err
	}
	// 反序列化
	err = json.Unmarshal(strings.StringToBytes(result), &videos)
	if err != nil {
		return nil, err
	}
	return videos, nil
}

func CreateVideo(tx *gorm.DB, video *Video) error {
	if err := tx.Model(Video{}).Create(&video).Error; err != nil {
		return err
	}
	return nil
}

// redis使用:
// 1、往redis写数据只有查询操作会做！！！
// 2、修改数据库数据时，在事务提交后需要将redis的数据删除，因此需要记住redis中的key是如何存储的

// Feed 暂时没有想到使用缓存的好方法直接查数据库
func Feed(tx *gorm.DB, latestTime time.Time) ([]Video, error) {
	var videos []Video
	if err := tx.Model(Video{}).Where("created_at<=?", latestTime).Limit(global.ServerConfig.FeedNumber).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
