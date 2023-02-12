package main

import (
	"context"
	"douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/client/kitex_gen/video"
	"douyin_rpc/server/cmd/comment/global"
	comment "douyin_rpc/server/cmd/comment/kitex_gen/comment"
	"douyin_rpc/server/cmd/comment/model"
	"douyin_rpc/server/cmd/comment/tools"
	"errors"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

// CommentServiceImpl implements the last service interface defined in the IDL.
type CommentServiceImpl struct{}

// CommentAction implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentAction(ctx context.Context, req *comment.CommentActionRequest) (resp *comment.CommentActionResponse, err error) {
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil || err != nil {
			tx.Rollback()
		}
	}()
	resp = new(comment.CommentActionResponse)
	//videoErr := make(chan error)
	//videoAuthorId := -1

	video1, err1 := global.VideoClient.GetVideo(ctx, &video.GetVideoRequest{
		VideoId:   req.VideoId,
		AuthorId:  0,
		QueryType: 1,
	})
	if err1 != nil {
		err = err1
		return
	}
	if video1.Video[0].Id == 0 {
		err = errors.New("视频不存在")
		return
	}
	user1, err := global.UserClient.User(ctx, &user.UserRequest{
		UserId: req.AuthId,
		AuthId: req.AuthId,
	})
	if err != nil {
		return nil, err
	}
	if req.ActionType == 1 {
		comment1 := model.Comment{
			Model:      gorm.Model{},
			Content:    req.CommentText,
			CreateTime: time.Now().Unix(),
			VideoID:    req.VideoId,
			UserID:     req.AuthId,
		}
		err = model.CreateComment(tx, &comment1)
		if err != nil {
			return
		}
		resp.Comment = &comment.Comment{
			Id: strconv.FormatInt(comment1.ID, 10),
			User: &comment.User{
				Id:            user1.User.Id,
				Name:          user1.User.Name,
				FollowCount:   user1.User.FollowCount,
				FollowerCount: user1.User.FollowerCount,
				IsFollow:      false,
			},
			Content:    comment1.Content,
			CreateDate: tools.GetMonthAndDay(comment1.CreatedAt),
		}

	} else {
		err = model.DeleteCommentByID(tx, req.CommentId)
		if err != nil {
			return
		}
		resp.Comment = &comment.Comment{
			Id: strconv.FormatInt(req.CommentId, 10),
			User: &comment.User{
				Id:            user1.User.Id,
				Name:          user1.User.Name,
				FollowCount:   user1.User.FollowCount,
				FollowerCount: user1.User.FollowerCount,
				IsFollow:      false,
			},
			Content:    req.CommentText,
			CreateDate: "",
		}

	}

	if err = tx.Commit().Error; err != nil {

		return
	}
	//TODO
	// 事务提交后删除缓存
	//go func() {
	//	key1 := models.CommentCachePrefix + "ID_" + strconv.Itoa(int(c.Comment.ID))
	//	key2 := models.CommentCachePrefix + "UserID_" + strconv.Itoa(int(c.Comment.User.ID))
	//	key3 := models.CommentCachePrefix + "VideoID_" + strconv.Itoa(int(c.VideoId))
	//	key4 := models.VideoCachePrefix + "AuthorID_" + strconv.Itoa(videoAuthorId)
	//	key5 := models.VideoCachePrefix + "ID_" + strconv.Itoa(int(c.VideoId))
	//	for {
	//		err := cache.Delete([]string{key1, key2, key3, key4, key5})
	//		if err == nil {
	//			break
	//		}
	//	}
	//}()

	return
}

