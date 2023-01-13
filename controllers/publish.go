package controllers

import (
	"douyin/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostPublishActionForm struct {
	Token string `json:"token"`
	Data  uint   `json:"data"`
	Title string `json:"content"`
}

// PublishAction
// @Summary 获取聊天记录
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param publishAction body PostPublishActionForm true "视频信息"
// @Success 200 object utils.Response 成功后返回值
// @Failure 409 object utils.Response 失败后返回值
// @Router /douyin/publish/action [post]
func PublishAction(c *gin.Context) {
	var publishActionForm *PostPublishActionForm
	if err := c.ShouldBind(&publishActionForm); err != nil {
		response := utils.Response{
			StatusCode: 1,
			StatusMsg:  "data数据有误",
		}

		c.JSON(409, response)
		return
	}

	//TODO
	// 业务代码

	response := utils.Response{
		StatusCode: 0,
		StatusMsg:  "",
	}
	c.JSON(http.StatusOK, response)
	return
}

type PublishListResponse struct {
	utils.Response
	VideoList []*Video `json:"video_list"`
}

// PublishList
// @Summary 获取聊天记录
// @Tags 发布
// @version 1.0
// @Accept application/x-json-stream
// @Param user_id query int true "用户id"
// @Param token query string true "token"
// @Success 200 object FavoriteListResponse 成功后返回值
// @Failure 409 object FavoriteListResponse 失败后返回值
// @Router /douyin/publish/list [get]
func PublishList(c *gin.Context) {

	//TODO
	// 业务代码

	response := PublishListResponse{
		Response: utils.Response{
			StatusCode: 0,
			StatusMsg:  "",
		},
		VideoList: nil,
	}

	c.JSON(http.StatusOK, response)
	return
}
