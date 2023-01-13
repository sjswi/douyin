package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
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
// @Param publishAction body PostPublishActionForm true "视频信息"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/user/info [get]
func UserInfo(c *gin.Context) {
	//TODO
	// 业务代码
	response := UserInfoResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		User: nil,
	}
	c.JSON(http.StatusOK, response)
}

// UserRegister
// @Summary 用户注册
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param publishAction body PostPublishActionForm true "视频信息"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/user/register [post]
func UserRegister(c *gin.Context) {
	//TODO
	// 业务代码

	response := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}
	c.JSON(http.StatusOK, response)
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

	response := UserLoginResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		Token:  "",
		UserId: 0,
	}
	c.JSON(http.StatusOK, response)
}
