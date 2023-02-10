// Code generated by hertz generator.

package api

import (
	"context"
	"douyin_rpc/server/cmd/api/constant"
	"douyin_rpc/server/cmd/api/global"
	"douyin_rpc/server/cmd/api/kitex_gen/user"
	"douyin_rpc/server/cmd/api/middleware"
	"douyin_rpc/server/cmd/api/models"
	"github.com/golang-jwt/jwt"
	"strconv"
	"time"

	api "douyin_rpc/server/cmd/api/biz/model/api"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Login
// @Summary 用户登录
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 object user.LoginResponse 成功后返回值
// @Failure 409 object user.LoginResponse 失败后返回值
// @Router /douyin/user/login/ [post]
// @router /douyin/user/login/ [POST]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	resp, err := global.UserClient.Login(ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return
	}
	j := middleware.NewJWT()
	claims := models.CustomClaims{
		ID: resp.UserId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + constant.ThirtyDays,
			Issuer:    constant.JWTIssuer,
		},
	}
	response := new(api.LoginResponse)
	token, err := j.CreateToken(claims)
	if err != nil {
		response.StatusCode = 1
		response.StatusMsg = err.Error()
		return
	}
	response.StatusCode = 0
	response.UserID = resp.UserId
	response.Token = token
	c.JSON(consts.StatusOK, response)
}

// Register
// @Summary 用户注册
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param username query string true "用户名"
// @Param password query string true "密码"
// @Success 200 object user.RegisterResponse 成功后返回值
// @Failure 409 object user.RegisterResponse 失败后返回值
// @Router /douyin/user/register/ [post]
// @router /douyin/user/register/ [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	response := new(api.RegisterResponse)
	resp, err := global.UserClient.Register(ctx, &user.RegisterRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil || resp.UserId == 0 {
		response.StatusCode = 1
		response.StatusMsg = err.Error()
		return
	}
	//resp := new(api.LoginResponse)
	j := middleware.NewJWT()
	claims := models.CustomClaims{
		ID: resp.UserId,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + constant.ThirtyDays,
			Issuer:    constant.JWTIssuer,
		},
	}

	token, err := j.CreateToken(claims)
	if err != nil {

		return
	}
	response.StatusCode = 0
	response.UserID = resp.UserId
	response.Token = token
	c.JSON(consts.StatusOK, response)
}

// GetUser
// @Summary 用户信息
// @Tags 用户
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object user.UserResponse 成功后返回值
// @Failure 409 object user.UserResponse 失败后返回值
// @Router /douyin/user/ [get]
// @router /douyin/user/ [GET]
func GetUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UserRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	value, exist := c.Get("accountID")
	if !exist {
		return
	}
	userID, err := strconv.ParseInt(req.UserID, 0, 64)
	if err != nil {
		return
	}
	resp, err := global.UserClient.User(ctx, &user.UserRequest{
		UserId: userID,
		AuthId: value.(int64),
	})
	if err != nil {
		return
	}
	//resp := new(api.UserResponse)

	c.JSON(consts.StatusOK, resp)
}
