package service

import (
	auth2 "douyin/auth"
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"douyin/vo"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

type UserInfoFlow struct {
	UserId   uint
	AuthUser *auth2.Auth
	User     *vo.User
}

func UserInfo(userId uint, auth *auth2.Auth) (*vo.User, error) {
	return (&UserInfoFlow{
		UserId:   userId,
		AuthUser: auth,
		User:     nil,
	}).Do()
}
func (c *UserInfoFlow) Do() (*vo.User, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.info(); err != nil {
		return nil, err
	}
	return c.User, nil
}
func (c *UserInfoFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *UserInfoFlow) info() error {
	// 2、查询数据库获取用户信息
	tx := driver.Db.Debug()
	user, err := models.QueryUserByIDWithCache(tx, c.UserId)
	if err != nil {
		return err
	}
	// 2.1、查询关注表，查看登录用户是否关注了user_id对应的用户
	isFollow := false
	relation, err := models.QueryRelationByUserIDAndTargetIDWithCache(tx, c.AuthUser.UserID, user.ID)
	if err != nil {
		return err
	}

	if relation.ID != 0 {
		isFollow = true
	}
	// 3、填充返回的用户信息
	c.User = &vo.User{
		ID:            user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	// 4、返回
	return nil
}

type UserRegisterFlow struct {
	UserName string
	Password string
	auth     *auth2.Auth
}

func UserRegister(username string, password string) (*auth2.Auth, error) {
	return (&UserRegisterFlow{
		UserName: username,
		Password: password,
		auth:     nil,
	}).Do()
}
func (c *UserRegisterFlow) Do() (*auth2.Auth, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.register(); err != nil {
		return nil, err
	}
	return c.auth, nil
}
func (c *UserRegisterFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *UserRegisterFlow) register() error {
	salt := utils.Salt()
	cryptoPassword := utils.CryptUserPassword(c.Password, salt)
	// 4、创建用户
	user := models.User{
		Model:         gorm.Model{},
		Name:          c.UserName,
		FollowCount:   0,
		FollowerCount: 0,
		Password:      cryptoPassword,
		Salt:          salt,
	}
	// 5、gorm创建用户
	// 5.1、事务开始
	tx := driver.Db.Debug().Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 5.2、创建用户
	err := models.CreateUser(tx, &user)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 5.2.1、用户登录
	c.auth = &auth2.Auth{
		UserID:           user.ID,
		UserName:         c.UserName,
		FollowCount:      0,
		FollowerCount:    0,
		RegisteredClaims: jwt.RegisteredClaims{},
	}

	// 5.3、事务提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

type UserLoginFlow struct {
	UserName string
	Password string
	auth     *auth2.Auth
}

func UserLogin(username string, password string) (*auth2.Auth, error) {
	return (&UserLoginFlow{
		UserName: username,
		Password: password,
		auth:     nil,
	}).Do()
}
func (c *UserLoginFlow) Do() (*auth2.Auth, error) {
	if err := c.checkParam(); err != nil {
		return nil, err
	}
	if err := c.login(); err != nil {
		return nil, err
	}
	return c.auth, nil
}
func (c *UserLoginFlow) checkParam() error {
	//TODO
	// redis bitmap验证参数
	return nil
}

func (c *UserLoginFlow) login() error {

	tx := driver.Db.Debug()
	user, err := models.QueryUserByNameWithCache(tx, c.UserName)
	if err != nil {
		return err
	}

	// 4、获取到盐值，加密后判断是否一致
	if !utils.VerifyUserPassword(user.Salt, c.Password, user.Password) {

		return errors.New("密码错误")
	}
	// 5、验证成功，创建token
	c.auth = &auth2.Auth{
		UserID:           user.ID,
		UserName:         user.Name,
		FollowCount:      user.FollowCount,
		FollowerCount:    user.FollowerCount,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	return nil
}
