package main

import (
	"context"
	"douyin_rpc/client/kitex_gen/comment"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	"douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/client/kitex_gen/video"
	"douyin_rpc/server/cmd/favorite/global"
	favorite "douyin_rpc/server/cmd/favorite/kitex_gen/favorite"
	"douyin_rpc/server/cmd/favorite/model"
	"errors"
	"sync"
)

// FavoriteServiceImpl implements the last service interface defined in the IDL.
type FavoriteServiceImpl struct{}

// FavoriteAction implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteAction(ctx context.Context, req *favorite.FavoriteActionRequest) (resp *favorite.FavoriteActionResponse, err error) {
	// 事务开始
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	//查询视频，判断视频是否存在
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
		return nil, errors.New("视频不存在")
	}
	// 3.3 查询favorite表获取信息，
	if req.ActionType == 1 {
		// 3.3.1 查看点赞是否存在，如果存在返回
		favorite1, err1 := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.AuthId, req.VideoId)
		if err1 != nil {
			tx.Rollback()
			return
		}
		if favorite1.ID == 0 {
			// 3.3.2 不存在，需要创建
			favorite1.UserID = req.AuthId
			favorite1.Exist = true
			favorite1.VideoID = req.VideoId
			err1 := model.UpdateOrCreateFavorite(tx, *favorite1)
			if err1 != nil {
				err = err1
				tx.Rollback()
				return
			}
		} else {
			// 3.3.4 赞现在存在
			err = errors.New("赞存在，无需点赞")
			tx.Rollback()
			return
		}
	} else {
		favorite1, err1 := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.AuthId, req.VideoId)
		if err1 != nil {
			err = err1
			return
		}

		if favorite1.ID == 0 {
			// 3.3.2 不存在，无法取消
			err = errors.New("赞不存在，无需取消")
			tx.Rollback()
			return
		} else {
			// 3.3.4 赞现在存在
			favorite1.Exist = false
			err1 := model.UpdateFavorite(tx, *favorite1)
			if err1 != nil {
				err = err1
				tx.Rollback()
				return
			}
		}
	}
	if err = tx.Commit().Error; err != nil {

		tx.Rollback()
		return
	}
	// 删除缓存
	//go func() {
	//	key1 := models.VideoCachePrefix + "ID_" + strconv.Itoa(int(video.ID))
	//	key2 := models.VideoCachePrefix + "AuthorID_" + strconv.Itoa(int(video.AuthorID))
	//	key3 := models.FavoriteCachePrefix + "UserID_" + strconv.Itoa(int(c.AuthUser.UserID)) + "_VideoID_" + strconv.Itoa(int(video.ID))
	//	for {
	//		cache.Delete([]string{key1, key2, key3})
	//		if err == nil {
	//			break
	//		}
	//	}
	//}()
	return
}

