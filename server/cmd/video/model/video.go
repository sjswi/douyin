package model

import (
	"douyin_rpc/server/cmd/video/global"
	"github.com/bytedance/sonic"
	"strconv"
	"sync"

	"github.com/bwmarrin/snowflake"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type Video struct {
	gorm.Model
	ID       int64  `gorm:"primary_key;not null"`
	AuthorID int64  `gorm:"index; not null"`                          //作者id
	Title    string `gorm:"type:varchar(50); not null; default='';"`  //视频标题
	PlayURL  string `gorm:"type:varchar(300); not null; default='';"` //播放地址
	CoverURL string `gorm:"type:varchar(500); not null; default='';"` //封面地址
}

func (b *Video) BeforeCreate(_ *gorm.DB) (err error) {
	sf, err := snowflake.NewNode(3)
	if err != nil {
		klog.Fatalf("generate id failed: %s", err.Error())
	}
	b.ID = sf.Generate().Int64()
	return nil
}

var mutex sync.Mutex
var update_keys []string

const VideoCachePrefix string = "video:video_"

// queryVideoByID 查询数据库的video
func queryVideoByID(tx *gorm.DB, videoId int64) (*Video, error) {
	// 直接查询数据库
	var video Video
	if err := tx.Table("video").Where("id=?", videoId).Find(&video).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// QueryVideoByIDWithCache 通过视频id查询视频信息
func QueryVideoByIDWithCache(tx *gorm.DB, videoID int64) (*Video, error) {

	//return queryVideoByID(tx, videoID)
	key := VideoCachePrefix + "ID_" + strconv.Itoa(int(videoID))

	var err error
	var video *Video
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		video, err = queryVideoByID(tx, videoID)

		if err != nil {
			return "", err
		}
		marshal, err2 := sonic.Marshal(video)
		if err2 != nil {
			return "", err2
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	if video != nil {
		return video, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &video)
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

// queryVideoByAuthorID 通过videoID查询video
func queryVideoByAuthorID(tx *gorm.DB, authorID int64) ([]Video, error) {
	var videos []Video
	if err := tx.Table("video").Where("author_id=?", authorID).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// QueryVideoByAuthorIDWithCache 先查缓存再查数据库
func QueryVideoByAuthorIDWithCache(tx *gorm.DB, authorID int64) ([]Video, error) {
	//return queryVideoByAuthorID(tx, authorID)
	key := VideoCachePrefix + "AuthorID_" + strconv.Itoa(int(authorID))
	// 查看key是否存在

	//不存在

	var err error
	var videos []Video
	fetch, err := global.RocksCacheClient.Fetch(key, 1*time.Hour, func() (string, error) {
		videos, err = queryVideoByAuthorID(tx, authorID)

		if err != nil {
			return "", err
		}
		marshal, err2 := sonic.Marshal(videos)
		if err2 != nil {
			return "", err2
		}
		return string(marshal), nil
	})
	if err != nil {
		return nil, err
	}
	if videos != nil {
		return videos, nil
	}
	err = sonic.Unmarshal([]byte(fetch), &videos)
	if err != nil {
		return nil, err
	}

	return videos, nil

}

func CreateVideo(tx *gorm.DB, video *Video) error {
	if err := tx.Table("video").Create(&video).Error; err != nil {
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
	if err := tx.Table("video").Where("created_at<=?", latestTime).Limit(global.ServerConfig.FeedNumber).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

func FeedWithoutTime(tx *gorm.DB) ([]Video, error) {
	var videos []Video
	if err := tx.Table("video").Limit(global.ServerConfig.FeedNumber).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}
