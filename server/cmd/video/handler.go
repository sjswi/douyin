package main

import (
	"context"
	"douyin_rpc/client/kitex_gen/comment"
	"douyin_rpc/client/kitex_gen/favorite"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	user "douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/server/cmd/video/global"
	video "douyin_rpc/server/cmd/video/kitex_gen/video"
	"douyin_rpc/server/cmd/video/model"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

// FeedServiceImpl implements the last service interface defined in the IDL.
type FeedServiceImpl struct{}

// Feed implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) Feed(ctx context.Context, req *video.FeedRequest) (resp *video.FeedResponse, err error) {
	tx := global.DB.Debug()
	var videos []model.Video
	if req.LatestTime == -1 {
		videos, err = model.FeedWithoutTime(tx)
		if err != nil {
			return nil, err
		}
	} else {
		latestTime := time.Unix(req.LatestTime, 0)
		if err != nil {
			return nil, err
		}
		videos, err = model.Feed(tx, latestTime)
	}

	if err != nil {
		return
	}
	var wg sync.WaitGroup
	wg.Add(len(videos))
	resp = new(video.FeedResponse)
	resp.VideoList = make([]*video.Video, len(videos))

	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 查询每个视频的作者，rpc调用user服务
			author, err1 := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    videos[i].AuthorID,
				Username:  "",
				QueryType: 1, // 0根据id查询，1根据名字查询
			})
			if err1 != nil {
				err = err1
				return
			}
			resp.VideoList[i] = &video.Video{
				Author: &video.User{
					Id:            strconv.FormatInt(author.User.Id, 10),
					Name:          author.User.Name,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      false,
				},
				Id:            strconv.FormatInt(videos[i].ID, 10),
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
				Title:         videos[i].Title,
				PlayUrl:       videos[i].PlayURL,
				CoverUrl:      videos[i].CoverURL,
			}
			// 异步查询出关注数，粉丝数，是否关注视频作者，是否给该视频点赞，视频点赞数，视频评论数
			var wg1 sync.WaitGroup
			wg1.Add(5)

			go func(followCount, followerCount *int64) {
				// 查询FollowCount和FollowerCount
				defer wg1.Done()
				count, err1 := global.RelationClient.GetCount(ctx, &relation2.GetCountRequest{UserId: author.User.Id})
				if err1 != nil {
					err = err1
					return
				}
				*followCount = count.FollowCount
				*followerCount = count.FollowerCount
			}(&resp.VideoList[i].Author.FollowCount, &resp.VideoList[i].Author.FollowerCount)
			go func(favoriteCount *int64) {
				// 查询FavoriteCount
				defer wg.Done()
				count, err1 := global.FavoriteClient.GetFavoriteCount(ctx, &favorite.GetFavoriteCountRequest{VideoId: videos[i].ID})
				if err1 != nil {
					err = err1
					return
				}
				*favoriteCount = count.Count
			}(&resp.VideoList[i].FavoriteCount)
			go func(commentCount *int64) {
				// 查询CommentCount
				defer wg1.Done()
				count, err1 := global.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountRequest{VideoId: videos[i].ID})
				if err1 != nil {
					err = err1
					return
				}
				*commentCount = count.Count
			}(&resp.VideoList[i].CommentCount)
			go func(isFavorite *bool) {
				// 查询登录用户是否点赞该视频
				defer wg1.Done()
				if req.AuthId != -1 {
					cache, err1 := global.FavoriteClient.GetFavorite(ctx, &favorite.GetFavoriteRequest{
						Id:        0,
						UserId:    req.AuthId,
						VideoId:   videos[i].ID,
						QueryType: 4,
					})
					if err1 != nil {

						err = err1
						return
					}
					if cache.Favorites[0].Id != 0 {
						*isFavorite = true
					}
				}
			}(&resp.VideoList[i].IsFavorite)
			go func(isFollow *bool) {
				//查询登录用户是否关注该视频作者
				defer wg1.Done()
				if req.AuthId != -1 {
					relation1, err1 := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
						Id:           0,
						UserId:       req.AuthId,
						TargetId:     author.User.Id,
						RelationType: 0,
						QueryType:    4,
					})
					if err1 != nil {

						err = err1
						return
					}
					if relation1.Relations[0].Id != 0 {
						*isFollow = true
					}
				}
			}(&resp.VideoList[i].Author.IsFollow)
			wg1.Wait()

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
	if err != nil {
		return
	}
	if len(resp.VideoList) > 0 {
		resp.NextTime = videos[0].CreatedAt.UnixMilli()
	}
	return
}

