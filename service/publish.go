package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/storage"
	"douyin/vo"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
	"sync"
)

type PublishListFlow struct {
	VideoList []vo.Video `json:"video_list"`
	AuthorId  uint
	AuthUser  *auth2.Auth
}

func PublishListGet(authorId uint, auth *auth2.Auth) ([]vo.Video, error) {
	return (&PublishListFlow{
		VideoList: nil,
		AuthorId:  0,
		AuthUser:  nil,
	}).Do()
}
func (c *PublishListFlow) Do() ([]vo.Video, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.publishList(); err != nil {
		return nil, err
	}
	return c.VideoList, nil
}
func (c *PublishListFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *PublishListFlow) publishList() error {

	tx := driver.Db.Debug()
	videos, err := models.QueryVideoByAuthorIDWithCache(tx, c.AuthorId)
	if err != nil {
		return err
	}
	var errList error
	var wg sync.WaitGroup
	// 3.2、查询视频的作者，填充返回的视频信息
	c.VideoList = make([]vo.Video, len(videos))
	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()
			c.VideoList[i].ID = videos[i].ID
			c.VideoList[i].Title = videos[i].Title
			c.VideoList[i].CommentCount = videos[i].CommentCount
			c.VideoList[i].CoverURL = videos[i].CoverURL
			c.VideoList[i].PlayURL = videos[i].PlayURL
			c.VideoList[i].FavoriteCount = videos[i].FavoriteCount
			c.VideoList[i].IsFavorite = false
			author, err := models.QueryUserByIDWithCache(tx, videos[i].AuthorID)
			if err != nil {
				errList = err
				return
			}
			c.VideoList[i].Author = &vo.User{
				ID:            author.ID,
				Name:          author.Name,
				FollowCount:   author.FollowCount,
				FollowerCount: author.FollowerCount,
				IsFollow:      false,
			}
			relation, err := models.QueryRelationByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, author.ID)
			if err != nil {
				errList = err
				return
			}
			if relation.ID != 0 {
				c.VideoList[i].Author.IsFollow = true
			}
		}()
	}
	wg.Wait()
	if errList != nil {
		return errList
	}
	return nil
}

type PublishActionFlow struct {
	Data     *multipart.FileHeader
	Title    string
	AuthUser *auth2.Auth
}

func PublishActionPost(title string, video *multipart.FileHeader, auth *auth2.Auth) error {
	return (&PublishActionFlow{
		Data:     video,
		Title:    title,
		AuthUser: auth,
	}).Do()
}
func (c *PublishActionFlow) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	if err := c.publish(); err != nil {
		return err
	}
	return nil
}
func (c *PublishActionFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *PublishActionFlow) publish() error {
	uid := uuid.New().String()
	video, err := c.Data.Open()
	videoURL := storage.OSS.Put(uid+c.Data.Filename, video)
	//coverURL := storage.OSS.Put(uid+".jpeg", snapshot)
	coverURL := videoURL + "?x-oss-process=video/snapshot,t_7000,f_jpg,w_800,h_600,m_fast"
	videoModel := models.Video{
		Model:         gorm.Model{},
		AuthorID:      c.AuthUser.UserID,
		Title:         c.Title,
		CommentCount:  0,
		FavoriteCount: 0,
		PlayURL:       videoURL,
		CoverURL:      coverURL,
	}
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = models.CreateVideo(tx, &videoModel)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
