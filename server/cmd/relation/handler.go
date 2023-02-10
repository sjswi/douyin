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

	user, err := global.UserClient.GetUser(ctx, &user2.GetUserRequest{
		UserId:    req.AuthId,
		Username:  "",
		QueryType: 0,
	})
	if err != nil {
		return
	}
	targetUser, err := global.UserClient.GetUser(ctx, &user2.GetUserRequest{
		UserId:    req.ToUserId,
		Username:  "",
		QueryType: 0,
	})
	if err != nil {
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
	var errAction error
	wg.Add(1)
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
			err := model.UpdateOrCreateRelation(tx, *relation1)
			if err != nil {
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
				err := model.UpdateRelation(tx, *relation1)
				if err != nil {
					tx.Rollback()
					return
				}
			} else {
				// Type为2，修改relation2的Type为1
				relation1.Exist = false
				relation1.Type = 1
				err := model.UpdateRelation(tx, *relation1)
				if err != nil {
					tx.Rollback()
					return
				}
				relation2.Type = 1
				err = model.UpdateRelation(tx, *relation2)
				if err != nil {
					tx.Rollback()
					return
				}
			}

		}
	}
	//注意使用gorm有可能修改到零值的需要使用Save而不能使用updates
	// 5、修改用户的关注数和粉丝数
	wg.Wait()
	if errAction != nil {
		tx.Rollback()
		return
	}
	if err := tx.Commit().Error; err != nil {
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
func getList(userId int64, userType int, authId int64) ([]*relation.User, error) {
	var relations []models.Relation
	var err error
	tx := driver.Db.Debug()
	if userType == 1 {
		relations, err = models.QueryRelationByUserIDWithCache(tx, userId)
		if err != nil {
			return nil, err
		}
	} else if userType == 2 {
		relations, err = models.QueryRelationByTargetIDWithCache(tx, userId)
		if err != nil {
			return nil, err
		}
	} else {
		relations, err = models.QueryRelationIsFriend(tx, userId)
		if err != nil {
			return nil, err
		}
	}

	userList := make([]vo.User, len(relations))
	var wg sync.WaitGroup
	wg.Add(len(relations))
	for j := 0; j < len(relations); j++ {
		i := j
		go func() {
			defer wg.Done()
			var user *models.User
			var relation *models.Relation
			// 此处为TargetID

			if userType == 2 {
				userList[i].ID = relations[i].UserID
			} else {
				userList[i].ID = relations[i].TargetID
			}
			user, err = models.QueryUserByIDWithCache(tx, userList[i].ID)
			if err != nil {
				return
			}

			userList[i].Name = user.Name
			userList[i].FollowerCount = user.FollowerCount
			userList[i].FollowCount = user.FollowCount
			userList[i].IsFollow = false
			relation, err = models.QueryRelationByUserIDAndTargetIDWithCache(tx, authId, userList[i].ID)
			if err != nil {
				return
			}
			// 再次判断是否存在
			if relation.ID != 0 {
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
	list, err := getList(req.UserId, 1, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// FollowerList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FollowerList(ctx context.Context, req *relation.RelationFollowerListRequest) (resp *relation.RelationFollowerListResponse, err error) {
	list, err := getList(req.UserId, 2, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// FriendList implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) FriendList(ctx context.Context, req *relation.RelationFriendListRequest) (resp *relation.RelationFriendListResponse, err error) {
	list, err := getList(req.UserId, 3, req.AuthId)
	if err != nil {
		return nil, err
	}
	resp.UserList = list
	return
}

// GetRelation implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetRelation(ctx context.Context, req *relation.GetRelationRequest) (resp *relation.GetRelationResponse, err error) {
	tx := global.DB.Debug()
	if req.QueryType == 1 {
		model.QueryRelationByUserIDWithCache(tx)
	} else if req.QueryType == 2 {

	}
}

// GetCount implements the RelationServiceImpl interface.
func (s *RelationServiceImpl) GetCount(ctx context.Context, req *relation.GetCountRequest) (resp *relation.GetCountResponse, err error) {
	tx := global.DB.Debug()
	countRelation, i, i2, err := model.CountRelation(tx, req.UserId)
	if err != nil {
		return nil, err
	}
}
