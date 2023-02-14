package main

import (
	"context"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	"douyin_rpc/server/cmd/user/global"
	"douyin_rpc/server/cmd/user/kitex_gen/user"
	"douyin_rpc/server/cmd/user/model"
	tools "douyin_rpc/server/cmd/user/tool"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Login implements the UserServiceImpl interface.
// 返回用户的id即可，有api进行jwt解析
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	tx := global.DB.Debug()

	temp_user, err := model.QueryUserByNameWithCache(tx, req.Username)
	if err != nil {
		return nil, err
	}
	if temp_user.ID == 0 {
		err = errors.New("该用户不存在")
		return
	}
	// 4、获取到盐值，加密后判断是否一致
	if !tools.VerifyUserPassword(temp_user.Salt, req.Password, temp_user.Password) {
		return nil, err
	}
	return &user.LoginResponse{
		StatusCode: 0,
		StatusMsg:  "",
		UserId:     temp_user.ID,
	}, nil
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {

	salt := tools.Salt()
	cryptoPassword := tools.CryptUserPassword(req.Password, salt)
	// 4、创建用户
	temp_user := model.User{
		Model:    gorm.Model{},
		Name:     req.Username,
		Password: cryptoPassword,
		Salt:     salt,
	}
	// 5、gorm创建用户
	// 5.1、事务开始
	tx := global.DB.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	user1, err := model.QueryUserByNameWithCache(tx, req.Username)
	fmt.Println(user1)
	if err != nil {
		return nil, err
	} else if user1.ID != 0 {
		return nil, errors.New("该用户已存在")
	}
	// 5.2、创建用户
	err = model.CreateUser(tx, &temp_user)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// 5.3、事务提交
	if err = tx.Commit().Error; err != nil {

		tx.Rollback()

		return nil, err
	}
	go model.DeleteCache()
	return &user.RegisterResponse{
		StatusCode: 0,
		StatusMsg:  "",
		UserId:     temp_user.ID,
	}, nil
}

// GetUser implements the UserServiceImpl interface.
// 该接口用于rpc内部获得用户信息，供其他rpc服务调用
func (s *UserServiceImpl) GetUser(ctx context.Context, req *user.GetUserRequest) (resp *user.GetUserResponse, err error) {
	tx := global.DB.Debug()
	resp = new(user.GetUserResponse)
	var cache *model.User
	if req.QueryType == 1 {
		cache, err = model.QueryUserByIDWithCache(tx, req.UserId)
		if err != nil {
			return nil, err
		}
		resp.User = &user.User1{
			Id:        cache.ID,
			Name:      cache.Name,
			Password:  cache.Password,
			CreatedAt: cache.CreatedAt.Unix(),
			UpdatedAt: cache.UpdatedAt.Unix(),
			Salt:      cache.Salt,
		}
	} else if req.QueryType == 2 {
		cache, err = model.QueryUserByNameWithCache(tx, req.Username)
		if err != nil {
			return nil, err
		}
		resp.User = &user.User1{
			Id:        cache.ID,
			Name:      cache.Name,
			Password:  cache.Password,
			CreatedAt: cache.CreatedAt.Unix(),
			UpdatedAt: cache.UpdatedAt.Unix(),
			Salt:      cache.Salt,
		}
	}
	return
}

// User implements the UserServiceImpl interface.
func (s *UserServiceImpl) User(ctx context.Context, req *user.UserRequest) (resp *user.UserResponse, err error) {
	tx := global.DB.Debug()
	cache, err1 := model.QueryUserByIDWithCache(tx, req.UserId)
	if err1 != nil {
		err = err1
		return
	}
	resp = new(user.UserResponse)
	resp.User = &user.User{
		Id:            strconv.FormatInt(cache.ID, 10),
		Name:          cache.Name,
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
	}
	//count, err1 := global.RelationClient.GetCount(ctx, &relation2.GetCountRequest{UserId: req.UserId})
	//if err1 != nil {
	//	err = err1
	//	return
	//}
	//resp.User.FollowCount = count.FollowCount
	//resp.User.FollowerCount = count.FollowerCount
	var wg sync.WaitGroup
	wg.Add(2)
	go func(followCount, followerCount *int64) {
		defer wg.Done()
		count, err1 := global.RelationClient.GetCount(ctx, &relation2.GetCountRequest{UserId: req.UserId})
		if err1 != nil {
			err = err1
			return
		}
		*followCount = count.FollowCount
		*followerCount = count.FollowerCount
	}(&resp.User.FollowCount, &resp.User.FollowerCount)
	go func(isFollow *bool) {
		defer wg.Done()
		/*
		   query_type=1  根据id查询
		   query_type=2  根据user_id查询
		   query_type=3  根据target_id查询
		   query_type=4  根据user_id和target_id查询
		   relation_type暂时不需要，考虑拿掉
		*/
		relation, err1 := global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
			Id:           0,
			UserId:       req.AuthId,
			TargetId:     req.UserId,
			RelationType: -1,
			QueryType:    4,
		})
		if err1 != nil {
			err = err1
			return
		}
		if len(relation.Relations) == 1 && relation.Relations[0].Id != 0 {
			*isFollow = true
		}
	}(&resp.User.IsFollow)
	wg.Wait()

	return
}
