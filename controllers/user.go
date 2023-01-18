package controllers

import (
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserInfoResponse struct {
	utils.Response
	User *vo.User `json:"user"`
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
// @Router /douyin/user/ [get]
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
	info, err := service.UserInfo(uint(id), auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 4、返回
	successResponse.User = info
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
// @Router /douyin/user/register/ [post]
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
	auth, err := service.UserRegister(username, password)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	token, err := MakeToken(auth)
	successResponse.Token = token
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
// @Router /douyin/user/login/ [post]
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
	auth, err := service.UserLogin(username, password)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	token, err := MakeToken(auth)
	if err != nil {
		failureResponse.StatusMsg = "创建token失败" + err.Error()
		c.JSON(409, failureResponse)
		return
	}
	// 6、返回token
	successResponse.Token = token
	successResponse.UserId = auth.UserID
	c.JSON(http.StatusOK, successResponse)
}
