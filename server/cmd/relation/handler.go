package main

import (
	"context"
	user2 "douyin_rpc/client/kitex_gen/user"
	"douyin_rpc/server/cmd/relation/global"
	relation "douyin_rpc/server/cmd/relation/kitex_gen/relation"
	"douyin_rpc/server/cmd/relation/model"
	"sync"
)

// RelationServiceImpl implements the last service interface defined in the IDL.
type RelationServiceImpl struct{}

// Action implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) Action(ctx context.Context, req *relation.RelationActionRequest) (resp *relation.RelationActionResponse, err error) {
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 3、查询数据库获取两个用户信息，使用for update加锁（用户一般都存在）

	user, err1 := global.UserClient.GetUser(ctx, &user2.GetUserRequest{
		UserId:    req.AuthId,
		Username:  "",
		QueryType: 0,
	})
	if err1 != nil {
		err = err1
		return
	}
	targetUser, err := global.UserClient.GetUser(ctx, &user2.GetUserRequest{
		UserId:    req.ToUserId,
		Username:  "",
		QueryType: 0,
	})
	if err1 != nil {
		err = err1
		return
	}
	var relation1, relation2 *model.Relation

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		relation1, err = model.QueryRelationByUserIDAndTargetIDWithCache(tx, user.User.Id, targetUser.User.Id)

	}()
	go func() {
		defer wg.Done()
		relation2, err = model.QueryRelationByUserIDAndTargetIDWithCache(tx, targetUser.User.Id, user.User.Id)

	}()

	wg.Wait()
	if err != nil {
		return
	}

	// 由于不能确定这两个关系同时存在，因此不要使用for update加锁（使用for update时确保索引存在。不存在会锁住表）
	// for update在数据存在时加的是行级锁，不存在加的是间隙锁。之后进行insert时容易形成死锁
	if req.ActionType == 1 {
		if relation1.ID == 0 {
			relation1.Exist = true
			relation1.Type = 1
			relation1.UserID = user.User.Id
			relation1.TargetID = targetUser.User.Id
			if relation2.ID != 0 {
				relation1.Type = 2
				relation2.Type = 2
				err = model.UpdateOrCreateRelation(tx, *relation2)
				if err != nil {
					tx.Rollback()
					return
				}
			}
			err1 := model.UpdateOrCreateRelation(tx, *relation1)
			if err1 != nil {
				err = err1
				tx.Rollback()
				return
			}

		} else {
			return
		}
	} else if req.ActionType == 2 {
		// 取消关注，数据不存在直接报错
		if relation1.ID == 0 {
			return
		} else {
			//数据存在
			if relation1.Type == 1 {
				// Type为1，只需要将Exist改为false
				relation1.Exist = false
				err1 := model.UpdateRelation(tx, *relation1)
				if err1 != nil {
					tx.Rollback()
					err = err1
					return
				}
			} else {
				// Type为2，修改relation2的Type为1
				relation1.Exist = false
				relation1.Type = 1
				err1 := model.UpdateRelation(tx, *relation1)
				if err1 != nil {
					tx.Rollback()
					err = err1
					return
				}
				relation2.Type = 1
				err1 = model.UpdateRelation(tx, *relation2)
				if err1 != nil {
					tx.Rollback()
					err = err1
					return
				}
			}

		}
	}
	//注意使用gorm有可能修改到零值的需要使用Save而不能使用updates
	// 5、修改用户的关注数和粉丝数

	if err = tx.Commit().Error; err != nil {
		tx.Rollback()
		return
	}
	//go func() {
	//	key1 := models.UserCachePrefix + "ID_" + strconv.Itoa(int(user.ID))
	//	key2 := models.UserCachePrefix + "Name_" + user.Name
	//	key3 := models.UserCachePrefix + "ID_" + strconv.Itoa(int(targetUser.ID))
	//	key4 := models.UserCachePrefix + "Name_" + targetUser.Name
	//	key5 := models.RelationCachePrefix + "UserID_" + strconv.Itoa(int(user.ID)) + "_TargetID_" + strconv.Itoa(int(targetUser.ID))
	//	key6 := models.RelationCachePrefix + "UserID_" + strconv.Itoa(int(targetUser.ID)) + "_TargetID_" + strconv.Itoa(int(user.ID))
	//	for {
	//		err := cache.Delete([]string{key1, key2, key3, key4, key6, key5})
	//		if err == nil {
	//			break
	//		}
	//	}
	//}()
	return
}
func getList(ctx context.Context, userId int64, userType int, authId int64) ([]*relation.User, error) {
	var relations []model.Relation
	var err error
	tx := global.DB.Debug()
	if userType == 1 {
		relations, err = model.QueryRelationByUserIDWithCache(tx, userId)
		if err != nil {
			return nil, err
		}
	} else if userType == 2 {
		relations, err = model.QueryRelationByTargetIDWithCache(tx, userId)
		if err != nil {
			return nil, err
		}
	} else {
		relations, err = model.QueryRelationIsFriend(tx, userId)
		if err != nil {
			return nil, err
		}
	}

	userList := make([]*relation.User, len(relations))
	var wg sync.WaitGroup
	wg.Add(len(relations))
	for j := 0; j < len(relations); j++ {
		i := j
		go func() {
			defer wg.Done()

			// 此处为TargetID

			if userType == 2 {
				userList[i].Id = relations[i].UserID
			} else {
				userList[i].Id = relations[i].TargetID
			}
			user1, err1 := global.UserClient.User(ctx, &user2.UserRequest{
				UserId: userList[i].Id,
				AuthId: authId,
			})
			if err1 != nil {
				err = err1
				return
			}

			userList[i].Name = user1.User.Name
			userList[i].FollowerCount = user1.User.FollowerCount
			userList[i].FollowCount = user1.User.FollowCount
			userList[i].IsFollow = false
			relation1, err1 := model.QueryRelationByUserIDAndTargetIDWithCache(tx, authId, userList[i].Id)
			if err1 != nil {
				err = err1
				return
			}
			// 再次判断是否存在
			if relation1.ID != 0 {
				userList[i].IsFollow = true
			}
		}()

	}
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return userList, nil

}

