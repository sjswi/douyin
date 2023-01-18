package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/vo"
	"sync"
	"time"
)

type FeedFlow struct {
	VideoList  []vo.Video `json:"video_list"`
	NextTime   int64      `json:"next_time"`
	LatestTime time.Time
	AuthUser   *auth2.Auth
}

func FeedGet(latestTime time.Time, auth *auth2.Auth) ([]vo.Video, int64, error) {
	return (&FeedFlow{
		VideoList:  nil,
		NextTime:   0,
		LatestTime: latestTime,
		AuthUser:   auth,
	}).Do()
}
func (c *FeedFlow) Do() ([]vo.Video, int64, error) {
	if err := c.checkParam(); err != nil {
		return nil, 0, err
	}
	if err := c.feed(); err != nil {
		return nil, 0, err
	}
	return c.VideoList, c.NextTime, nil
}
func (c *FeedFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *FeedFlow) feed() error {
	tx := driver.Db.Debug()
	videos, err := models.Feed(tx, c.LatestTime)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(len(videos))
	c.VideoList = make([]vo.Video, len(videos))
	var errFeed error
	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 查询每个视频的作者
			author, err := models.QueryUserByIDWithCache(tx, videos[i].AuthorID)
			if err != nil {
				errFeed = err
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
				ID:            videos[i].ID,
				FavoriteCount: videos[i].FavoriteCount,
				CommentCount:  videos[i].CommentCount,
				IsFavorite:    false,
				Title:         videos[i].Title,
				PlayURL:       videos[i].PlayURL,
				CoverURL:      videos[i].CoverURL,
			}
			if c.AuthUser != nil {
				favorite, err := models.QueryFavoriteByUserIDAndVideoIDWithCache(tx, author.ID, videos[i].ID)
				if err != nil {
					errFeed = err
					return
				}
				relation, err := models.QueryRelationByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, author.ID)
				if err != nil {
					errFeed = err
					return
				}
				if favorite.ID != 0 {
					c.VideoList[i].IsFavorite = true
				}
				if relation.ID != 0 {
					c.VideoList[i].Author.IsFollow = true
				}
			}
		}()

	}
	wg.Wait()
	if errFeed != nil {
		return err
	}
	if len(c.VideoList) > 0 {
		c.NextTime = videos[0].CreatedAt.UnixMilli()
	}
	return nil
}
