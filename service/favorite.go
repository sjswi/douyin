package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/cache"
	"douyin/models"
	"douyin/vo"
	"errors"
	"strconv"
	"sync"
)

type FavoriteActionFlow struct {
	VideoId    uint
	AuthUser   auth2.Auth
	ActionType int
}

func FavoriteActionPost(videoId uint, action int, auth auth2.Auth) error {
	return (&FavoriteActionFlow{
		VideoId:    videoId,
		AuthUser:   auth,
		ActionType: action,
	}).Do()
}
func (c *FavoriteActionFlow) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	if err := c.favorite(); err != nil {
		return err
	}
	return nil
}
func (c *FavoriteActionFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *FavoriteActionFlow) favorite() error {
	// 事务开始
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	var errAction error
	// 3.2 查询数据库获取视频信息
	video, err := models.QueryVideoByIDWithCache(tx, uint(c.VideoId))
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if c.ActionType == 1 {
			video.FavoriteCount += 1
		} else {
			video.FavoriteCount -= 1
		}

		err := models.UpdateVideo(tx, video)
		if err != nil {
			errAction = err
		}
	}()
	// 3.3 查询favorite表获取信息，
	if c.ActionType == 1 {
		// 3.3.1 查看点赞是否存在，如果存在返回
		favorite, err := models.QueryFavoriteByUserIDAndVideoIDWithCache(tx, c.AuthUser.UserID, uint(c.VideoId))
		if err != nil {
			tx.Rollback()
			return err
		}
		if favorite.ID == 0 {
			// 3.3.2 不存在，需要创建
			favorite.UserID = c.AuthUser.UserID
			favorite.Exist = true
			favorite.VideoID = video.ID
			err := models.UpdateOrCreateFavorite(tx, *favorite)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 3.3.4 赞现在存在
			tx.Rollback()
			return errors.New("赞存在")
		}
	} else {
		favorite, err := models.QueryFavoriteByUserIDAndVideoIDWithCache(tx, c.AuthUser.UserID, uint(c.VideoId))
		if err != nil {
			return err
		}

		if favorite.ID == 0 {
			// 3.3.2 不存在，无法取消
			tx.Rollback()
			return errors.New("赞不存在")

		} else {
			// 3.3.4 赞现在存在
			favorite.Exist = false
			err := models.UpdateFavorite(tx, *favorite)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	wg.Wait()
	if errAction != nil {
		tx.Rollback()
		return errAction
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	// 删除缓存
	go func() {
		key1 := models.VideoCachePrefix + "ID_" + strconv.Itoa(int(video.ID))
		key2 := models.VideoCachePrefix + "AuthorID_" + strconv.Itoa(int(video.AuthorID))
		key3 := models.FavoriteCachePrefix + "UserID_" + strconv.Itoa(int(c.AuthUser.UserID)) + "_VideoID_" + strconv.Itoa(int(video.ID))
		for {
			cache.Delete([]string{key1, key2, key3})
			if err == nil {
				break
			}
		}
	}()
	return nil
}
func FavoriteListGet(userId uint, auth auth2.Auth) ([]vo.Video, error) {
	return (&FavoriteListFlow{
		UserId:   userId,
		AuthUser: auth,
	}).Do()
}

type FavoriteListFlow struct {
	UserId    uint
	AuthUser  auth2.Auth
	VideoList []vo.Video `json:"video_list"`
}

func (c *FavoriteListFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}
func (c *FavoriteListFlow) Do() ([]vo.Video, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.favoriteList(); err != nil {
		return nil, err
	}
	return c.VideoList, nil
}

func (c *FavoriteListFlow) favoriteList() error {
	//var favoriteList []models.Favorite
	tx := driver.Db.Debug()
	favorites, err := models.QueryFavoriteByUserIDWithCache(tx, c.UserId)
	if err != nil {
		return err
	}
	var errList error
	var wg sync.WaitGroup
	c.VideoList = make([]vo.Video, len(favorites))
	wg.Add(len(favorites))
	for j := 0; j < len(favorites); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 3.2、查询视频
			video, err := models.QueryVideoByIDWithCache(tx, favorites[i].VideoID)
			if err != nil {
				errList = err
				return
			}

			// 3.3 查询视频author信息\
			author, err := models.QueryUserByIDWithCache(tx, video.AuthorID)
			if err != nil {
				errList = err
				return
			}
			// 3.4 查询视频作者与用户auth的关系
			relation, err := models.QueryRelationByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, author.ID)
			if err != nil {
				errList = err
				return
			}
			c.VideoList[i] = vo.Video{
				Author: &vo.User{
					ID:            author.ID,
					Name:          author.Name,
					FollowCount:   author.FollowCount,
					FollowerCount: author.FollowerCount,
					IsFollow:      false,
				},
				ID:            video.ID,
				FavoriteCount: video.FavoriteCount,
				CommentCount:  video.CommentCount,
				IsFavorite:    false,
				Title:         video.Title,
				PlayURL:       video.PlayURL,
				CoverURL:      video.CoverURL,
			}
			if relation.ID != 0 && relation.Exist {
				c.VideoList[i].Author.IsFollow = true
			}
			// 3.5 查询视频自己是否点过赞
			favorite, err := models.QueryFavoriteByUserIDAndVideoIDWithCache(tx, c.AuthUser.UserID, video.ID)
			if err != nil {
				errList = err
				return
			}
			if favorite.ID != 0 {
				c.VideoList[i].IsFavorite = true
			}
		}()
	}
	wg.Wait()
	if errList != nil {
		return errList
	}
	return nil
}