// FavoriteList implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) FavoriteList(ctx context.Context, req *favorite.FavoriteListRequest) (resp *favorite.FavoriteListResponse, err error) {
	tx := global.DB.Debug()
	favorites, err1 := model.QueryFavoriteByUserIDWithCache(tx, req.UserId)
	if err1 != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	resp.VideoList = make([]*favorite.Video, len(favorites))
	wg.Add(len(favorites))
	for j := 0; j < len(favorites); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 3.2、查询视频
			videos, err1 := global.VideoClient.GetVideo(ctx, &video.GetVideoRequest{
				VideoId:   favorites[i].VideoID,
				AuthorId:  0,
				QueryType: 1,
			})

			// 3.3 查询视频author信息\
			author, err1 := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    videos.Video[0].AuthorId,
				Username:  "",
				QueryType: 1,
			})

			if err1 != nil {
				err = err1
				return
			}
			// 3.4 查询视频作者与用户auth的关系

			resp.VideoList[i] = &favorite.Video{
				Author: &favorite.User{
					Id:            author.User.Id,
					Name:          author.User.Name,
					FollowCount:   0,
					FollowerCount: 0,
					IsFollow:      false,
				},
				Id:            videos.Video[0].Id,
				FavoriteCount: 0,
				CommentCount:  0,
				IsFavorite:    false,
				Title:         videos.Video[0].Title,
				PlayUrl:       videos.Video[0].PlayUrl,
				CoverUrl:      videos.Video[0].CoverUrl,
			}
			// 异步查询出关注数，粉丝数，是否关注视频作者，是否给该视频点赞，视频点赞数，视频评论数
			var wg1 sync.WaitGroup
			wg1.Add(5)

			go func(followCount, followerCount *int64) {
				// 查询FollowCount和FollowerCount
				defer wg.Done()
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
				count, err1 := model.CountFavoriteByVideoID(tx, favorites[i].VideoID)
				if err1 != nil {
					err = err1
					return
				}
				*favoriteCount = count
			}(&resp.VideoList[i].FavoriteCount)
			go func(commentCount *int64) {
				// 查询CommentCount
				defer wg.Done()
				count, err1 := global.CommentClient.GetCommentCount(ctx, &comment.GetCommentCountRequest{VideoId: favorites[i].VideoID})
				if err1 != nil {
					err = err1
					return
				}
				*commentCount = count.Count
			}(&resp.VideoList[i].CommentCount)
			go func(isFavorite *bool) {
				// 查询登录用户是否点赞该视频
				defer wg.Done()
				if req.AuthId != -1 {
					cache, err1 := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.AuthId, favorites[i].VideoID)
					if err1 != nil {

						err = err1
						return
					}
					if cache.ID != 0 {
						*isFavorite = true
					}
				}
			}(&resp.VideoList[i].IsFavorite)
			go func(isFollow *bool) {
				//查询登录用户是否关注该视频作者
				defer wg.Done()
				if req.AuthId != -1 {
					relation1, err1 := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
						Id:           0,
						UserId:       req.AuthId,
						TargetId:     author.User.Id,
						RelationType: 0,
						QueryType:    1,
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

		}()
	}
	wg.Wait()
	return
}

// GetFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavorite(ctx context.Context, req *favorite.GetFavoriteRequest) (resp *favorite.GetFavoriteResponse, err error) {
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		// 通过user_id和video_id查找
		cache, err1 := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.UserId, req.VideoId)
		if err1 != nil {
			err = err1
			return
		}
		resp.Favorites = make([]*favorite.Favorite1, 1)
		resp.Favorites[0] = &favorite.Favorite1{
			Id:        cache.ID,
			UserId:    cache.UserID,
			VideoId:   cache.VideoID,
			CreatedAt: cache.CreatedAt.Unix(),
			UpdatedAt: cache.UpdatedAt.Unix(),
		}
		return
	} else if req.QueryType == 2 {
		//通过user_id查找
		cache, err1 := model.QueryFavoriteByUserIDWithCache(tx, req.UserId)
		if err1 != nil {
			err = err1
			return
		}
		resp.Favorites = make([]*favorite.Favorite1, len(cache))
		var wg sync.WaitGroup
		wg.Add(len(cache))
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				defer wg.Done()
				resp.Favorites[i] = &favorite.Favorite1{
					Id:        cache[i].ID,
					UserId:    cache[i].UserID,
					VideoId:   cache[i].VideoID,
					CreatedAt: cache[i].CreatedAt.Unix(),
					UpdatedAt: cache[i].UpdatedAt.Unix(),
				}
			}()
		}
		wg.Wait()
		return
	} else if req.QueryType == 3 {
		// 通过video_id查找
		cache, err1 := model.QueryFavoriteByVideoIDWithCache(tx, req.VideoId)
		if err1 != nil {
			err = err1
			return
		}
		resp.Favorites = make([]*favorite.Favorite1, len(cache))
		var wg sync.WaitGroup
		wg.Add(len(cache))
		for j := 0; j < len(cache); j++ {
			i := j
			go func() {
				defer wg.Done()
				resp.Favorites[i] = &favorite.Favorite1{
					Id:        cache[i].ID,
					UserId:    cache[i].UserID,
					VideoId:   cache[i].VideoID,
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

// GetFavoriteCount implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavoriteCount(ctx context.Context, req *favorite.GetFavoriteCountRequest) (resp *favorite.GetFavoriteCountResponse, err error) {
	tx := global.DB.Debug()
	cache, err := model.CountFavoriteByVideoID(tx, req.VideoId)
	if err != nil {
		return nil, err
	}
	resp.Count = cache
	return
}
