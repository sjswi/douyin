package main

import (
	"bytes"
	"context"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	user "douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/server/cmd/video/global"
	video "douyin_rpc/server/cmd/video/kitex_gen/video"
	"douyin_rpc/server/cmd/video/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
	"time"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	tx := global.DB.Debug()
	latestTime, err := time.Parse("", req.LatestTime)
	if err != nil {
		return nil, err
	}
	videos, err := model.Feed(tx, latestTime)
	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(videos))
	resp.VideoList = make([]*video.Video, len(videos))
	var errFeed error
	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 查询每个视频的作者，rpc调用user服务
			getUserResp, err := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    videos[i].AuthorID,
				Username:  "",
				QueryType: 0, // 0根据id查询，1根据名字查询
			})
			if err != nil {
				errFeed = err
				return
			}
			resp.VideoList[i] = &video.Video{
				Author: &video.User{
					Id:            getUserResp.User.Id,
					Name:          getUserResp.User.Name,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      false,
				},
				Id:            videos[i].ID,
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
				Title:         videos[i].Title,
				PlayUrl:       videos[i].PlayURL,
				CoverUrl:      videos[i].CoverURL,
			}
			var wg1 sync.WaitGroup
			wg1.Add(5)
			var tempErr error
			go func(followCount, followerCount *int64) {
				// 查询FollowCount和FollowerCount
				defer wg.Done()

			}(&resp.VideoList[i].Author.FollowCount, &resp.VideoList[i].Author.FollowerCount)
			go func(favoriteCount *int64) {
				// 查询FavoriteCount
				defer wg.Done()

			}(&resp.VideoList[i].FavoriteCount)
			go func(commentCount *int64) {
				// 查询CommentCount
				defer wg.Done()

			}(&resp.VideoList[i].CommentCount)
			go func(isFavorite *bool) {
				// 查询登录用户是否点赞该视频
				defer wg.Done()
				if req.AuthId != -1 {

				}
			}(&resp.VideoList[i].IsFavorite)
			go func(isFollow *bool) {
				//查询登录用户是否关注该视频作者
				defer wg.Done()
				if req.AuthId != -1 {

				}
			}(&resp.VideoList[i].Author.IsFollow)
			wg1.Wait()
			if tempErr != nil {
				err = tempErr
				return
			}
			//if c.AuthUser != nil {
			//	favorite, err := global.FavoriteClient.GetFavorite(tx, author.ID, videos[i].ID)
			//	if err != nil {
			//		errFeed = err
			//		return
			//	}
			//	getRelationResp, err := global.RelationClient.GetRelation(ctx, c.AuthUser.UserID, author.ID)
			//	if err != nil {
			//		errFeed = err
			//		return
			//	}
			//	if len(getRelationResp.Relations) == 1 {
			//		resp.VideoList[i].IsFavorite = true
			//	}
			//	if len(getRelationResp.Relations) == 1 {
			//		resp.VideoList[i].Author.IsFollow = true
			//	}
			//}
		}()

	}
	wg.Wait()
	if errFeed != nil {
		err = errFeed
		return
	}
	if len(resp.VideoList) > 0 {
		resp.NextTime = videos[0].CreatedAt.UnixMilli()
	}
	return
}

// PublishAction implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {
	uid := uuid.New().String()
	reader := bytes.NewReader(req.Data)
	videoURL := global.OSS.Put(uid+req.Filename, reader)
	//coverURL := storage.OSS.Put(uid+".jpeg", snapshot)
	coverURL := videoURL + "?x-oss-process=video/snapshot,t_7000,f_jpg,w_800,h_600,m_fast"
	videoModel := model.Video{
		Model:    gorm.Model{},
		AuthorID: req.AuthId,
		Title:    req.Title,
		PlayURL:  videoURL,
		CoverURL: coverURL,
	}
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	err = model.CreateVideo(tx, &videoModel)
	if err != nil {
		tx.Rollback()
		return
	}

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}
	return
}

