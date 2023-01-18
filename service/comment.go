package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/cache"
	"douyin/models"
	"douyin/utils"
	"douyin/vo"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

type CommentActionFlow struct {
	AuthUser    auth2.Auth
	VideoId     uint
	ActionType  int
	CommentText string
	CommentId   uint
	Comment     *vo.Comment `json:"comment"`
}

func CommentActionPost(videoID, commentID uint, actionType int, commentText string, auth auth2.Auth) (*vo.Comment, error) {
	return (&CommentActionFlow{
		VideoId:     videoID,
		ActionType:  actionType,
		CommentText: commentText,
		CommentId:   commentID,
		Comment:     nil,
		AuthUser:    auth,
	}).Do()
}
func (c *CommentActionFlow) Do() (*vo.Comment, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.comment(); err != nil {
		return nil, err
	}
	return c.Comment, nil
}
func (c *CommentActionFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *CommentActionFlow) comment() error {
	var err error
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	videoErr := make(chan error)
	videoAuthorId := -1
	go func() {
		video, err := models.QueryVideoByIDWithCache(tx, c.VideoId)
		if err != nil {
			videoErr <- err
			return
		}
		if c.ActionType == 1 {
			video.CommentCount += 1
		} else {
			video.CommentCount -= 1
		}
		// 更新video信息
		err = models.UpdateVideo(tx, video)
		if err != nil {
			videoErr <- err
			return
		}
		videoErr <- nil
		videoAuthorId = int(video.AuthorID)
		return
	}()

	if c.ActionType == 1 {
		comment := models.Comment{
			Model:      gorm.Model{},
			Content:    c.CommentText,
			CreateTime: time.Now().UTC(),
			VideoID:    c.VideoId,
			UserID:     c.AuthUser.UserID,
		}
		err = models.CreateComment(tx, comment)
		if err != nil {
			return err
		}
		c.Comment = &vo.Comment{
			ID: comment.ID,
			User: &vo.User{
				ID:            c.AuthUser.UserID,
				Name:          c.AuthUser.UserName,
				FollowCount:   c.AuthUser.FollowCount,
				FollowerCount: c.AuthUser.FollowerCount,
				IsFollow:      false,
			},
			Content:    comment.Content,
			CreateDate: utils.GetMonthAndDay(comment.CreateTime),
		}

	} else {
		err = models.DeleteCommentByID(tx, c.CommentId)
		if err != nil {
			return err
		}
		c.Comment = &vo.Comment{
			ID: c.CommentId,
			User: &vo.User{
				ID:            c.AuthUser.UserID,
				Name:          c.AuthUser.UserName,
				FollowCount:   c.AuthUser.FollowCount,
				FollowerCount: c.AuthUser.FollowerCount,
				IsFollow:      false,
			},
			Content:    c.CommentText,
			CreateDate: "",
		}

	}
	if err = <-videoErr; err != nil {
		return err

	}

	if err = tx.Commit().Error; err != nil {
		return err
	}
	//TODO
	// 事务提交后删除缓存
	go func() {
		key1 := models.CommentCachePrefix + "ID_" + strconv.Itoa(int(c.Comment.ID))
		key2 := models.CommentCachePrefix + "UserID_" + strconv.Itoa(int(c.Comment.User.ID))
		key3 := models.CommentCachePrefix + "VideoID_" + strconv.Itoa(int(c.VideoId))
		key4 := models.VideoCachePrefix + "AuthorID_" + strconv.Itoa(videoAuthorId)
		key5 := models.VideoCachePrefix + "ID_" + strconv.Itoa(int(c.VideoId))
		for {
			err := cache.Delete([]string{key1, key2, key3, key4, key5})
			if err == nil {
				break
			}
		}
	}()

	return nil
}

type CommentListFlow struct {
	VideoId     int
	AuthUser    auth2.Auth
	comments    []models.Comment
	CommentList []vo.Comment `json:"comment_list"`
}

func CommentListGet(videoId int, auth auth2.Auth) ([]vo.Comment, error) {
	return (&CommentListFlow{
		VideoId:     videoId,
		AuthUser:    auth,
		comments:    nil,
		CommentList: nil,
	}).Do()
}
func (c *CommentListFlow) Do() ([]vo.Comment, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.commentList(); err != nil {
		return nil, err
	}
	return c.CommentList, nil
}
func (c *CommentListFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *CommentListFlow) commentList() error {
	var errList error

	tx := driver.Db.Debug()
	var err error
	c.comments, err = models.QueryCommentByVideoIDWithCache(tx, uint(c.VideoId))
	if err != nil {
		errList = err
	}

	var wg sync.WaitGroup
	wg.Add(len(c.comments))
	// 4、装配返回值
	c.CommentList = make([]vo.Comment, len(c.comments))
	for j := 0; j < len(c.comments); j++ {
		i := j
		go func() {
			defer wg.Done()
			//var relation models.Relation
			//var user models.User
			c.CommentList[i].ID = c.comments[i].ID
			c.CommentList[i].Content = c.comments[i].Content
			c.CommentList[i].CreateDate = utils.GetMonthAndDay(c.comments[i].CreateTime)
			user, err := models.QueryUserByIDWithCache(tx, c.comments[i].UserID)
			if err != nil {
				errList = err
				return
			}
			c.CommentList[i].User = &vo.User{
				ID:            user.ID,
				Name:          user.Name,
				FollowCount:   user.FollowCount,
				FollowerCount: user.FollowerCount,
				IsFollow:      false,
			}
			relation, err := models.QueryRelationByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, user.ID)
			if err != nil {
				errList = err
				return
			}
			if relation.ID != 0 && relation.Exist {
				c.CommentList[i].User.IsFollow = true
			}
		}()

	}
	wg.Wait()
	if errList != nil {
		return errList
	}
	return nil
}