// CommentList implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) CommentList(ctx context.Context, req *comment.CommentListRequest) (resp *comment.CommentListResponse, err error) {

	tx := global.DB.Debug()
	comments, err1 := model.QueryCommentByVideoIDWithCache(tx, req.VideoId)
	if err1 != nil {
		err = err1
		return
	}
	resp = new(comment.CommentListResponse)
	var wg sync.WaitGroup
	wg.Add(len(comments))
	// 4、装配返回值
	resp.CommentList = make([]*comment.Comment, len(comments))
	for j := 0; j < len(comments); j++ {
		i := j
		go func() {
			defer wg.Done()
			//var relation models.Relation
			//var user models.User
			resp.CommentList[i] = &comment.Comment{
				Id:         strconv.FormatInt(comments[i].ID, 10),
				User:       nil,
				Content:    comments[i].Content,
				CreateDate: tools.GetMonthAndDay(comments[i].CreatedAt),
			}
			user1, err1 := global.UserClient.User(ctx, &user.UserRequest{
				UserId: comments[i].UserID,
				AuthId: req.AuthId,
			})
			if err1 != nil {
				err = err1
				return
			}
			resp.CommentList[i].User = &comment.User{
				Id:            user1.User.Id,
				Name:          user1.User.Name,
				FollowCount:   user1.User.FollowCount,
				FollowerCount: user1.User.FollowerCount,
				IsFollow:      user1.User.IsFollow,
			}
		}()

	}
	wg.Wait()

	return
}

// GetComment implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetComment(ctx context.Context, req *comment.GetCommentRequest) (resp *comment.GetCommentResponse, err error) {
	resp = new(comment.GetCommentResponse)
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		//  通过comment_id查询，几乎不会用到。只返回一个评论，因为评论id是不重复的
		cache, err := model.QueryCommentByIDWithCache(tx, req.Id)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, 1)
		// 设置resp.Comments
		resp.Comments[0] = &comment.Comment1{
			Id:        cache.ID,
			UserId:    cache.UserID,
			VideoId:   cache.VideoID,
			CreatedAt: cache.CreatedAt.Unix(),
			UpdatedAt: cache.UpdatedAt.Unix(),
			Content:   cache.Content,
		}
	} else if req.QueryType == 2 {
		// 通过user_id和video_id获得评论
		cache, err := model.QueryCommentByUserIDAndVideoIDWithCache(tx, req.UserId, req.VideoId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, len(cache))

		// 设置resp.Comments
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				resp.Comments[i] = &comment.Comment1{
					Id:        cache[i].ID,
					UserId:    cache[i].UserID,
					VideoId:   cache[i].VideoID,
					CreatedAt: cache[i].CreatedAt.Unix(),
					UpdatedAt: cache[i].UpdatedAt.Unix(),
					Content:   cache[i].Content,
				}
			}()
		}
	} else if req.QueryType == 3 {
		// 通过user_id查询评论内容
		cache, err := model.QueryCommentByUserIDWithCache(tx, req.UserId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, len(cache))
		// 设置resp.Comments
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				resp.Comments[i] = &comment.Comment1{
					Id:        cache[i].ID,
					UserId:    cache[i].UserID,
					VideoId:   cache[i].VideoID,
					CreatedAt: cache[i].CreatedAt.Unix(),
					UpdatedAt: cache[i].UpdatedAt.Unix(),
					Content:   cache[i].Content,
				}
			}()
		}
	} else if req.QueryType == 4 {
		// 通过video_id查询评论内容
		cache, err := model.QueryCommentByVideoIDWithCache(tx, req.VideoId)
		if err != nil {
			return nil, err
		}
		resp.Comments = make([]*comment.Comment1, len(cache))
		// 设置resp.Comments
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				resp.Comments[i] = &comment.Comment1{
					Id:        cache[i].ID,
					UserId:    cache[i].UserID,
					VideoId:   cache[i].VideoID,
					CreatedAt: cache[i].CreatedAt.Unix(),
					UpdatedAt: cache[i].UpdatedAt.Unix(),
					Content:   cache[i].Content,
				}
			}()
		}
	}

	return
}

// GetCommentCount implements the CommentServiceImpl interface.
func (s *CommentServiceImpl) GetCommentCount(ctx context.Context, req *comment.GetCommentCountRequest) (resp *comment.GetCommentCountResponse, err error) {
	tx := global.DB.Debug()
	resp = new(comment.GetCommentCountResponse)
	id, err := model.CountCommentByVideoID(tx, req.VideoId)
	resp.Count = id
	return
}