// PublishList implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) PublishList(ctx context.Context, req *video.PublishListRequest) (resp *video.PublishListResponse, err error) {
	tx := global.DB.Debug()
	videos, err := model.QueryVideoByAuthorIDWithCache(tx, req.AuthId)
	if err != nil {
		return err
	}
	var errList error
	var wg sync.WaitGroup
	// 3.2、查询视频的作者，填充返回的视频信息
	resp.VideoList = make([]*video.Video, len(videos))
	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()

			resp.VideoList[i].Id = videos[i].ID
			resp.VideoList[i].Title = videos[i].Title
			resp.VideoList[i].CommentCount = 0
			resp.VideoList[i].CoverUrl = videos[i].CoverURL
			resp.VideoList[i].PlayUrl = videos[i].PlayURL
			resp.VideoList[i].FavoriteCount = 0
			resp.VideoList[i].IsFavorite = false
			author, err := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    0,
				Username:  "",
				QueryType: 0,
			})
			if err != nil {
				errList = err
				return
			}

			resp.VideoList[i] = &video.Video{
				Id: videos[i].ID,
				Author: &video.User{
					Id:            author.User.Id,
					Name:          author.User.Name,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      false,
				},
				PlayUrl:       videos[i].PlayURL,
				CoverUrl:      videos[i].CoverURL,
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
				Title:         videos[i].Title,
			}
			relation, err := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
				Id:           0,
				UserId:       0,
				TargetId:     0,
				RelationType: 0,
				QueryType:    0,
			})
			if err != nil {
				errList = err
				return
			}
			var wg1 sync.WaitGroup
			wg1.Add(5)
			var tempErr error
			go func(followCount, followerCount *int64) {
				// 查询FollowCount和FollowerCount
				defer wg.Done()

			}(&resp.VideoList[i].Author.FollowCount, &resp.VideoList[i].Author.FollowerCount)
			go func(favoriteCount *int64) {
				// 查询FavoriteCount
				defer wg.Done()

			}(&resp.VideoList[i].FavoriteCount)
			go func(commentCount *int64) {
				// 查询CommentCount
				defer wg.Done()

			}(&resp.VideoList[i].CommentCount)
			go func(isFavorite *bool) {
				// 查询登录用户是否点赞该视频
				defer wg.Done()
				if req.AuthId != -1 {

				}
			}(&resp.VideoList[i].IsFavorite)
			go func(isFollow *bool) {
				//查询登录用户是否关注该视频作者
				defer wg.Done()
				if req.AuthId != -1 {

				}
			}(&resp.VideoList[i].Author.IsFollow)
			wg1.Wait()
			if tempErr != nil {
				err = tempErr
				return
			}
		}()
	}
	wg.Wait()
	if errList != nil {
		err = errList

	}
	return
}

// GetVideo implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoRequest) (resp *video.GetVideoResponse, err error) {
	// 根据video_id 查询video
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		cache, err := model.QueryVideoByIDWithCache(tx, req.VideoId)
		if err != nil {
			return nil, err
		}
		resp.Video = make([]*video.Video1, 1)
		resp.Video[0] = &video.Video1{
			Id:        cache.ID,
			AuthorId:  cache.AuthorID,
			Title:     cache.Title,
			PlayUrl:   cache.PlayURL,
			CoverUrl:  cache.CoverURL,
			CreatedAt: cache.CreatedAt.Unix(),
			UpdatedAt: cache.UpdatedAt.Unix(),
		}
	} else if req.QueryType == 2 {
		cache, err := model.QueryVideoByAuthorIDWithCache(tx, req.AuthorId)
		if err != nil {
			return nil, err
		}
		resp.Video = make([]*video.Video1, len(cache))
		var wg sync.WaitGroup
		wg.Add(len(cache))
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				defer wg.Done()
				resp.Video[i] = &video.Video1{
					Id:        cache[i].ID,
					AuthorId:  cache[i].AuthorID,
					Title:     cache[i].Title,
					PlayUrl:   cache[i].PlayURL,
					CoverUrl:  cache[i].CoverURL,
					CreatedAt: cache[i].CreatedAt.Unix(),
					UpdatedAt: cache[i].UpdatedAt.Unix(),
				}
			}()
		}
		wg.Wait()
		return

	}
	return
}
