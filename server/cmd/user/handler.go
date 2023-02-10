package main

import (
	"context"
	relation2 "douyin_rpc/client/kitex_gen/relation"
	"douyin_rpc/server/cmd/user/global"
	user "douyin_rpc/server/cmd/user/kitex_gen/user"
	"douyin_rpc/server/cmd/user/model"
	tools "douyin_rpc/server/cmd/user/tool"
	"errors"
	"fmt"
	"gorm.io/gorm"
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
	if req.QueryType == 1 {
		cache, err := model.QueryUserByIDWithCache(tx, req.UserId)
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
		cache, err := model.QueryUserByNameWithCache(tx, req.Username)
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
	cache, err := model.QueryUserByIDWithCache(tx, req.UserId)
	if err != nil {
		return nil, err
	}
	count, err := global.RelationClient.GetCount(ctx, &relation2.GetCountRequest{UserId: req.UserId})
	if err != nil {
		return nil, err
	}
	global.RelationClient.GetRelation(ctx, &relation2.GetRelationRequest{
		Id:           0,
		UserId:       req.AuthId,
		TargetId:     req.UserId,
		RelationType: 1,
		QueryType:    1,
	})
	resp.User = &user.User{
		Id:            cache.ID,
		Name:          cache.Name,
		FollowCount:   count.FollowCount,
		FollowerCount: count.FollowerCount,
		IsFollow:      false,
	}
	return
}