// FollowList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowList(ctx context.Context, req *relation.RelationFollowListRequest) (resp *relation.RelationFollowListResponse, err error) {
	list, err := getList(ctx, req.UserId, 1, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	list, err := getList(ctx, req.UserId, 2, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	list, err := getList(ctx, req.UserId, 3, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// GetRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetRelation(ctx context.Context, req *relation.GetRelationRequest) (resp *relation.GetRelationResponse, err error) {
	/*
	   query_type=1  根据id查询
	   query_type=2  根据user_id查询
	   query_type=3  根据target_id查询
	   query_type=4  根据user_id和target_id查询
	   relation_type暂时不需要，考虑拿掉
	*/
	tx := global.DB.Debug()
	var cache []model.Relation
	if req.QueryType == 1 {
		//model.Q(tx, req.)
	} else if req.QueryType == 2 {
		cache, err = model.QueryRelationByUserIDWithCache(tx, req.UserId)
		if err != nil {
			return
		}
	} else if req.QueryType == 3 {
		cache, err = model.QueryRelationByTargetIDWithCache(tx, req.TargetId)
		if err != nil {
			return
		}
	} else if req.QueryType == 4 {
		cache1, err1 := model.QueryRelationByUserIDAndTargetIDWithCache(tx, req.UserId, req.TargetId)
		if err1 != nil {
			err = err1
			return
		}
		resp.Relations = make([]*relation.Relation1, 1)
		resp.Relations[0] = &relation.Relation1{
			Id:       cache1.ID,
			UserId:   cache1.UserID,
			TargetId: cache1.TargetID,
			Type:     int32(cache1.Type),
		}
		return
	}
	resp.Relations = make([]*relation.Relation1, len(cache))
	var wg sync.WaitGroup
	wg.Add(len(cache))
	for j := 0; j < len(cache); j++ {
		i := j
		go func() {
			resp.Relations[0] = &relation.Relation1{
				Id:       cache[i].ID,
				UserId:   cache[i].UserID,
				TargetId: cache[i].TargetID,
				Type:     int32(cache[i].Type),
			}
		}()
	}
	wg.Wait()
	return
}

// GetCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetCount(ctx context.Context, req *relation.GetCountRequest) (resp *relation.GetCountResponse, err error) {
	tx := global.DB.Debug()
	friendCount, followCount, followerCount, err := model.CountRelation(tx, req.UserId)
	if err != nil {
		return nil, err
	}
	resp.FriendCount = friendCount
	resp.FollowCount = followCount
	resp.FollowerCount = followerCount
	return
}
