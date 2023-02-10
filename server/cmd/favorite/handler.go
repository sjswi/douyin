package main

import (
	"context"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	"douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/client/kitex_gen/video"
	"douyin_rpc/server/cmd/favorite/global"
	favorite "douyin_rpc/server/cmd/favorite/kitex_gen/favorite"
	"douyin_rpc/server/cmd/favorite/model"
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
	var errAction error
	//查询视频，判断视频是否存在
	video, err := global.VideoClient.GetVideo(ctx, &video.GetVideoRequest{
		VideoId:   req.VideoId,
		AuthorId:  0,
		QueryType: 1,
	})
	if err != nil {
		return nil, err
	}
	// 3.3 查询favorite表获取信息，
	if req.ActionType == 1 {
		// 3.3.1 查看点赞是否存在，如果存在返回
		favorite, err := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.AuthId, req.VideoId)
		if err != nil {
			tx.Rollback()
			return
		}
		if favorite.ID == 0 {
			// 3.3.2 不存在，需要创建
			favorite.UserID = req.AuthId
			favorite.Exist = true
			favorite.VideoID = req.VideoId
			err := model.UpdateOrCreateFavorite(tx, *favorite)
			if err != nil {
				tx.Rollback()
				return
			}
		} else {
			// 3.3.4 赞现在存在
			tx.Rollback()
			return
		}
	} else {
		favorite, err := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.AuthId, req.VideoId)
		if err != nil {
			return
		}

		if favorite.ID == 0 {
			// 3.3.2 不存在，无法取消
			tx.Rollback()
			return
		} else {
			// 3.3.4 赞现在存在
			favorite.Exist = false
			err := model.UpdateFavorite(tx, *favorite)
			if err != nil {
				tx.Rollback()
				return
			}
		}
	}

	if errAction != nil {
		tx.Rollback()
		return
	}
	if err := tx.Commit().Error; err != nil {
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
	favorites, err := model.QueryFavoriteByUserIDWithCache(tx, req.UserId)
	if err != nil {
		return err
	}
	var errList error
	var wg sync.WaitGroup
	resp.VideoList = make([]*favorite.Video, len(favorites))
	wg.Add(len(favorites))
	for j := 0; j < len(favorites); j++ {
		i := j
		go func() {
			defer wg.Done()
			// 3.2、查询视频
			videos, err := global.VideoClient.GetVideo(ctx, &video.GetVideoRequest{
				VideoId:   req.VideoId,
				AuthorId:  0,
				QueryType: 1,
			})

			// 3.3 查询视频author信息\
			author, err := global.UserClient.GetUser(ctx, &user.GetUserRequest{
				UserId:    videos.Video[0].AuthorId,
				Username:  "",
				QueryType: 1,
			})

			if err != nil {
				errList = err
				return
			}
			// 3.4 查询视频作者与用户auth的关系
			relation, err := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
				Id:           0,
				UserId:       0,
				TargetId:     0,
				RelationType: 0,
				QueryType:    0,
			})
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
		return errList
	}
	return nil
}

// GetFavorite implements the FavoriteServiceImpl interface.
func (s *FavoriteServiceImpl) GetFavorite(ctx context.Context, req *favorite.GetFavoriteRequest) (resp *favorite.GetFavoriteResponse, err error) {
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		// 通过user_id和video_id查找
		cache, err := model.QueryFavoriteByUserIDAndVideoIDWithCache(tx, req.UserId, req.VideoId)
		if err != nil {
			return nil, err
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
		cache, err := model.QueryFavoriteByUserIDWithCache(tx, req.UserId)
		if err != nil {
			return nil, err
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
		cache, err := model.QueryFavoriteByVideoIDWithCache(tx, req.VideoId)
		if err != nil {
			return nil, err
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
