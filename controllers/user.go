package controllers

import (
	"douyin/bootstrap/driver"
	"douyin/models"
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserInfoResponse struct {
	utils.Response
	User *User `json:"user"`
}

// UserInfo
// @Summary 用户信息
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object UserInfoResponse 成功后返回值
// @Failure 409 object UserInfoResponse 失败后返回值
// @Router /douyin/user [get]
func UserInfo(c *gin.Context) {

	successResponse := UserInfoResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		User: nil,
	}

	failureResponse := UserInfoResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		User: nil,
	}
	// 1、解析参数(不需要解析token，在中间件已经解析了)
	userID := c.Query("user_id")
	token := c.Query("token")
	auth, err := ParseToken(token)
	if err != nil {
		failureResponse.StatusMsg = "验证失败"
		c.JSON(409, failureResponse)
	}
	id, err := strconv.Atoi(userID)
	if err != nil || id < 0 {
		failureResponse.StatusMsg = "user_id 非数字"
		c.JSON(409, failureResponse)
		return
	}
	if id == 0 {
		id = int(auth.UserID)
	}
	// 2、查询数据库获取用户信息
	var user models.User
	if err := driver.Db.Debug().Model(user).Where(" id = ?", uint(id)).Find(&user).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 2.1、查询关注表，查看登录用户是否关注了user_id对应的用户
	isFollow := false
	var relation models.Relation
	if err := driver.Db.Debug().Model(relation).Where(" user_id = ? ", auth.UserID).Where("target_id = ?", user.ID).Where("exist=1").Where("type=1 or type=2").Find(&relation).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		return

	}
	if relation.ID != 0 {
		isFollow = true
	}
	// 3、填充返回的用户信息
	returnUser := User{
		ID:            user.ID,
		Name:          user.Name,
		FollowCount:   user.FollowCount,
		FollowerCount: user.FollowerCount,
		IsFollow:      isFollow,
	}
	// 4、返回
	successResponse.User = &returnUser
	c.JSON(http.StatusOK, successResponse)
	return
}

// UserRegister
// @Summary 用户注册
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/user/register [post]
func UserRegister(c *gin.Context) {

	successResponse := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}

	failureResponse := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}

	// 1、解析参数
	username := c.Query("username")
	password := c.Query("password")
	// 2、验证参数
	if !utils.VerifyParam(username, password) {
		failureResponse.StatusMsg = "用户名或密码长度大于32字符"
		c.JSON(409, failureResponse)
		return
	}
	// 3、生成盐值和加密后的密码
	salt := utils.Salt()
	cryptoPassword := utils.CryptUserPassword(password, salt)
	// 4、创建用户
	user := models.User{
		Model:         gorm.Model{},
		Name:          username,
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
	if err := tx.Model(user).Create(&user).Error; err != nil {
		tx.Rollback()
		failureResponse.StatusMsg = "创建用户失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	successResponse.UserId = user.ID
	// 5.2.1、用户登录
	auth := Auth{
		UserID:           user.ID,
		UserName:         username,
		FollowCount:      0,
		FollowerCount:    0,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token, err := MakeToken(&auth)
	if err != nil {
		tx.Rollback()
		failureResponse.StatusMsg = "创建token失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	successResponse.Token = token
	// 5.3、事务提交
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		failureResponse.StatusMsg = "创建用户失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 结果返回
	c.JSON(http.StatusOK, successResponse)

}

type UserLoginResponse struct {
	utils.Response
	Token  string `json:"token"`
	UserId uint   `json:"user_id"`
}

// UserLogin
// @Summary 用户登录
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/user/login [post]
func UserLogin(c *gin.Context) {
	//TODO
	// 业务代码

	successResponse := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}
	failureResponse := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}
	// 1、解析参数
	username := c.Query("username")
	password := c.Query("password")
	// 2、验证参数
	if !utils.VerifyParam(username, password) {
		failureResponse.StatusMsg = "用户名或密码长度大于32字符"
		c.JSON(409, failureResponse)
		return
	}
	// 3、解析用户名和密码是否与数据库一致
	// 3.1、根据用户名查询user信息，用户名必须唯一，可以添加unique索引。仅仅是查询，无需使用事务
	var user models.User
	if err := driver.Db.Debug().Model(user).Where(" name = ?", username).Find(&user).Error; err != nil {
		failureResponse.StatusMsg = "查询数据库失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 4、获取到盐值，加密后判断是否一致
	if !utils.VerifyUserPassword(user.Salt, password, user.Password) {
		failureResponse.StatusMsg = "验证失败"
		c.JSON(409, failureResponse)
		return
	}
	// 5、验证成功，创建token
	auth := Auth{
		UserID:           user.ID,
		UserName:         user.Name,
		FollowCount:      user.FollowCount,
		FollowerCount:    user.FollowerCount,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token, err := MakeToken(&auth)
	if err != nil {
		failureResponse.StatusMsg = "创建token失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 6、返回token
	successResponse.Token = token
	successResponse.UserId = user.ID
	c.JSON(http.StatusOK, successResponse)
}
