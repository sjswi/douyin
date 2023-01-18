package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/cache"
	"douyin/models"
	"douyin/utils"
	"douyin/vo"
	"errors"
	"strconv"
	"sync"
)

//TODO
// 关注，粉丝，好友列表机器相似，可以整合为一个函数

type RelationList struct {
	utils.Response
	UserList []vo.User `json:"user_list"`
}

type RelationAction struct {
	ToUserId   uint
	ActionType int
	AuthUser   *auth2.Auth
}

func RelationActionPost(toUserId uint, actionType int, auth *auth2.Auth) error {
	return (&RelationAction{
		ToUserId:   toUserId,
		ActionType: actionType,
		AuthUser:   auth,
	}).Do()
}
func (c *RelationAction) Do() error {
	if err := c.checkParam(); err != nil {
		return err
	}
	if err := c.action(); err != nil {
		return err
	}
	return nil
}
func (c *RelationAction) action() error {
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 3、查询数据库获取两个用户信息，使用for update加锁（用户一般都存在）

	user, err := models.QueryUserByIDWithCache(tx, c.AuthUser.UserID)
	if err != nil {
		return err
	}
	targetUser, err := models.QueryUserByIDWithCache(tx, c.ToUserId)
	if err != nil {
		return err
	}
	var relation1, relation2 *models.Relation

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		relation1, err = models.QueryRelationByUserIDAndTargetIDWithCache(tx, user.ID, targetUser.ID)

	}()
	go func() {
		defer wg.Done()
		relation2, err = models.QueryRelationByUserIDAndTargetIDWithCache(tx, targetUser.ID, user.ID)

	}()

	wg.Wait()
	if err != nil {
		return err
	}
	var errAction error
	wg.Add(1)
	go func() {
		defer wg.Done()
		if c.ActionType == 1 {
			user.FollowCount += 1
			targetUser.FollowerCount += 1
		} else {
			user.FollowCount -= 1
			targetUser.FollowerCount -= 1
		}
		errAction = models.UpdateUser(tx, *user)
		if errAction != nil {
			return
		}
		errAction = models.UpdateUser(tx, *targetUser)
		if errAction != nil {
			return
		}
	}()
	// 由于不能确定这两个关系同时存在，因此不要使用for update加锁（使用for update时确保索引存在。不存在会锁住表）
	// for update在数据存在时加的是行级锁，不存在加的是间隙锁。之后进行insert时容易形成死锁
	if c.ActionType == 1 {
		if relation1.ID == 0 {
			relation1.Exist = true
			relation1.Type = 1
			relation1.UserID = user.ID
			relation1.TargetID = targetUser.ID
			if relation2.ID != 0 {
				relation1.Type = 2
				relation2.Type = 2
				err = models.UpdateOrCreateRelation(tx, *relation2)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
			err := models.UpdateOrCreateRelation(tx, *relation1)
			if err != nil {
				tx.Rollback()
				return err
			}

		} else {
			return errors.New("已经关注了")
		}
	} else if c.ActionType == 2 {
		// 取消关注，数据不存在直接报错
		if relation1.ID == 0 {
			return errors.New("并未关注，无需取消")
		} else {
			//数据存在
			if relation1.Type == 1 {
				// Type为1，只需要将Exist改为false
				relation1.Exist = false
				err := models.UpdateRelation(tx, *relation1)
				if err != nil {
					tx.Rollback()
					return err
				}
			} else {
				// Type为2，修改relation2的Type为1
				relation1.Exist = false
				relation1.Type = 1
				err := models.UpdateRelation(tx, *relation1)
				if err != nil {
					tx.Rollback()
					return err
				}
				relation2.Type = 1
				err = models.UpdateRelation(tx, *relation2)
				if err != nil {
					tx.Rollback()
					return err
				}
			}

		}
	}
	//注意使用gorm有可能修改到零值的需要使用Save而不能使用updates
	// 5、修改用户的关注数和粉丝数
	wg.Wait()
	if errAction != nil {
		tx.Rollback()
		return errAction
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	go func() {
		key1 := models.UserCachePrefix + "ID_" + strconv.Itoa(int(user.ID))
		key2 := models.UserCachePrefix + "Name_" + user.Name
		key3 := models.UserCachePrefix + "ID_" + strconv.Itoa(int(targetUser.ID))
		key4 := models.UserCachePrefix + "Name_" + targetUser.Name
		key5 := models.RelationCachePrefix + "UserID_" + strconv.Itoa(int(user.ID)) + "_TargetID_" + strconv.Itoa(int(targetUser.ID))
		key6 := models.RelationCachePrefix + "UserID_" + strconv.Itoa(int(targetUser.ID)) + "_TargetID_" + strconv.Itoa(int(user.ID))
		for {
			err := cache.Delete([]string{key1, key2, key3, key4, key6, key5})
			if err == nil {
				break
			}
		}
	}()
	return nil
}
func (c *RelationAction) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

// 三个列表获取通用
func checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}
func getList(userId uint, userType int, authId uint) ([]vo.User, error) {
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

type RelationFollowListFlow struct {
	AuthUser *auth2.Auth
	UserId   uint
	UserList []vo.User `json:"user_list"`
}

func RelationFollowListGet(userId uint, auth *auth2.Auth) ([]vo.User, error) {
	return (&RelationFollowListFlow{
		AuthUser: auth,
		UserId:   userId,
		UserList: nil,
	}).Do()
}
func (c *RelationFollowListFlow) Do() ([]vo.User, error) {
	if err := checkParam(); err != nil {
		return nil, err
	}
	var err error
	if c.UserList, err = getList(c.UserId, 1, c.AuthUser.UserID); err != nil {
		return nil, err
	}
	return c.UserList, err
}

type RelationFollowerListFlow struct {
	AuthUser *auth2.Auth
	UserId   uint
	UserList []vo.User `json:"user_list"`
}

func RelationFollowerListGet(userId uint, auth *auth2.Auth) ([]vo.User, error) {
	return (&RelationFollowerListFlow{
		AuthUser: auth,
		UserId:   userId,
		UserList: nil,
	}).Do()
}
func (c *RelationFollowerListFlow) Do() ([]vo.User, error) {
	if err := checkParam(); err != nil {
		return nil, err
	}
	var err error
	if c.UserList, err = getList(c.UserId, 2, c.AuthUser.UserID); err != nil {
		return nil, err
	}
	return c.UserList, err
}

type RelationFriendListFlow struct {
	AuthUser *auth2.Auth
	UserId   uint
	UserList []vo.User `json:"user_list"`
}

func RelationFriendListGet(userId uint, auth *auth2.Auth) ([]vo.User, error) {
	return (&RelationFriendListFlow{
		AuthUser: auth,
		UserId:   userId,
		UserList: nil,
	}).Do()
}
func (c *RelationFriendListFlow) Do() ([]vo.User, error) {
	if err := checkParam(); err != nil {
		return nil, err
	}
	var err error
	if c.UserList, err = getList(c.UserId, 3, c.AuthUser.UserID); err != nil {
		return nil, err
	}
	return c.UserList, err
}
