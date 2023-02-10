package main

import (
	"context"
	"douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/client/kitex_gen/video"
	"douyin_rpc/server/cmd/comment/global"
	comment "douyin_rpc/server/cmd/comment/kitex_gen/comment"
	"douyin_rpc/server/cmd/comment/model"
	"gorm.io/gorm"
	"time"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	var err error
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()

	videoErr := make(chan error)
	videoAuthorId := -1

	video, err := global.VideoClient.GetVideo(ctx, &video.GetVideoRequest{
		VideoId:   0,
		AuthorId:  0,
		QueryType: 0,
	})
	if err != nil {
		videoErr <- err
		return
	}
	global.UserClient.User(ctx, user.UserRequest{
		UserId: "",
		AuthId: 0,
	})
	if req.ActionType == "1" {
		comment1 := model.Comment{
			Model:      gorm.Model{},
			Content:    req.CommentText,
			CreateTime: time.Now().Unix(),
			VideoID:    req.VideoId,
			UserID:     req.AuthId,
		}
		err = model.CreateComment(tx, comment1)
		if err != nil {
			return err
		}
		resp.Comment = &comment.Comment{
			Id: comment1.ID,
			User: &comment.User{
				Id:            c.AuthUser.UserID,
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

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {
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

// GetComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetComment(ctx context.Context, req *comment.GetCommentRequest) (resp *comment.GetCommentResponse, err error) {
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		cache, err := model.QueryCommentByIDWithCache(tx, req.Id)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, 1)
		// 设置resp.Comments
	} else if req.QueryType == 2 {
		cache, err := model.QueryCommentByUserIDAndVideoIDWithCache(tx, req.UserId, req.VideoId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, 1)
		// 设置resp.Comments
	} else if req.QueryType == 3 {
		cache, err := model.QueryCommentByUserIDWithCache(tx, req.UserId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, len(cache))
		// 设置resp.Comments
	} else if req.QueryType == 4 {
		cache, err := model.QueryCommentByVideoIDWithCache(tx, req.UserId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, len(cache))
		// 设置resp.Comments
	}

	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.GetCommentCountRequest) (resp *comment.GetCommentCountResponse, err error) {
	tx := global.DB.Debug()
	id, err := model.CountCommentByVideoID(tx, req.VideoId)
	resp.Count = *id
	return
}