// PublishAction implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) PublishAction(ctx context.Context, req *video.PublishActionRequest) (resp *video.PublishActionResponse, err error) {

	videoModel := model.Video{
		Model:    gorm.Model{},
		AuthorID: req.AuthId,
		Title:    req.Title,
		PlayURL:  req.PlayUrl,
		CoverURL: req.CoverUrl,
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
	videos, err1 := model.QueryVideoByAuthorIDWithCache(tx, req.AuthId)
	if err1 != nil {
		err = err1
		return
	}
	resp = new(video.PublishListResponse)
	var wg sync.WaitGroup
	// 3.2、查询视频的作者，填充返回的视频信息
	wg.Add(len(videos))
	resp.VideoList = make([]*video.Video, len(videos))
	for j := 0; j < len(videos); j++ {
		i := j
		go func() {
			defer wg.Done()

			author, err1 := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    videos[i].AuthorID,
				Username:  "",
				QueryType: 1,
			})
			if err1 != nil {
				err = err1
				return
			}

			resp.VideoList[i] = &video.Video{
				Id: strconv.FormatInt(videos[i].ID, 10),
				Author: &video.User{
					Id:            strconv.FormatInt(author.User.Id, 10),
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

			var wg1 sync.WaitGroup
			wg1.Add(5)

			go func(followCount, followerCount *int64) {
				// 查询FollowCount和FollowerCount
				defer wg1.Done()
				count, err1 := global.RelationClient.GetCount(ctx, &relation2.GetCountRequest{UserId: author.User.Id})
				if err1 != nil {
					err = err1
					return
				}
				*followCount = count.FollowCount
				*followerCount = count.FollowerCount
			}(&resp.VideoList[i].Author.FollowCount, &resp.VideoList[i].Author.FollowerCount)
			go func(favoriteCount *int64) {
				// 查询FavoriteCount
				defer wg1.Done()
				count, err1 := global.FavoriteClient.GetFavoriteCount(ctx, &favorite.GetFavoriteCountRequest{VideoId: videos[i].ID})
				if err1 != nil {
					err = err1
					return
				}
				*favoriteCount = count.Count
			}(&resp.VideoList[i].FavoriteCount)
			go func(commentCount *int64) {
				// 查询CommentCount
				defer wg1.Done()
				count, err1 := global.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountRequest{VideoId: videos[i].ID})
				if err1 != nil {
					err = err1
					return
				}
				*commentCount = count.Count
			}(&resp.VideoList[i].CommentCount)
			go func(isFavorite *bool) {
				// 查询登录用户是否点赞该视频
				defer wg1.Done()
				// 需要登录
				if req.AuthId != -1 {
					getFavorite, err1 := global.FavoriteClient.GetFavorite(ctx, &favorite.GetFavoriteRequest{
						Id:        0,
						UserId:    req.AuthId,
						VideoId:   videos[i].ID,
						QueryType: 4,
					})
					if err1 != nil {

						err = err1
						return
					}
					if len(getFavorite.Favorites) == 1 && getFavorite.Favorites[0].Id != 0 {
						*isFavorite = true
					}

				}
			}(&resp.VideoList[i].IsFavorite)
			go func(isFollow *bool) {
				//查询登录用户是否关注该视频作者
				defer wg1.Done()
				if req.AuthId != -1 {
					getRelation, err1 := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
						Id:           0,
						UserId:       req.AuthId,
						TargetId:     author.User.Id,
						RelationType: 0,
						QueryType:    4,
					})
					if err1 != nil {
						err = err1
						return
					}
					if len(getRelation.Relations) == 1 && getRelation.Relations[0].Id != 0 {
						*isFollow = true
					}
				}
			}(&resp.VideoList[i].Author.IsFollow)
			wg1.Wait()

		}()
	}
	wg.Wait()

	return
}

// GetVideo implements the FeedServiceImpl interface.
func (s *FeedServiceImpl) GetVideo(ctx context.Context, req *video.GetVideoRequest) (resp *video.GetVideoResponse, err error) {
	// 根据video_id 查询video
	tx := global.DB.Debug()
	/**
	* query_type=1 根据视频id查询
	* query_type=2 根据作者id查询
	*
	**/
	resp = new(video.GetVideoResponse)
	if req.QueryType == 1 {
		cache, err1 := model.QueryVideoByIDWithCache(tx, req.VideoId)
		if err1 != nil {
			err = err1
			return
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
		cache, err1 := model.QueryVideoByAuthorIDWithCache(tx, req.AuthorId)
		if err1 != nil {
			err = err1
			return
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
