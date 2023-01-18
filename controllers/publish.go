package controllers

import (
	auth2 "douyin/auth"
	"douyin/service"
	"douyin/utils"
	"douyin/vo"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// PublishAction
// @Summary 发布视频
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param token formData string true "token"
// @Param data formData file true "视频"
// @Param title formData string true "标题"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/publish/action/ [post]
func PublishAction(c *gin.Context) {

	auth := auth2.Auth{}.GetAuth(c)

	successResponse := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	failureResponse := utils.Response{
		StatusCode: 1,
		StatusMsg:  "",
	}
	// 1、解析参数
	form, err := c.MultipartForm()
	if err != nil {
		failureResponse.StatusMsg = "解析form表单错误"
		c.JSON(409, failureResponse)
		return
	}
	title := c.PostForm("title")
	// 2、验证参数
	data := form.File["data"]
	if len(data) != 1 {
		failureResponse.StatusMsg = "视频数量只能为1"
		c.JSON(409, failureResponse)
		return
	}
	err = service.PublishActionPost(title, data[0], &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 4、装配返回值
	c.JSON(http.StatusOK, successResponse)
	return
}

type PublishListResponse struct {
	utils.Response
	VideoList []vo.Video `json:"video_list"`
}

// PublishList
// @Summary 获取发布视频列表
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/publish/list/ [get]
func PublishList(c *gin.Context) {

	auth := auth2.Auth{}.GetAuth(c)

	successResponse := PublishListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: nil,
	}
	failureResponse := PublishListResponse{
		Response: utils.Response{
			StatusCode: 1,
			StatusMsg:  "",
		},
		VideoList: nil,
	}
	// 1、解析参数
	userId := c.Query("user_id")

	// 2、验证参数
	// user_id必须为正整数
	userID, err := strconv.Atoi(userId)
	if err != nil || userID < 0 {
		failureResponse.StatusMsg = "user_id必须为正整数"
		c.JSON(409, failureResponse)
		return
	}
	if userID == 0 {
		userID = int(auth.UserID)
	}
	// 3、具体业务
	// 3.1、查询用户的所有视频
	returnVideoList, err := service.PublishListGet(uint(userID), &auth)
	if err != nil {
		c.JSON(409, failureResponse)
		return
	}
	// 4、装配返回值
	successResponse.VideoList = returnVideoList
	c.JSON(http.StatusOK, successResponse)
	return
}
